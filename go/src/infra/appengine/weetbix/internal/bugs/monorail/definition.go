// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package monorail

import (
	"fmt"
	"regexp"
	"strings"

	"google.golang.org/genproto/protobuf/field_mask"

	"infra/appengine/weetbix/internal/bugs"
	"infra/appengine/weetbix/internal/clustering"
	configpb "infra/appengine/weetbix/internal/config/proto"
	mpb "infra/monorailv2/api/v3/api_proto"
)

const (
	DescriptionTemplate = `%s

This bug has been automatically filed by Weetbix in response to a cluster of test failures.`

	LinkTemplate = `See failure impact and configure the failure association rule for this bug at: %s`
)

const (
	manualPriorityLabel = "Weetbix-Manual-Priority"
	restrictViewLabel   = "Restrict-View-Google"
	autoFiledLabel      = "Weetbix-Auto-Filed"
)

// whitespaceRE matches blocks of whitespace, including new lines tabs and
// spaces.
var whitespaceRE = regexp.MustCompile(`[ \t\n]+`)

// priorityRE matches chromium monorail priority values.
var priorityRE = regexp.MustCompile(`^Pri-([0123])$`)

// AutomationUsers are the identifiers of Weetbix automation users in monorail.
var AutomationUsers = []string{
	"users/3816576959", // chops-weetbix@appspot.gserviceaccount.com
	"users/4149141945", // chops-weetbix-dev@appspot.gserviceaccount.com
}

// ChromiumDefaultAssignee is the default issue assignee for chromium.
// This should be deleted in future when auto-assignment is implemented.
const ChromiumDefaultAssignee = "users/2581171748" // mwarton@chromium.org

// VerifiedStatus is that status of bugs that have been fixed and verified.
const VerifiedStatus = "Verified"

// AssignedStatus is the status of bugs that are open and assigned to an owner.
const AssignedStatus = "Assigned"

// UntriagedStatus is the status of bugs that have just been opened.
const UntriagedStatus = "Untriaged"

// Generator provides access to a methods to generate a new bug and/or bug
// updates for a cluster.
type Generator struct {
	// The impact of the cluster to generate monorail changes for.
	impact *bugs.ClusterImpact
	// The monorail configuration to use.
	monorailCfg *configpb.MonorailProject
}

// NewGenerator initialises a new Generator.
func NewGenerator(impact *bugs.ClusterImpact, monorailCfg *configpb.MonorailProject) (*Generator, error) {
	if len(monorailCfg.Priorities) == 0 {
		return nil, fmt.Errorf("invalid configuration for monorail project %q; no monorail priorities configured", monorailCfg.Project)
	}
	return &Generator{
		impact:      impact,
		monorailCfg: monorailCfg,
	}, nil
}

// PrepareNew prepares a new bug from the given cluster. title and description
// are the cluster-specific bug title and description.
func (g *Generator) PrepareNew(description *clustering.ClusterDescription) *mpb.MakeIssueRequest {
	issue := &mpb.Issue{
		Summary: fmt.Sprintf("Tests are failing: %v", sanitiseTitle(description.Title, 150)),
		State:   mpb.IssueContentState_ACTIVE,
		Status:  &mpb.Issue_StatusValue{Status: UntriagedStatus},
		FieldValues: []*mpb.FieldValue{
			{
				Field: g.priorityFieldName(),
				Value: g.clusterPriority(),
			},
		},
		Labels: []*mpb.Issue_LabelValue{{
			Label: restrictViewLabel,
		}, {
			Label: autoFiledLabel,
		}},
	}
	for _, fv := range g.monorailCfg.DefaultFieldValues {
		issue.FieldValues = append(issue.FieldValues, &mpb.FieldValue{
			Field: fmt.Sprintf("projects/%s/fieldDefs/%v", g.monorailCfg.Project, fv.FieldId),
			Value: fv.Value,
		})
	}
	if g.monorailCfg.Project == "chromium" {
		// mwarton@chromium.org in both prod and staging monorail.
		issue.Owner = &mpb.Issue_UserValue{User: ChromiumDefaultAssignee}
	}

	return &mpb.MakeIssueRequest{
		Parent:      fmt.Sprintf("projects/%s", g.monorailCfg.Project),
		Issue:       issue,
		Description: fmt.Sprintf(DescriptionTemplate, description.Description),
		NotifyType:  mpb.NotifyType_EMAIL,
	}
}

