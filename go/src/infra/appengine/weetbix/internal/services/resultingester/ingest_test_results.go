// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package resultingester

import (
	"context"
	"fmt"
	"time"

	bbpb "go.chromium.org/luci/buildbucket/proto"
	"go.chromium.org/luci/common/clock"
	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/common/logging"
	"go.chromium.org/luci/common/retry/transient"
	"go.chromium.org/luci/common/tsmon/field"
	"go.chromium.org/luci/common/tsmon/metric"
	rdbbutil "go.chromium.org/luci/resultdb/pbutil"
	rdbpb "go.chromium.org/luci/resultdb/proto/v1"
	"go.chromium.org/luci/server"
	"go.chromium.org/luci/server/tq"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"infra/appengine/weetbix/internal/analysis"
	"infra/appengine/weetbix/internal/analysis/clusteredfailures"
	"infra/appengine/weetbix/internal/buildbucket"
	"infra/appengine/weetbix/internal/clustering/chunkstore"
	"infra/appengine/weetbix/internal/clustering/ingestion"
	"infra/appengine/weetbix/internal/config"
	"infra/appengine/weetbix/internal/resultdb"
	"infra/appengine/weetbix/internal/services/resultcollector"
	"infra/appengine/weetbix/internal/tasks/taskspb"
	"infra/appengine/weetbix/utils"
)

const (
	resultIngestionTaskClass = "result-ingestion"
	resultIngestionQueue     = "result-ingestion"

	// ingestionEarliest is the oldest data that may be ingested by Weetbix.
	// This is an offset relative to the current time, and should be kept
	// in sync with the data retention period in Spanner and BigQuery.
	ingestionEarliest = -90 * 24 * time.Hour

	// ingestionLatest is the newest data that may be ingested by Weetbix.
	// This is an offset relative to the current time. It is designed to
	// allow for clock drift.
	ingestionLatest = 24 * time.Hour
)

// maxResultDBPages is the maximum number of pages of results to ingest from
// ResultDB, per build. The page size is 1000 results.
const maxResultDBPages = 10

var (
	taskCounter = metric.NewCounter(
		"weetbix/ingestion/task_completion",
		"The number of completed Weetbix ingestion tasks, by build project and outcome.",
		nil,
		// The LUCI Project.
		field.String("project"),
		// "success", "failed_validation",
		// "ignored_no_bb_access", "ignored_no_project_config",
		// "ignored_cq_dry_run", "ignored_no_invocation".
		field.String("outcome"))
)

// Options configures test result ingestion.
type Options struct {
}

type resultIngester struct {
	clustering *ingestion.Ingester
}

var resultIngestion = tq.RegisterTaskClass(tq.TaskClass{
	ID:        resultIngestionTaskClass,
	Prototype: &taskspb.IngestTestResults{},
	Queue:     resultIngestionQueue,
	Kind:      tq.Transactional,
})

// RegisterTaskHandler registers the handler for result ingestion tasks.
func RegisterTaskHandler(srv *server.Server) error {
	ctx := srv.Context
	cfg, err := config.Get(ctx)
	if err != nil {
		return err
	}
	chunkStore, err := chunkstore.NewClient(ctx, cfg.ChunkGcsBucket)
	if err != nil {
		return err
	}
	srv.RegisterCleanup(func(ctx context.Context) {
		chunkStore.Close()
	})
	cf := clusteredfailures.NewClient(srv.Options.CloudProject)
	analysis := analysis.NewClusteringHandler(cf)
	ri := &resultIngester{
		clustering: ingestion.New(chunkStore, analysis),
	}
	handler := func(ctx context.Context, payload proto.Message) error {
		task := payload.(*taskspb.IngestTestResults)
		return ri.ingestTestResults(ctx, task)
	}
	resultIngestion.AttachHandler(handler)
	return nil
}

// Schedule enqueues a task to ingest test results from a build.
func Schedule(ctx context.Context, task *taskspb.IngestTestResults) {
	tq.MustAddTask(ctx, &tq.Task{
		Title:   fmt.Sprintf("%s-%s-%d", task.Build.Project, task.Build.Host, task.Build.Id),
		Payload: task,
	})
}

