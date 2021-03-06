// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package updater

import (
	"context"
	"encoding/hex"
	"sort"

	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/common/logging"
	"go.chromium.org/luci/server/span"

	"infra/appengine/weetbix/internal/analysis"
	"infra/appengine/weetbix/internal/bugs"
	"infra/appengine/weetbix/internal/clustering"
	"infra/appengine/weetbix/internal/clustering/algorithms"
	"infra/appengine/weetbix/internal/clustering/algorithms/rulesalgorithm"
	"infra/appengine/weetbix/internal/clustering/rules"
	"infra/appengine/weetbix/internal/clustering/rules/lang"
	"infra/appengine/weetbix/internal/clustering/runs"
	"infra/appengine/weetbix/internal/config/compiledcfg"
	pb "infra/appengine/weetbix/proto/v1"
)

// testnameThresholdInflationPercent is the percentage factor by which
// the bug filing threshold is inflated when applied to test-name clusters.
// This is to bias bug-filing towards failure reason clusters, which are
// seen as generally better scoped and more actionable (because they
// focus on one reason for the test failing.)
//
// The value of 34% was selected as it is sufficient to inflate any threshold
// values which are a '3' (e.g. CV runs rejected) to a '4'. Otherwise integer
// discretization of the statistics would cancel out any intended bias.
//
// If changing this value, please also update the comment in
// project_config.proto.
const testnameThresholdInflationPercent = 34

// BugManager implements bug creation and bug updates for a bug-tracking
// system. The BugManager determines bug content and priority given a
// cluster.
type BugManager interface {
	// Create creates a new bug for the given request, returning its name,
	// or any encountered error.
	Create(ctx context.Context, cluster *bugs.CreateRequest) (string, error)
	// Update updates the specified list of bugs.
	Update(ctx context.Context, bugs []*bugs.BugToUpdate) error
}

// BugUpdater performs updates to Monorail bugs and failure association
// rules to keep them in sync with clusters generated by analysis.
type BugUpdater struct {
	// project is the LUCI project to act on behalf of.
	project string
	// analysisClient provides access to cluster analysis.
	analysisClient AnalysisClient
	// managers stores the manager responsible for updating bugs for each
	// bug tracking system (monorail, buganizer, etc.).
	managers map[string]BugManager
	// projectCfg is the snapshot of project configuration to use for
	// the auto-bug filing run.
	projectCfg *compiledcfg.ProjectConfig
	// MaxBugsFiledPerRun is the maximum number of bugs to file each time
	// BugUpdater runs. This throttles the rate of changes to monorail.
	MaxBugsFiledPerRun int
}

// NewBugUpdater initialises a new BugUpdater. The specified impact thresholds are used
// when determining whether to a file a bug.
func NewBugUpdater(project string, mgrs map[string]BugManager, ac AnalysisClient, projectCfg *compiledcfg.ProjectConfig) *BugUpdater {
	return &BugUpdater{
		project:            project,
		managers:           mgrs,
		analysisClient:     ac,
		projectCfg:         projectCfg,
		MaxBugsFiledPerRun: 1, // Default value.
	}
}