// LinkCommentRequest represents a request that adds links to Weetbix to
// a monorail bug.
type LinkCommentRequest struct {
	// The GAE app id, e.g. "chops-weetbix".
	AppID string
	// The LUCI Project.
	Project string
	// The internal bug name, e.g. "chromium/100".
	BugName string
}

// PrepareLinkComment prepares a request that adds links to Weetbix to
// a monorail bug.
func PrepareLinkComment(request LinkCommentRequest) (*mpb.ModifyIssuesRequest, error) {
	issueName, err := toMonorailIssueName(request.BugName)
	if err != nil {
		return nil, err
	}

	bugLink := fmt.Sprintf("https://%s.appspot.com/b/%s", request.AppID, request.BugName)

	comment := fmt.Sprintf(LinkTemplate, bugLink)

	result := &mpb.ModifyIssuesRequest{
		Deltas: []*mpb.IssueDelta{
			{
				Issue: &mpb.Issue{
					Name: issueName,
				},
				UpdateMask: &field_mask.FieldMask{},
			},
		},
		NotifyType:     mpb.NotifyType_NO_NOTIFICATION,
		CommentContent: comment,
	}
	return result, nil
}

func (g *Generator) priorityFieldName() string {
	return fmt.Sprintf("projects/%s/fieldDefs/%v", g.monorailCfg.Project, g.monorailCfg.PriorityFieldId)
}

// NeedsUpdate determines if the bug for the given cluster needs to be updated.
func (g *Generator) NeedsUpdate(issue *mpb.Issue) bool {
	// Bugs must have restrict view label to be updated.
	if !hasLabel(issue, restrictViewLabel) {
		return false
	}
	// Cases that a bug may be updated follow.
	switch {
	case !g.isCompatibleWithVerified(issueVerified(issue)):
		return true
	case !hasLabel(issue, manualPriorityLabel) &&
		!issueVerified(issue) &&
		!g.isCompatibleWithPriority(g.IssuePriority(issue)):
		// The priority has changed on a cluster which is not verified as fixed
		// and the user isn't manually controlling the priority.
		return true
	default:
		return false
	}
}

// MakeUpdate prepares an updated for the bug associated with a given cluster.
// Must ONLY be called if NeedsUpdate(...) returns true.
func (g *Generator) MakeUpdate(issue *mpb.Issue, comments []*mpb.Comment) *mpb.ModifyIssuesRequest {
	delta := &mpb.IssueDelta{
		Issue: &mpb.Issue{
			Name: issue.Name,
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{},
		},
	}

	var commentary []string
	notify := false
	issueVerified := issueVerified(issue)
	if !g.isCompatibleWithVerified(issueVerified) {
		// Verify or reopen the issue.
		comment := g.prepareBugVerifiedUpdate(issue, delta)
		commentary = append(commentary, comment)
		notify = true
		// After the update, whether the issue was verified will have changed.
		issueVerified = g.clusterResolved()
	}
	if !hasLabel(issue, manualPriorityLabel) &&
		!issueVerified &&
		!g.isCompatibleWithPriority(g.IssuePriority(issue)) {

		if hasManuallySetPriority(comments) {
			// We were not the last to update the priority of this issue.
			// Set the 'manually controlled priority' label to reflect
			// the state of this bug and avoid further attempts to update.
			comment := prepareManualPriorityUpdate(issue, delta)
			commentary = append(commentary, comment)
		} else {
			// We were the last to update the bug priority.
			// Apply the priority update.
			comment := g.preparePriorityUpdate(issue, delta)
			commentary = append(commentary, comment)
			// Notify if new priority is higher than existing priority.
			notify = notify || g.isHigherPriority(g.clusterPriority(), g.IssuePriority(issue))
		}
	}

	update := &mpb.ModifyIssuesRequest{
		Deltas: []*mpb.IssueDelta{
			delta,
		},
		NotifyType:     mpb.NotifyType_NO_NOTIFICATION,
		CommentContent: strings.Join(commentary, "\n\n"),
	}
	if notify {
		update.NotifyType = mpb.NotifyType_EMAIL
	}
	return update
}

