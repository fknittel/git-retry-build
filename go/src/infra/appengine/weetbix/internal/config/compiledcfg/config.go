// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package compiledcfg contains compiled versions of the Weetbix config.
// (E.g. Regular expressions are compiled for efficiency.)
package compiledcfg

import (
	"context"
	"time"

	"go.chromium.org/luci/common/data/caching/lru"
	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/server/caching"

	"infra/appengine/weetbix/internal/clustering/algorithms/testname/rules"
	"infra/appengine/weetbix/internal/config"
	configpb "infra/appengine/weetbix/internal/config/proto"
)

// TODO(crbug.com/1243174). Instrument the size of this cache so that we
// can monitor it.
var configCache = caching.RegisterLRUCache(0)

// NotExistsErr is returned if no matching configuration could be found
// for the specified project.
var NotExistsErr = errors.New("no config exists for the specified project")

// ProjectConfig is a compiled version of Weetbix project configuration.
type ProjectConfig struct {
	// Config is the raw, uncompiled, configuration.
	Config *configpb.ProjectConfig

	// TestNameRules are the set of rules to use to cluster test results
	// by test name.
	TestNameRules []rules.Evaluator

	// LastUpdated is the time the configuration was last updated.
	LastUpdated time.Time
}

// NewConfig compiles the given clustering configuration into a Config
// object.
func NewConfig(config *configpb.ProjectConfig) (*ProjectConfig, error) {
	rs := config.Clustering.GetTestNameRules()
	compiledRules := make([]rules.Evaluator, len(rs))
	for i, rule := range rs {
		eval, err := rules.Compile(rule)
		if err != nil {
			return nil, errors.Annotate(err, "compiling test name clustering rule").Err()
		}
		compiledRules[i] = eval
	}
	return &ProjectConfig{
		Config:        config,
		TestNameRules: compiledRules,
		LastUpdated:   config.LastUpdated.AsTime(),
	}, nil
}

// Project returns the clustering configuration for the given project,
// with a LastUpdated time of at least minimumVersion. If no particular
// minimum version is desired, pass time.Time{} to minimumVersion.
func Project(ctx context.Context, project string, minimumVersion time.Time) (*ProjectConfig, error) {
	cache := configCache.LRU(ctx)
	if cache == nil {
		// A fallback useful in unit tests that may not have the process cache
		// available. Production environments usually have the cache installed
		// by the framework code that initializes the root context.
		projectCfg, err := config.ProjectWithMinimumVersion(ctx, project, minimumVersion)
		if err != nil {
			if err == config.NotExistsErr {
				return nil, NotExistsErr
			}
			return nil, err
		}
		config, err := NewConfig(projectCfg)
		if err != nil {
			return nil, err
		}
		return config, nil
	} else {
		var err error
		val, _ := cache.Mutate(ctx, project, func(it *lru.Item) *lru.Item {
			var projectCfg *configpb.ProjectConfig
			// Fetch the latest configuration for the given project, with
			// the specified minimum version.
			projectCfg, err = config.ProjectWithMinimumVersion(ctx, project, minimumVersion)
			if err != nil {
				// Delete cached value.
				return nil
			}

			if it != nil {
				cfg := it.Value.(*ProjectConfig)
				if cfg.LastUpdated.Equal(projectCfg.LastUpdated.AsTime()) {
					// Cached value is already up to date.
					return it
				}
			}
			var config *ProjectConfig
			config, err = NewConfig(projectCfg)
			if err != nil {
				// Delete cached value.
				return nil
			}
			return &lru.Item{
				Value: config,
				Exp:   0, // No expiry.
			}
		})
		if err != nil {
			if err == config.NotExistsErr {
				return nil, NotExistsErr
			}
			return nil, errors.Annotate(err, "obtain compiled configuration").Err()
		}
		cfg := val.(*ProjectConfig)
		return cfg, nil
	}
}