// Run updates files/updates bugs to match high-impact clusters as
// identified by analysis. Each bug has a corresponding failure association
// rule.
// The passed progress should reflect the progress of re-clustering as captured
// in the latest analysis.
func (b *BugUpdater) Run(ctx context.Context, progress *runs.ReclusteringProgress) error {
	if progress.IsReclusteringToNewAlgorithms() {
		logging.Warningf(ctx, "Auto-bug filing paused for project %s as re-clustering to new algorithms is in progress.", b.project)
		return nil
	}
	if progress.IsReclusteringToNewConfig() {
		logging.Warningf(ctx, "Auto-bug filing paused for project %s as re-clustering to new configuration is in progress.", b.project)
		return nil
	}
	if algorithms.AlgorithmsVersion != progress.LatestAlgorithmsVersion {
		logging.Warningf(ctx, "Auto-bug filing paused for project %s as bug-filing is running mismatched algorithms version %v (want %v).",
			b.project, algorithms.AlgorithmsVersion, progress.LatestAlgorithmsVersion)
		return nil
	}
	if !b.projectCfg.LastUpdated.Equal(progress.LatestConfigVersion) {
		logging.Warningf(ctx, "Auto-bug filing paused for project %s as bug-filing is running mismatched config version %v (want %v).",
			b.project, b.projectCfg.LastUpdated, progress.LatestConfigVersion)
		return nil
	}

	ruleByID, err := b.readActiveFailureAssociationRules(ctx)
	if err != nil {
		return errors.Annotate(err, "read active failure association rules").Err()
	}

	// blockedSourceClusterIDs is the set of source cluster IDs for which
	// filing new bugs should be suspended.
	blockedSourceClusterIDs := make(map[string]struct{})
	for _, r := range ruleByID {
		if !progress.IncorporatesRulesVersion(r.CreationTime) {
			// If a bug cluster was recently filed for a source cluster, and
			// re-clustering and analysis is not yet complete (to move the
			// impact from the source cluster to the bug cluster), do not file
			// another bug for the source cluster.
			// (Of course, if a bug cluster was filed for a source cluster,
			// but the bug cluster's failure association rule was subsequently
			// modified (e.g. narrowed), it is allowed to file another bug
			// if the residual impact justifies it.)
			blockedSourceClusterIDs[r.SourceCluster.Key()] = struct{}{}
		}
	}

	// We want to read analysis for two categories of clusters:
	// - Bug Clusters: to update the priority of filed bugs.
	// - Impactful Suggested Clusters: if any suggested clusters have
	//    reached the threshold to file a new bug for, we want to read
	//    them, so we can file a bug.
	clusterSummaries, err := b.analysisClient.ReadImpactfulClusters(ctx, analysis.ImpactfulClusterReadOptions{
		Project:                  b.project,
		Thresholds:               b.projectCfg.Config.BugFilingThreshold,
		AlwaysIncludeBugClusters: true,
	})
	if err != nil {
		return errors.Annotate(err, "read impactful clusters").Err()
	}

	sortByBugFilingPreference(clusterSummaries)

	var toCreateBugsFor []*analysis.ClusterSummary
	impactByRuleID := make(map[string]*bugs.ClusterImpact)
	for _, clusterSummary := range clusterSummaries {
		if clusterSummary.ClusterID.IsBugCluster() {
			if clusterSummary.ClusterID.Algorithm == rulesalgorithm.AlgorithmName {
				// Use only impact from latest algorithm version.
				ruleID := clusterSummary.ClusterID.ID
				impactByRuleID[ruleID] = ExtractResidualPreWeetbixImpact(clusterSummary)
			}
			// Never file another bug for a bug cluster.
			continue
		}

		// Was a bug recently filed for this source cluster?
		// We want to avoid race conditions whereby we file multiple bug
		// clusters for the same source cluster, because re-clustering and
		// re-analysis has not yet run and moved residual impact from the
		// source (suggested) cluster to the bug cluster.
		_, ok := blockedSourceClusterIDs[clusterSummary.ClusterID.Key()]
		if ok {
			// Do not file a bug.
			continue
		}

		// Only file a bug if the residual impact exceeds the threshold.
		impact := ExtractResidualPreWeetbixImpact(clusterSummary)
		bugFilingThreshold := b.projectCfg.Config.BugFilingThreshold
		if clusterSummary.ClusterID.IsTestNameCluster() {
			// Use an inflated threshold for test name clusters to bias
			// bug creation towards failure reason clusters.
			bugFilingThreshold =
				bugs.InflateThreshold(b.projectCfg.Config.BugFilingThreshold,
					testnameThresholdInflationPercent)
		}
		if !impact.MeetsThreshold(bugFilingThreshold) {
			continue
		}

		toCreateBugsFor = append(toCreateBugsFor, clusterSummary)
	}

	bugsFiled := 0
	for _, clusterSummary := range toCreateBugsFor {
		// Throttle how many bugs may be filed each time.
		if bugsFiled >= b.MaxBugsFiledPerRun {
			break
		}
		created, err := b.createBug(ctx, clusterSummary)
		if err != nil {
			return err
		}
		if created {
			bugsFiled++
		}
	}

	// Iterate over all active bug clusters (except those we just created).
	bugUpdatesBySystem := make(map[string][]*bugs.BugToUpdate)
	for id, r := range ruleByID {
		impact, ok := impactByRuleID[id]
		if !ok {
			// If there is no analysis, this usually means the cluster is
			// empty, so we should use empty impact.
			impact = &bugs.ClusterImpact{}
		}

		// Only update the bug if it is managed by this rule.
		if !r.IsManagingBug {
			continue
		}

		// Only update the bug if re-clustering and analysis ran on the latest
		// version of this failure association rule. This avoids bugs getting
		// erroneous priority changes while impact information is incomplete.
		if !progress.IncorporatesRulesVersion(r.PredicateLastUpdated) {
			continue
		}

		bugUpdates := bugUpdatesBySystem[r.BugID.System]
		bugUpdates = append(bugUpdates, &bugs.BugToUpdate{
			BugName: r.BugID.ID,
			Impact:  impact,
		})
		bugUpdatesBySystem[r.BugID.System] = bugUpdates
	}

	for system, bugsToUpdate := range bugUpdatesBySystem {
		manager, ok := b.managers[system]
		if !ok {
			logging.Warningf(ctx, "Encountered bug(s) with an unrecognised manager: %q", manager)
			continue
		}
		if err := manager.Update(ctx, bugsToUpdate); err != nil {
			return err
		}
	}
	return nil
}