func (g *Generator) prepareBugVerifiedUpdate(issue *mpb.Issue, update *mpb.IssueDelta) string {
	resolved := g.clusterResolved()
	var status string
	var message strings.Builder
	if resolved {
		status = VerifiedStatus

		oldPriorityIndex := len(g.monorailCfg.Priorities) - 1
		// A priority index of len(g.monorailCfg.Priorities) indicates
		// a priority lower than the lowest defined priority (i.e. bug verified.)
		newPriorityIndex := len(g.monorailCfg.Priorities)

		message.WriteString("Because:\n")
		message.WriteString(g.priorityDecreaseJustification(oldPriorityIndex, newPriorityIndex))
		message.WriteString("Weetbix is marking the issue verified.")
	} else {
		if issue.GetOwner().GetUser() != "" {
			status = AssignedStatus
		} else {
			status = UntriagedStatus
		}

		// A priority index of len(g.monorailCfg.Priorities) indicates
		// a priority lower than the lowest defined priority (i.e. bug verified.)
		oldPriorityIndex := len(g.monorailCfg.Priorities)
		newPriorityIndex := len(g.monorailCfg.Priorities) - 1

		message.WriteString("Because:\n")
		message.WriteString(g.priorityIncreaseJustification(oldPriorityIndex, newPriorityIndex))
		message.WriteString("Weetbix has re-opened the bug.")
	}
	update.Issue.Status = &mpb.Issue_StatusValue{Status: status}
	update.UpdateMask.Paths = append(update.UpdateMask.Paths, "status")
	return message.String()
}

func prepareManualPriorityUpdate(issue *mpb.Issue, update *mpb.IssueDelta) string {
	update.Issue.Labels = []*mpb.Issue_LabelValue{{
		Label: manualPriorityLabel,
	}}
	update.UpdateMask.Paths = append(update.UpdateMask.Paths, "labels")
	return fmt.Sprintf("The bug priority has been manually set. To re-enable automatic priority updates by Weetbix, remove the %s label.", manualPriorityLabel)
}

func (g *Generator) preparePriorityUpdate(issue *mpb.Issue, update *mpb.IssueDelta) string {
	newPriority := g.clusterPriority()

	update.Issue.FieldValues = []*mpb.FieldValue{
		{
			Field: g.priorityFieldName(),
			Value: newPriority,
		},
	}
	update.UpdateMask.Paths = append(update.UpdateMask.Paths, "field_values")

	oldPriority := g.IssuePriority(issue)
	oldPriorityIndex := g.indexOfPriority(oldPriority)
	newPriorityIndex := g.indexOfPriority(newPriority)

	if newPriorityIndex < oldPriorityIndex {
		var message strings.Builder
		message.WriteString("Because:\n")
		message.WriteString(g.priorityIncreaseJustification(oldPriorityIndex, newPriorityIndex))
		message.WriteString(fmt.Sprintf("Weetbix has increased the bug priority from %v to %v.", oldPriority, newPriority))
		return message.String()
	} else {
		var message strings.Builder
		message.WriteString("Because:\n")
		message.WriteString(g.priorityDecreaseJustification(oldPriorityIndex, newPriorityIndex))
		message.WriteString(fmt.Sprintf("Weetbix has decreased the bug priority from %v to %v.", oldPriority, newPriority))
		return message.String()
	}
}

