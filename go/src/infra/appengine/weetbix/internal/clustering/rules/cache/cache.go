// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cache

import (
	"context"
	"fmt"
	"time"

	"go.chromium.org/luci/common/clock"
	"go.chromium.org/luci/common/data/caching/lru"
	"go.chromium.org/luci/server/caching"
)

// refreshInterval controls how often rulesets are refreshed.
const refreshInterval = time.Minute

// RulesCache is an in-process cache of failure association rules used
// by LUCI projects.
type RulesCache struct {
	cache caching.LRUHandle
}

// NewRulesCache initialises a new RulesCache.
func NewRulesCache(c caching.LRUHandle) *RulesCache {
	return &RulesCache{
		cache: c,
	}
}

// Ruleset obtains the Ruleset for a particular project from the cache, or if
// it does not exist, retrieves it from Spanner. MinimumPredicatesVersion
// specifies the minimum version of rule predicates that must be incorporated
// in the given Ruleset. If no particular version is desired, pass
// rules.StartingEpoch.
func (c *RulesCache) Ruleset(ctx context.Context, project string, minimumPredicatesVersion time.Time) (*Ruleset, error) {
	var err error
	now := clock.Now(ctx)
	value, _ := c.cache.LRU(ctx).Mutate(ctx, project, func(it *lru.Item) *lru.Item {
		var ruleset *Ruleset
		if it != nil {
			ruleset = it.Value.(*Ruleset)
			if ruleset.LastRefresh.Add(refreshInterval).After(now) && !ruleset.Version.Predicates.Before(minimumPredicatesVersion) {
				// The ruleset is up-to-date. Do not mutate it further.
				return it
			}
		} else {
			ruleset = newEmptyRuleset(project)
		}
		ruleset, err = ruleset.refresh(ctx)
		if err != nil {
			// Issue refreshing ruleset. Keep the cached value (if any) for now.
			return it
		}
		return &lru.Item{
			Value: ruleset,
			Exp:   0, // Never.
		}
	})
	if err != nil {
		return nil, err
	}
	ruleset := value.(*Ruleset)
	if ruleset.Version.Predicates.Before(minimumPredicatesVersion) {
		return nil, fmt.Errorf("could not obtain ruleset of requested minimum predicate version (%v)", minimumPredicatesVersion)
	}
	return ruleset, nil
}
