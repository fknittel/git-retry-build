// Copyright 2020 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"bufio"
	"context"
	"encoding/json"
	"math"
	"os"
	"path/filepath"

	"cloud.google.com/go/bigquery"
	"github.com/maruel/subcommands"
	"golang.org/x/sync/errgroup"

	"go.chromium.org/luci/auth"
	"go.chromium.org/luci/common/cli"
	"go.chromium.org/luci/common/data/text"
	"go.chromium.org/luci/common/errors"

	"infra/rts/filegraph/git"
	"infra/rts/presubmit/eval"
)

func cmdCreateModel(authOpt *auth.Options) *subcommands.Command {
	return &subcommands.Command{
		UsageLine: `create-model -model-dir <path>`,
		ShortDesc: "create a model to be used by select subcommand",
		LongDesc:  "Create a model to be used by select subcommand",
		CommandRun: func() subcommands.CommandRun {
			r := &createModelRun{authOpt: authOpt}
			r.Flags.StringVar(&r.modelDir, "model-dir", "", text.Doc(`
				Path to the directory where to write the model files.
				The directory will be created if it does not exist.
			`))

			r.Flags.StringVar(&r.checkout, "checkout", "", "Path to a src.git checkout")
			r.Flags.IntVar(&r.loadOptions.MaxCommitSize, "fg-max-commit-size", 100, text.Doc(`
				Maximum number of files touched by a commit.
				Commits that exceed this limit are ignored.
				The rationale is that large commits provide a weak signal of file
				relatedness and are expensive to process, O(N^2).
			`))

			r.ev.RegisterFlags(&r.Flags)
			return r
		},
	}
}

type createModelRun struct {
	baseCommandRun
	modelDir string

	checkout    string
	loadOptions git.LoadOptions
	fg          *git.Graph

	ev eval.Eval

	authOpt  *auth.Options
	bqClient *bigquery.Client
}

func (r *createModelRun) validateFlags() error {
	if err := r.ev.ValidateFlags(); err != nil {
		return err
	}
	switch {
	case r.modelDir == "":
		return errors.New("-model-dir is required")
	case r.checkout == "":
		return errors.New("-checkout is required")
	default:
		return nil
	}
}

func (r *createModelRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	ctx := cli.GetContext(a, r, env)
	if len(args) != 0 {
		return r.done(errors.New("unexpected positional arguments"))
	}
	if err := r.validateFlags(); err != nil {
		return r.done(err)
	}

	var err error
	if r.bqClient, err = newBQClient(ctx, auth.NewAuthenticator(ctx, auth.InteractiveLogin, *r.authOpt)); err != nil {
		return r.done(errors.Annotate(err, "failed to create BigQuery client").Err())
	}

	return r.done(r.writeModel(ctx, r.modelDir))
}

// writeModel writes the model files to the directory.
func (r *createModelRun) writeModel(ctx context.Context, dir string) error {
	// Ensure model dir exists.
	if err := os.MkdirAll(dir, 0777); err != nil {
		return errors.Annotate(err, "failed to create model dir at %q", dir).Err()
	}

	eg, ctx := errgroup.WithContext(ctx)
	defer eg.Wait()

	eg.Go(func() error {
		err := r.writeFileGraphModel(ctx, filepath.Join(dir, "git-file-graph"))
		return errors.Annotate(err, "failed to write file graph model").Err()
	})

	eg.Go(func() error {
		err := r.writeTestFileSet(ctx, filepath.Join(dir, "test-files.jsonl"))
		return errors.Annotate(err, "failed to write test file set").Err()
	})

	return eg.Wait()
}

// writeFileGraphModel writes the file graph model to the model dir.
func (r *createModelRun) writeFileGraphModel(ctx context.Context, dir string) error {
	var err error
	if r.fg, err = git.Load(ctx, r.checkout, r.loadOptions); err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)
	defer eg.Wait()

	eg.Go(func() error {
		err := r.writeFileGraph(ctx, filepath.Join(dir, "graph.fg"))
		return errors.Annotate(err, "failed to write file graph").Err()
	})

	eg.Go(func() error {
		err := r.writeEvalResults(ctx, filepath.Join(dir, "eval-results.json"))
		return errors.Annotate(err, "failed to write thresholds").Err()
	})

	return eg.Wait()
}

// writeFileGraph writes the graph file.
func (r *createModelRun) writeFileGraph(ctx context.Context, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	bufW := bufio.NewWriter(f)
	if err := r.fg.Write(bufW); err != nil {
		return err
	}
	return bufW.Flush()
}

// writeEvalResults evaluates the selection strategy, prints results and writes
// them to the file.
func (r *createModelRun) writeEvalResults(ctx context.Context, fileName string) error {
	r.ev.Strategy = r.selectTests
	res, err := r.ev.Run(ctx)
	if err != nil {
		return err
	}

	res.Print(os.Stdout, 0.97)

	// JSON does not support infinity. Replace it with a large number.
	for _, t := range res.Thresholds {
		if math.IsInf(t.Value.Distance, 1) {
			t.Value.Distance = 1e10
		}
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(res)
}

// writeTestFileSet writes the test file set in Chromium to the file.
// It skips tests that match neverSkipTestFileRegexp.
//
// The file format is JSON Lines of TestFile protobufs.
func (r *createModelRun) writeTestFileSet(ctx context.Context, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	bufW := bufio.NewWriter(f)

	if err := writeTestFiles(ctx, r.bqClient, bufW); err != nil {
		return err
	}

	if err := bufW.Flush(); err != nil {
		return err
	}
	return f.Close()
}