// hasManuallySetPriority returns whether the the given issue has a manually
// controlled priority, based on its comments.
func hasManuallySetPriority(comments []*mpb.Comment) bool {
	// Example comment showing a user changing priority:
	// {
	// 	name: "projects/chromium/issues/915761/comments/1"
	// 	state: ACTIVE
	// 	type: COMMENT
	// 	commenter: "users/2627516260"
	// 	create_time: {
	// 	  seconds: 1632111572
	// 	}
	// 	amendments: {
	// 	  field_name: "Labels"
	// 	  new_or_delta_value: "Pri-1"
	// 	}
	// }
	for i := len(comments) - 1; i >= 0; i-- {
		c := comments[i]

		isManualPriorityUpdate := false
		isRevertToAutomaticPriority := false
		for _, a := range c.Amendments {
			if a.FieldName == "Labels" {
				deltaLabels := strings.Split(a.NewOrDeltaValue, " ")
				for _, lbl := range deltaLabels {
					if lbl == "-"+manualPriorityLabel {
						isRevertToAutomaticPriority = true
					}
					if priorityRE.MatchString(lbl) {
						if !isAutomationUser(c.Commenter) {
							isManualPriorityUpdate = true
						}
					}
				}
			}
		}
		if isRevertToAutomaticPriority {
			return false
		}
		if isManualPriorityUpdate {
			return true
		}
	}
	// No manual changes to priority indicates the bug is still under
	// automatic control.
	return false
}

func isAutomationUser(user string) bool {
	for _, u := range AutomationUsers {
		if u == user {
			return true
		}
	}
	return false
}

// hasLabel returns whether the bug the specified label.
func hasLabel(issue *mpb.Issue, label string) bool {
	for _, l := range issue.Labels {
		if l.Label == label {
			return true
		}
	}
	return false
}

// IssuePriority returns the priority of the given issue.
func (g *Generator) IssuePriority(issue *mpb.Issue) string {
	priorityFieldName := g.priorityFieldName()
	for _, fv := range issue.FieldValues {
		if fv.Field == priorityFieldName {
			return fv.Value
		}
	}
	return ""
}

func issueVerified(issue *mpb.Issue) bool {
	return issue.Status.Status == VerifiedStatus
}

// isHigherPriority returns whether priority p1 is higher than priority p2.
// The passed strings are the priority field values as used in monorail. These
// must be matched against monorail project configuration in order to
// identify the ordering of the priorities.
func (g *Generator) isHigherPriority(p1 string, p2 string) bool {
	i1 := g.indexOfPriority(p1)
	i2 := g.indexOfPriority(p2)
	// Priorities are configured from highest to lowest, so higher priorities
	// have lower indexes.
	return i1 < i2
}

func (g *Generator) indexOfPriority(priority string) int {
	for i, p := range g.monorailCfg.Priorities {
		if p.Priority == priority {
			return i
		}
	}
	// If we can't find the priority, treat it as one lower than
	// the lowest priority we know about.
	return len(g.monorailCfg.Priorities)
}

// isCompatibleWithVerified returns whether the impact of the current cluster
// is compatible with the issue having the given verified status, based on
// configured thresholds and hysteresis.
func (g *Generator) isCompatibleWithVerified(verified bool) bool {
	hysteresisPerc := g.monorailCfg.PriorityHysteresisPercent
	lowestPriority := g.monorailCfg.Priorities[len(g.monorailCfg.Priorities)-1]
	if verified {
		// The issue is verified. Only reopen if there is enough impact
		// to exceed the threshold with hysteresis.
		inflatedThreshold := bugs.InflateThreshold(lowestPriority.Threshold, hysteresisPerc)
		return !g.impact.MeetsThreshold(inflatedThreshold)
	} else {
		// The issue is not verified. Only close if the impact falls
		// below the threshold with hysteresis.
		deflatedThreshold := bugs.InflateThreshold(lowestPriority.Threshold, -hysteresisPerc)
		return g.impact.MeetsThreshold(deflatedThreshold)
	}
}