// sortByBugFilingPreference sorts cluster summaries based on our preference
// to file bugs for these clusters.
func sortByBugFilingPreference(cs []*analysis.ClusterSummary) {
	// The current ranking approach prefers filing bugs for clusters with more
	// impact, with a bias towards reason clusters.
	//
	// The order of this ranking is only important where there are
	// multiple competing clusters which meet the bug-filing threshold.
	// As bug filing runs relatively often, except in cases of contention,
	// the first bug to meet the threshold will be filed.
	sort.Slice(cs, func(i, j int) bool {
		presubmitRejects := func(cs *analysis.ClusterSummary) analysis.Counts { return cs.PresubmitRejects7d }
		testRuns := func(cs *analysis.ClusterSummary) analysis.Counts { return cs.TestRunFails7d }
		failures := func(cs *analysis.ClusterSummary) analysis.Counts { return cs.Failures7d }

		if equal, less := rankByMetric(cs[i], cs[j], presubmitRejects); !equal {
			return less
		}
		if equal, less := rankByMetric(cs[i], cs[j], testRuns); !equal {
			return less
		}
		if equal, less := rankByMetric(cs[i], cs[j], failures); !equal {
			return less
		}
		// If all else fails, sort by cluster ID. This is mostly to ensure
		// the code behaves deterministically when under unit testing.
		if cs[i].ClusterID.Algorithm != cs[j].ClusterID.Algorithm {
			return cs[i].ClusterID.Algorithm < cs[j].ClusterID.Algorithm
		}
		return cs[i].ClusterID.ID < cs[j].ClusterID.ID
	})
}

func rankByMetric(a, b *analysis.ClusterSummary, accessor func(*analysis.ClusterSummary) analysis.Counts) (equal bool, less bool) {
	valueA := accessor(a).ResidualPreWeetbix
	valueB := accessor(b).ResidualPreWeetbix
	// If one cluster we are comparing with is a test name cluster,
	// give the other cluster an impact boost in the comparison, so
	// that we bias towards filing it (instead of the test name cluster).
	if b.ClusterID.IsTestNameCluster() {
		valueA = (valueA * (100 + testnameThresholdInflationPercent)) / 100
	}
	if a.ClusterID.IsTestNameCluster() {
		valueB = (valueB * (100 + testnameThresholdInflationPercent)) / 100
	}
	equal = (valueA == valueB)
	// a less than b in the sort order is defined as a having more impact
	// than b, so that cluster summaries are sorted in descending impact order.
	less = (valueA > valueB)
	return equal, less
}