func (i *resultIngester) ingestTestResults(ctx context.Context, payload *taskspb.IngestTestResults) error {
	if err := validateRequest(ctx, payload); err != nil {
		project := "(unknown)"
		if payload.Build != nil && payload.Build.Project != "" {
			project = payload.Build.Project
		}
		taskCounter.Add(ctx, 1, project, "failed_validation")
		return tq.Fatal.Apply(err)
	}

	if _, err := config.Project(ctx, payload.Build.Project); err != nil {
		if err == config.NotExistsErr {
			// Project not configured in Weetbix, ignore it.
			taskCounter.Add(ctx, 1, payload.Build.Project, "ignored_no_project_config")
			return nil
		} else {
			// Transient error.
			return transient.Tag.Apply(errors.Annotate(err, "get project config").Err())
		}
	}

	// Buildbucket build only has builder, infra.resultdb, status populated.
	b, err := retrieveBuild(ctx, payload)
	code := status.Code(err)
	if code == codes.NotFound {
		// Build not found, end the task gracefully.
		logging.Warningf(ctx, "Buildbucket build %s/%d for project %s not found (or Weetbix does not have access to read it).",
			payload.Build.Host, payload.Build.Id, payload.Build.Project)
		taskCounter.Add(ctx, 1, payload.Build.Project, "ignored_no_bb_access")
		return nil
	}
	if err != nil {
		return err
	}

	if b.Infra.GetResultdb().GetInvocation() == "" {
		// Build does not have a ResultDB invocation to ingest.
		logging.Debugf(ctx, "Skipping ingestion of build %s-%d because it has no ResultDB invocation.",
			payload.Build.Host, payload.Build.Id)
		taskCounter.Add(ctx, 1, payload.Build.Project, "ignored_no_invocation")
		return nil
	}
	if payload.PresubmitRun != nil && payload.PresubmitRun.Mode != "FULL_RUN" {
		// CQ Dry Runs currently add a lot of noise to the analysis, which
		// the analysis is not yet set up to deal with. Skip for now.
		logging.Debugf(ctx, "Skipping ingestion of build %s-%d because it was a CQ Dry Run.",
			payload.Build.Host, payload.Build.Id)
		taskCounter.Add(ctx, 1, payload.Build.Project, "ignored_cq_dry_run")
		return nil
	}

	rdbHost := b.Infra.Resultdb.Hostname
	invName := b.Infra.Resultdb.Invocation
	builder := b.Builder.Builder
	rc, err := resultdb.NewClient(ctx, rdbHost)
	if err != nil {
		return transient.Tag.Apply(err)
	}
	inv, err := rc.GetInvocation(ctx, invName)
	code = status.Code(err)
	if code == codes.NotFound {
		// Invocation not found, end the task gracefully.
		logging.Warningf(ctx, "Invocation %s for project %s not found (or Weetbix does not have access to read it).",
			invName, payload.Build.Project)
		taskCounter.Add(ctx, 1, payload.Build.Project, "ignored_no_resultdb_access")
		return nil
	}
	if err != nil {
		return transient.Tag.Apply(err)
	}

	project, _ := utils.SplitRealm(inv.Realm)
	if project == "" {
		return fmt.Errorf("invocation has invalid realm: %q", inv.Realm)
	}

	realmCfg, err := config.Realm(ctx, inv.Realm)
	if err != nil && err != config.RealmNotExistsErr {
		return transient.Tag.Apply(err)
	}
	ingestForTestVariantAnalysis := realmCfg != nil &&
		shouldIngestForTestVariants(realmCfg, payload)

	// Setup clustering ingestion.
	invID, err := rdbbutil.ParseInvocationName(invName)
	if err != nil {
		// This should never happen.
		return transient.Tag.Apply(err)
	}
	opts := ingestion.Options{
		Project:       project,
		InvocationID:  invID,
		PartitionTime: payload.PartitionTime.AsTime(),
		Realm:         inv.Realm,
		// In case of Success, Cancellation, or Infra Failure automatically
		// exonerate failures of tests which were invocation-blocking,
		// even if the recipe did not upload an exoneration to ResultDB.
		// The build status implies the test result could not have been
		// responsible for causing the build (or consequently, the CQ run)
		// to fail.
		ImplicitlyExonerateBlockingFailures: b.Status != bbpb.Status_FAILURE,
	}
	if payload.PresubmitRun != nil {
		opts.PresubmitRunID = payload.PresubmitRun.PresubmitRunId
		opts.PresubmitRunOwner = payload.PresubmitRun.Owner
		opts.PresubmitRunCls = payload.PresubmitRun.Cls
		if !payload.PresubmitRun.Critical {
			// CQ did not consider the build critical.
			opts.ImplicitlyExonerateBlockingFailures = true
		}
		if payload.PresubmitRun.Critical && b.Status == bbpb.Status_FAILURE &&
			payload.PresubmitRun.PresubmitRunSucceeded {
			logging.Warningf(ctx, "Inconsistent data from LUCI CV: build %v/%v was critical to presubmit run %v/%v and failed, but presubmit run did not fail.",
				payload.Build.Host, payload.Build.Id, payload.PresubmitRun.PresubmitRunId.System, payload.PresubmitRun.PresubmitRunId.Id)
		}
	}
	clusterIngestion := i.clustering.Open(opts)

	// Query test variants from ResultDB and save/update the corresponding
	// AnalyzedTestVariant rows.
	// We read test variants from ResultDB in pages, and the func will be called
	// once per page of test variants.
	f := func(tvs []*rdbpb.TestVariant) error {
		if ingestForTestVariantAnalysis {
			if err := createOrUpdateAnalyzedTestVariants(ctx, inv.Realm, builder, tvs); err != nil {
				return errors.Annotate(err, "ingesting for test variant analysis").Err()
			}
		}
		// Clustering ingestion is designed to behave gracefully in case of
		// a task retry. Given the same options and same test variants (in
		// the same order), the IDs and content of the chunks it writes is
		// designed to be stable. If chunks already exist, it will skip them.
		if err := clusterIngestion.Put(ctx, tvs); err != nil {
			return errors.Annotate(err, "ingesting for clustering").Err()
		}
		return nil
	}
	req := &rdbpb.QueryTestVariantsRequest{
		Invocations: []string{invName},
		Predicate: &rdbpb.TestVariantPredicate{
			Status: rdbpb.TestVariantStatus_UNEXPECTED_MASK,
		},
		PageSize: 1000,
	}
	err = rc.QueryTestVariants(ctx, req, f, maxResultDBPages)
	if err != nil {
		return transient.Tag.Apply(err)
	}
	if err := clusterIngestion.Flush(ctx); err != nil {
		return errors.Annotate(err, "ingesting for clustering").Err()
	}

	if ingestForTestVariantAnalysis {
		isPreSubmit := payload.PresubmitRun != nil
		contributedToCLSubmission := payload.PresubmitRun != nil && payload.PresubmitRun.PresubmitRunSucceeded
		if err = resultcollector.Schedule(ctx, inv, rdbHost, b.Builder.Builder, isPreSubmit, contributedToCLSubmission); err != nil {
			return transient.Tag.Apply(err)
		}
	}

	taskCounter.Add(ctx, 1, project, "success")
	return nil
}