// isCompatibleWithPriority returns whether the impact of the current cluster
// is compatible with the issue having the given priority, based on
// configured thresholds and hysteresis.
func (g *Generator) isCompatibleWithPriority(issuePriority string) bool {
	index := g.indexOfPriority(issuePriority)
	if index >= len(g.monorailCfg.Priorities) {
		// Unknown priority in use. The priority should be updated to
		// one of the configured priorities.
		return false
	}
	hysteresisPerc := g.monorailCfg.PriorityHysteresisPercent
	lowestAllowedPriority := g.clusterPriorityWithInflatedThresholds(hysteresisPerc)
	highestAllowedPriority := g.clusterPriorityWithInflatedThresholds(-hysteresisPerc)

	// Check the cluster has a priority no less than lowest priority
	// and no greater than highest priority allowed by hysteresis.
	// Note that a lower priority index corresponds to a higher
	// priority (e.g. P0 <-> index 0, P1 <-> index 1, etc.)
	return g.indexOfPriority(lowestAllowedPriority) >= index &&
		index >= g.indexOfPriority(highestAllowedPriority)
}

// priorityDecreaseJustification outputs a human-readable justification
// explaining why bug priority was decreased (including to the point where
// a priority no longer applied, and the issue was marked as verified.)
//
// priorityIndex(s) are indices into the per-project priority list:
//   g.monorailCfg.Priorities
// The special index len(g.monorailCfg.Priorities) indicates a verified
// issue, i.e. an issue with so low priority that it does not deserve
// to be open.
//
// Example output:
// "- Presubmit Runs Failed (1-day) < 15, and
//  - Test Runs Failed (1-day) < 100"
func (g *Generator) priorityDecreaseJustification(oldPriorityIndex, newPriorityIndex int) string {
	if newPriorityIndex <= oldPriorityIndex {
		// Priority did not change or increased.
		return ""
	}

	// Priority decreased.
	// To justify the decrease, it is sufficient to explain why we could no
	// longer meet the criteria for the next-higher priority.

	hysteresisPerc := g.monorailCfg.PriorityHysteresisPercent

	// The next-higher priority level that we failed to meet.
	failedToMeetThreshold := g.monorailCfg.Priorities[newPriorityIndex-1].Threshold
	if newPriorityIndex == oldPriorityIndex+1 {
		// We only dropped one priority level. That means we failed to meet the
		// old threshold, even after applying hysteresis.
		failedToMeetThreshold = bugs.InflateThreshold(failedToMeetThreshold, -hysteresisPerc)
	}

	var message strings.Builder
	explanation := bugs.ExplainThresholdNotMet(failedToMeetThreshold)

	// As there may be multiple ways in which we could have met the
	// threshold for the next-higher priority (due to the OR-
	// disjunction of different metric thresholds), we must explain
	// we did not meet any of them.
	for i, exp := range explanation {
		message.WriteString(fmt.Sprintf("- %s (%v-day) < %v", exp.Metric, exp.TimescaleDays, exp.Threshold))
		if i < (len(explanation) - 1) {
			message.WriteString(", and")
		}
		message.WriteString("\n")
	}
	return message.String()
}

