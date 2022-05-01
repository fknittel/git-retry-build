// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cache

import (
	"sort"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go.chromium.org/luci/common/clock/testclock"
	"go.chromium.org/luci/server/caching"

	"infra/appengine/weetbix/internal/bugs"
	"infra/appengine/weetbix/internal/clustering/rules"
	"infra/appengine/weetbix/internal/testutil"
)

var cache = caching.RegisterLRUCache(50)

func TestRulesCache(t *testing.T) {
	Convey(`With Spanner Test Database`, t, func() {
		ctx := testutil.SpannerTestContext(t)
		ctx, tc := testclock.UseTime(ctx, testclock.TestRecentTimeUTC)
		ctx = caching.WithEmptyProcessCache(ctx)

		rc := NewRulesCache(cache)
		rules.SetRulesForTesting(ctx, nil)

		test := func(minimumPredicatesVerison time.Time, expectedRules []*rules.FailureAssociationRule, expectedVersion rules.Version) {
			// Tests the content of the cache is as expected.
			ruleset, err := rc.Ruleset(ctx, "myproject", minimumPredicatesVerison)
			So(err, ShouldBeNil)
			So(ruleset.Version, ShouldResemble, expectedVersion)

			activeRules := 0
			for _, e := range expectedRules {
				if e.IsActive {
					activeRules++
				}
			}
			So(len(ruleset.ActiveRulesSorted), ShouldEqual, activeRules)
			So(len(ruleset.ActiveRulesByID), ShouldEqual, activeRules)

			sortedExpectedRules := sortRulesByPredicateLastUpdated(expectedRules)

			actualRuleIndex := 0
			for _, e := range sortedExpectedRules {
				if e.IsActive {
					a := ruleset.ActiveRulesSorted[actualRuleIndex]
					So(a.Rule, ShouldResemble, *e)
					// Technically (*lang.Expr).String() may not get us
					// back the original rule if RuleDefinition didn't use
					// normalised formatting. But for this test, we use
					// normalised formatting, so that is not an issue.
					So(a.Expr, ShouldNotBeNil)
					So(a.Expr.String(), ShouldEqual, e.RuleDefinition)
					actualRuleIndex++

					a2, ok := ruleset.ActiveRulesByID[a.Rule.RuleID]
					So(ok, ShouldBeTrue)
					So(a2.Rule, ShouldResemble, *e)
				}
			}
			So(len(ruleset.ActiveRulesWithPredicateUpdatedSince(rules.StartingEpoch)), ShouldEqual, activeRules)
			So(len(ruleset.ActiveRulesWithPredicateUpdatedSince(time.Date(2100, time.January, 1, 1, 0, 0, 0, time.UTC))), ShouldEqual, 0)
		}

		Convey(`Initially Empty`, func() {
			err := rules.SetRulesForTesting(ctx, nil)
			So(err, ShouldBeNil)
			test(rules.StartingEpoch, nil, rules.StartingVersion)

			Convey(`Then Empty`, func() {
				// Test cache.
				test(rules.StartingEpoch, nil, rules.StartingVersion)

				tc.Add(refreshInterval)

				test(rules.StartingEpoch, nil, rules.StartingVersion)
				test(rules.StartingEpoch, nil, rules.StartingVersion)
			})
			Convey(`Then Non-Empty`, func() {
				// Spanner commit timestamps are in microsecond
				// (not nanosecond) granularity, and some Spanner timestamp
				// operators truncates to microseconds. For this
				// reason, we use microsecond resolution timestamps
				// when testing.
				reference := time.Date(2020, 1, 2, 3, 4, 5, 6000, time.UTC)

				rs := []*rules.FailureAssociationRule{
					rules.NewRule(100).
						WithLastUpdated(reference.Add(-1 * time.Hour)).
						WithPredicateLastUpdated(reference.Add(-2 * time.Hour)).
						Build(),
					rules.NewRule(101).WithActive(false).
						WithLastUpdated(reference.Add(1 * time.Hour)).
						WithPredicateLastUpdated(reference).
						Build(),
				}
				err := rules.SetRulesForTesting(ctx, rs)
				So(err, ShouldBeNil)

				expectedRulesVersion := rules.Version{
					Total:      reference.Add(1 * time.Hour),
					Predicates: reference,
				}

				Convey(`By Strong Read`, func() {
					test(StrongRead, rs, expectedRulesVersion)
					test(StrongRead, rs, expectedRulesVersion)
				})
				Convey(`By Requesting Version`, func() {
					test(expectedRulesVersion.Predicates, rs, expectedRulesVersion)
				})
				Convey(`By Cache Expiry`, func() {
					// Test cache is working and still returning the old value.
					tc.Add(refreshInterval / 2)
					test(rules.StartingEpoch, nil, rules.StartingVersion)

					tc.Add(refreshInterval)

					test(rules.StartingEpoch, rs, expectedRulesVersion)
					test(rules.StartingEpoch, rs, expectedRulesVersion)
				})
			})
		})
		Convey(`Initially Non-Empty`, func() {
			reference := time.Date(2021, 1, 2, 3, 4, 5, 6000, time.UTC)

			ruleOne := rules.NewRule(100).
				WithLastUpdated(reference.Add(-2 * time.Hour)).
				WithPredicateLastUpdated(reference.Add(-3 * time.Hour))
			ruleTwo := rules.NewRule(101).
				WithLastUpdated(reference.Add(-2 * time.Hour)).
				WithPredicateLastUpdated(reference.Add(-3 * time.Hour))
			ruleThree := rules.NewRule(102).WithActive(false).
				WithLastUpdated(reference).
				WithPredicateLastUpdated(reference.Add(-1 * time.Hour))

			rs := []*rules.FailureAssociationRule{
				ruleOne.Build(),
				ruleTwo.Build(),
				ruleThree.Build(),
			}
			err := rules.SetRulesForTesting(ctx, rs)
			So(err, ShouldBeNil)

			expectedRulesVersion := rules.Version{
				Total:      reference,
				Predicates: reference.Add(-1 * time.Hour),
			}
			test(rules.StartingEpoch, rs, expectedRulesVersion)

			Convey(`Then Empty`, func() {
				// Mark all rules inactive.
				newRules := []*rules.FailureAssociationRule{
					ruleOne.WithActive(false).
						WithLastUpdated(reference.Add(4 * time.Hour)).
						WithPredicateLastUpdated(reference.Add(3 * time.Hour)).
						Build(),
					ruleTwo.WithActive(false).
						WithLastUpdated(reference.Add(2 * time.Hour)).
						WithPredicateLastUpdated(reference.Add(1 * time.Hour)).
						Build(),
					ruleThree.WithActive(false).
						WithLastUpdated(reference.Add(2 * time.Hour)).
						WithPredicateLastUpdated(reference.Add(1 * time.Hour)).
						Build(),
				}
				err := rules.SetRulesForTesting(ctx, newRules)
				So(err, ShouldBeNil)

				oldRulesVersion := expectedRulesVersion
				expectedRulesVersion := rules.Version{
					Total:      reference.Add(4 * time.Hour),
					Predicates: reference.Add(3 * time.Hour),
				}

				Convey(`By Strong Read`, func() {
					test(StrongRead, newRules, expectedRulesVersion)
					test(StrongRead, newRules, expectedRulesVersion)
				})
				Convey(`By Requesting Version`, func() {
					test(expectedRulesVersion.Predicates, newRules, expectedRulesVersion)
				})
				Convey(`By Cache Expiry`, func() {
					// Test cache is working and still returning the old value.
					tc.Add(refreshInterval / 2)
					test(rules.StartingEpoch, rs, oldRulesVersion)

					tc.Add(refreshInterval)

					test(rules.StartingEpoch, newRules, expectedRulesVersion)
					test(rules.StartingEpoch, newRules, expectedRulesVersion)
				})
			})
			Convey(`Then Non-Empty`, func() {
				newRules := []*rules.FailureAssociationRule{
					// Mark an existing rule inactive.
					ruleOne.WithActive(false).
						WithLastUpdated(reference.Add(time.Hour)).
						WithPredicateLastUpdated(reference.Add(time.Hour)).
						Build(),
					// Make a non-predicate change on an active rule.
					ruleTwo.
						WithBug(bugs.BugID{System: "monorail", ID: "project/123"}).
						WithLastUpdated(reference.Add(time.Hour)).
						Build(),
					// Make an existing rule active.
					ruleThree.WithActive(true).
						WithLastUpdated(reference.Add(time.Hour)).
						WithPredicateLastUpdated(reference.Add(time.Hour)).
						Build(),
					// Add a new active rule.
					rules.NewRule(103).
						WithPredicateLastUpdated(reference.Add(time.Hour)).
						WithLastUpdated(reference.Add(time.Hour)).
						Build(),
					// Add a new inactive rule.
					rules.NewRule(104).WithActive(false).
						WithPredicateLastUpdated(reference.Add(2 * time.Hour)).
						WithLastUpdated(reference.Add(3 * time.Hour)).
						Build(),
				}
				err := rules.SetRulesForTesting(ctx, newRules)
				So(err, ShouldBeNil)

				oldRulesVersion := expectedRulesVersion
				expectedRulesVersion := rules.Version{
					Total:      reference.Add(3 * time.Hour),
					Predicates: reference.Add(2 * time.Hour),
				}

				Convey(`By Strong Read`, func() {
					test(StrongRead, newRules, expectedRulesVersion)
					test(StrongRead, newRules, expectedRulesVersion)
				})
				Convey(`By Forced Eviction`, func() {
					test(expectedRulesVersion.Predicates, newRules, expectedRulesVersion)
				})
				Convey(`By Cache Expiry`, func() {
					// Test cache is working and still returning the old value.
					tc.Add(refreshInterval / 2)
					test(rules.StartingEpoch, rs, oldRulesVersion)

					tc.Add(refreshInterval)

					test(rules.StartingEpoch, newRules, expectedRulesVersion)
					test(rules.StartingEpoch, newRules, expectedRulesVersion)
				})
			})
		})
	})
}

func sortRulesByPredicateLastUpdated(rs []*rules.FailureAssociationRule) []*rules.FailureAssociationRule {
	result := make([]*rules.FailureAssociationRule, len(rs))
	copy(result, rs)
	sort.Slice(result, func(i, j int) bool {
		if result[i].PredicateLastUpdated.Equal(result[j].PredicateLastUpdated) {
			return result[i].RuleID < result[j].RuleID
		}
		return result[i].PredicateLastUpdated.After(result[j].PredicateLastUpdated)
	})
	return result
}