func validateRequest(ctx context.Context, payload *taskspb.IngestTestResults) error {
	if !payload.PartitionTime.IsValid() {
		return errors.New("partition time must be specified and valid")
	}
	t := payload.PartitionTime.AsTime()
	now := clock.Now(ctx)
	if t.Before(now.Add(ingestionEarliest)) {
		return fmt.Errorf("partition time (%v) is too long ago", t)
	} else if t.After(now.Add(ingestionLatest)) {
		return fmt.Errorf("partition time (%v) is too far in the future", t)
	}
	if payload.Build == nil {
		return errors.New("build must be specified")
	}
	return nil
}

func retrieveBuild(ctx context.Context, payload *taskspb.IngestTestResults) (*bbpb.Build, error) {
	bbHost := payload.Build.Host
	id := payload.Build.Id
	bc, err := buildbucket.NewClient(ctx, bbHost)
	if err != nil {
		return nil, err
	}
	request := &bbpb.GetBuildRequest{
		Id: id,
		Mask: &bbpb.BuildMask{
			Fields: &field_mask.FieldMask{
				Paths: []string{"builder", "infra.resultdb", "status"},
			},
		},
	}
	b, err := bc.GetBuild(ctx, request)
	switch {
	case err != nil:
		return nil, err
	}
	return b, nil
}