// priorityIncreaseJustification outputs a human-readable justification
// explaining why bug priority was increased (including for the case
// where a bug was re-opened.)
//
// priorityIndex(s) are indices into the per-project priority list:
//   g.monorailCfg.Priorities
// The special index len(g.monorailCfg.Priorities) indicates a verified
// issue, i.e. an issue with so low priority that it does not deserve
// to be open.
//
// Example output:
// "- Presubmit Runs Failed (1-day) >= 15"
func (g *Generator) priorityIncreaseJustification(oldPriorityIndex, newPriorityIndex int) string {
	if newPriorityIndex >= oldPriorityIndex {
		// Priority did not change or decreased.
		return ""
	}

	// Priority increased.
	// To justify the increase, we must show that we met the criteria for
	// each successively higher priority level.

	hysteresisPerc := g.monorailCfg.PriorityHysteresisPercent

	// Visit priorities in increasing priority order.
	var explanations []bugs.ThresholdExplanation
	for i := oldPriorityIndex - 1; i >= newPriorityIndex; i-- {
		metThreshold := g.monorailCfg.Priorities[i].Threshold
		if i == oldPriorityIndex-1 {
			// For the first priority step up, we must have also exceeded
			// hysteresis.
			metThreshold = bugs.InflateThreshold(metThreshold, hysteresisPerc)
		}

		// There may be multiple ways in which we could have met the
		// threshold for the next-higher priority (due to the OR-
		// disjunction of different metric thresholds). This obtains
		// just one of the ways in which we met it.
		explanations = append(explanations, g.impact.ExplainThresholdMet(metThreshold))
	}

	// Remove redundant explanations.
	// E.g. "Presubmit Runs Failed (1-day) >= 15"
	// and "Presubmit Runs Failed (1-day) >= 30" can be merged to just
	// "Presubmit Runs Failed (1-day) >= 30", because the latter
	// trivially implies the former.
	explanations = bugs.MergeThresholdMetExplanations(explanations)

	var message strings.Builder
	for i, exp := range explanations {
		message.WriteString(fmt.Sprintf("- %s (%v-day) >= %v", exp.Metric, exp.TimescaleDays, exp.Threshold))
		if i < (len(explanations) - 1) {
			message.WriteString(", and")
		}
		message.WriteString("\n")
	}
	return message.String()
}

// clusterPriority returns the desired priority of the bug, if no hysteresis
// is applied.
func (g *Generator) clusterPriority() string {
	return g.clusterPriorityWithInflatedThresholds(0)
}

// clusterPriority returns the desired priority of the bug, if thresholds
// are inflated or deflated with the given percentage.
//
// See bugs.InflateThreshold for the interpretation of inflationPercent.
func (g *Generator) clusterPriorityWithInflatedThresholds(inflationPercent int64) string {
	// Default to using the lowest priority.
	priority := g.monorailCfg.Priorities[len(g.monorailCfg.Priorities)-1]
	for i := len(g.monorailCfg.Priorities) - 2; i >= 0; i-- {
		p := g.monorailCfg.Priorities[i]
		adjustedThreshold := bugs.InflateThreshold(p.Threshold, inflationPercent)
		if !g.impact.MeetsThreshold(adjustedThreshold) {
			// A cluster cannot reach a higher priority unless it has
			// met the thresholds for all lower priorities.
			break
		}
		priority = p
	}
	return priority.Priority
}

// clusterResolved returns the desired state of whether the cluster has been
// verified, if no hysteresis has been applied.
func (g *Generator) clusterResolved() bool {
	lowestPriority := g.monorailCfg.Priorities[len(g.monorailCfg.Priorities)-1]
	return !g.impact.MeetsThreshold(lowestPriority.Threshold)
}

// sanitiseTitle removes tabs and line breaks from input, replacing them with
// spaces, and truncates the output to the given number of runes.
func sanitiseTitle(input string, maxLength int) string {
	// Replace blocks of whitespace, including new lines and tabs, with just a
	// single space.
	strippedInput := whitespaceRE.ReplaceAllString(input, " ")

	// Truncate to desired length.
	runes := []rune(strippedInput)
	if len(runes) > maxLength {
		return string(runes[0:maxLength-3]) + "..."
	}
	return strippedInput
}