// createBug files a new bug for the given suggested cluster,
// and stores the association from bug to failures through a new
// failure association rule.
func (b *BugUpdater) createBug(ctx context.Context, cs *analysis.ClusterSummary) (created bool, err error) {
	alg, err := algorithms.SuggestingAlgorithm(cs.ClusterID.Algorithm)
	if err == algorithms.ErrAlgorithmNotExist {
		// The cluster is for an old algorithm that no longer exists, or
		// for a new algorithm that is not known by us yet.
		// Do not file a bug. This is not an error, it is expected during
		// algorithm version changes.
		return false, nil
	}

	summary := clusterSummaryFromAnalysis(cs)

	// Double-check the failure matches the cluster. Generating a
	// failure association rule that does not match the suggested cluster
	// could result in indefinite creation of new bugs, as the system
	// will repeatedly create new failure association rules for the
	// same suggested cluster.
	// Mismatches should usually be transient as re-clustering will fix
	// up any incorrect clustering.
	if hex.EncodeToString(alg.Cluster(b.projectCfg, &summary.Example)) != cs.ClusterID.ID {
		return false, errors.New("example failure did not match cluster ID")
	}
	rule, err := b.generateFailureAssociationRule(alg, &summary.Example)
	if err != nil {
		return false, errors.Annotate(err, "obtain failure association rule").Err()
	}

	ruleID, err := rules.GenerateID()
	if err != nil {
		return false, errors.Annotate(err, "generating rule ID").Err()
	}

	description, err := alg.ClusterDescription(b.projectCfg, summary)
	if err != nil {
		return false, errors.Annotate(err, "prepare bug description").Err()
	}

	request := &bugs.CreateRequest{
		Description: description,
		Impact:      ExtractResidualPreWeetbixImpact(cs),
	}

	// For now, the only issue system supported is monorail.
	system := bugs.MonorailSystem
	mgr := b.managers[system]
	name, err := mgr.Create(ctx, request)
	if err == bugs.ErrCreateSimulated {
		// Create did not do anything because it is in simulation mode.
		// This is expected.
		return false, nil
	}
	if err != nil {
		return false, errors.Annotate(err, "create issue in %v", mgr).Err()
	}

	// Create a failure association rule associating the failures with a bug.
	r := &rules.FailureAssociationRule{
		Project:        b.project,
		RuleID:         ruleID,
		RuleDefinition: rule,
		BugID:          bugs.BugID{System: system, ID: name},
		IsActive:       true,
		IsManagingBug:  true,
		SourceCluster:  cs.ClusterID,
	}
	create := func(ctx context.Context) error {
		user := rules.WeetbixSystem
		return rules.Create(ctx, r, user)
	}
	if _, err := span.ReadWriteTransaction(ctx, create); err != nil {
		return false, errors.Annotate(err, "create bug cluster").Err()
	}

	return true, nil
}

func clusterSummaryFromAnalysis(cs *analysis.ClusterSummary) *clustering.ClusterSummary {
	example := clustering.Failure{
		TestID: cs.ExampleTestID(),
	}
	if cs.ExampleFailureReason.Valid {
		example.Reason = &pb.FailureReason{PrimaryErrorMessage: cs.ExampleFailureReason.StringVal}
	}
	// A list of 5 commonly occuring tests are included in bugs created
	// for failure reason clusters, to improve searchability by test name.
	var topTests []string
	for _, tt := range cs.TopTestIDs {
		topTests = append(topTests, tt.Value)
	}
	return &clustering.ClusterSummary{
		Example:  example,
		TopTests: topTests,
	}
}

func (b *BugUpdater) generateFailureAssociationRule(alg algorithms.Algorithm, failure *clustering.Failure) (string, error) {
	rule := alg.FailureAssociationRule(b.projectCfg, failure)

	// Check the generated rule is valid and matches the failure.
	// An improperly generated failure association rule could result
	// in uncontrolled creation of new bugs.
	expr, err := lang.Parse(rule)
	if err != nil {
		return "", errors.Annotate(err, "rule generated by %s did not parse", alg.Name()).Err()
	}
	match := expr.Evaluate(failure)
	if !match {
		return "", errors.New("rule generated by %s did not match example failure")
	}
	return rule, nil
}

func (b *BugUpdater) readActiveFailureAssociationRules(ctx context.Context) (map[string]*rules.FailureAssociationRule, error) {
	rs, err := rules.ReadActive(span.Single(ctx), b.project)
	if err != nil {
		return nil, err
	}

	ruleByID := make(map[string]*rules.FailureAssociationRule)
	for _, r := range rs {
		ruleByID[r.RuleID] = r
	}
	return ruleByID, nil
}
