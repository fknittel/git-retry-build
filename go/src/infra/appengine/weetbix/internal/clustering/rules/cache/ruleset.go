package cache

import (
	"context"
	"sort"
	"time"

	"go.chromium.org/luci/common/clock"
	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/common/trace"
	"go.chromium.org/luci/server/span"

	"infra/appengine/weetbix/internal/clustering/rules"
	"infra/appengine/weetbix/internal/clustering/rules/lang"
)

// CachedRule represents a "compiled" version of a failure
// association rule.
// It should be treated as immutable, and is therefore safe to
// share across multiple threads.
type CachedRule struct {
	// The failure association rule.
	Rule rules.FailureAssociationRule
	// The parsed and compiled failure association rule.
	Expr *lang.Expr
}

// NewCachedRule initialises a new CachedRule from the given failure
// association rule.
func NewCachedRule(rule *rules.FailureAssociationRule) (*CachedRule, error) {
	expr, err := lang.Parse(rule.RuleDefinition)
	if err != nil {
		return nil, err
	}
	return &CachedRule{
		Rule: *rule,
		Expr: expr,
	}, nil
}

// Ruleset represents a version of the set of failure
// association rules in use by a LUCI Project.
// It should be treated as immutable, and therefore safe to share
// across multiple threads.
type Ruleset struct {
	// The LUCI Project.
	Project string
	// ActiveRulesSorted is the set of active failure association rules
	// (should be used by Weetbix for matching), sorted in descending
	// PredicateLastUpdated time order.
	ActiveRulesSorted []*CachedRule
	// ActiveRulesByID stores active failure association
	// rules by their Rule ID.
	ActiveRulesByID map[string]*CachedRule
	// Version versions the contents of the Ruleset. These timestamps only
	// change if a rule is modified.
	Version rules.Version
	// LastRefresh contains the monotonic clock reading when the last ruleset
	// refresh was initiated. The refresh is guaranteed to contain all rules
	// changes made prior to this timestamp.
	LastRefresh time.Time
}

// ActiveRulesWithPredicateUpdatedSince returns the set of rules that are
// active and whose predicates have been updated since (but not including)
// the given time.
// Rules which have been made inactive since the given time will NOT be
// returned. To check if a previous rule has been made inactive, consider
// using IsRuleActive instead.
// The returned slice must not be mutated.
func (r *Ruleset) ActiveRulesWithPredicateUpdatedSince(t time.Time) []*CachedRule {
	// Use the property that ActiveRules is sorted by descending
	// LastUpdated time.
	for i, rule := range r.ActiveRulesSorted {
		if !rule.Rule.PredicateLastUpdated.After(t) {
			// This is the first rule that has not been updated since time t.
			// Return all rules up to (but not including) this rule.
			return r.ActiveRulesSorted[:i]
		}
	}
	return r.ActiveRulesSorted
}

// Returns whether the given ruleID is an active rule.
func (r *Ruleset) IsRuleActive(ruleID string) bool {
	_, ok := r.ActiveRulesByID[ruleID]
	return ok
}

// newEmptyRuleset initialises a new empty ruleset.
// This initial ruleset is invalid and must be refreshed before use.
func newEmptyRuleset(project string) *Ruleset {
	return &Ruleset{
		Project:           project,
		ActiveRulesSorted: nil,
		ActiveRulesByID:   make(map[string]*CachedRule),
		// The zero predicate last updated time is not valid and will be
		// rejected by clustering state validation if we ever try to save
		// it to Spanner as a chunk's RulesVersion.
		Version:     rules.Version{},
		LastRefresh: time.Time{},
	}
}

// NewRuleset creates a new ruleset with the given project,
// active rules, rules last updated and last refresh time.
func NewRuleset(project string, activeRules []*CachedRule, version rules.Version, lastRefresh time.Time) *Ruleset {
	return &Ruleset{
		Project:           project,
		ActiveRulesSorted: sortByDescendingPredicateLastUpdated(activeRules),
		ActiveRulesByID:   rulesByID(activeRules),
		Version:           version,
		LastRefresh:       lastRefresh,
	}
}

// refresh updates the ruleset. To ensure existing users of the rulset
// do not observe changes while they are using it, a new copy is returned.
func (r *Ruleset) refresh(ctx context.Context) (ruleset *Ruleset, err error) {
	// Under our design assumption of 10,000 active rules per project,
	// pulling and compiling all rules could take a meaningful amount
	// of time (@ 1KB per rule, = ~10MB).
	ctx, s := trace.StartSpan(ctx, "infra/appengine/weetbix/internal/clustering/rules/cache.Refresh")
	s.Attribute("project", r.Project)
	defer func() { s.End(err) }()

	// Use clock reading before refresh. The refresh is guaranteed
	// to contain all rule changes committed to Spanner prior to
	// this timestamp.
	lastRefresh := clock.Now(ctx)

	txn, cancel := span.ReadOnlyTransaction(ctx)
	defer cancel()

	var activeRules []*CachedRule
	if r.Version == (rules.Version{}) {
		// On the first refresh, query all active rules.
		ruleRows, err := rules.ReadActive(txn, r.Project)
		if err != nil {
			return nil, err
		}
		activeRules, err = cachedRulesFromFullRead(ruleRows)
		if err != nil {
			return nil, err
		}
	} else {
		// On subsequent refreshes, query just the differences.
		delta, err := rules.ReadDelta(txn, r.Project, r.Version.Total)
		if err != nil {
			return nil, err
		}
		activeRules, err = cachedRulesFromDelta(r.ActiveRulesSorted, delta)
		if err != nil {
			return nil, err
		}
	}

	// Get the version of set of rules read by ReadActive/ReadDelta.
	// Must occur in the same spanner transaction as ReadActive/ReadDelta.
	// If the project has no rules, this returns rules.StartingEpoch.
	rulesVersion, err := rules.ReadVersion(txn, r.Project)
	if err != nil {
		return nil, err
	}

	return NewRuleset(r.Project, activeRules, rulesVersion, lastRefresh), nil
}

// cachedRulesFromFullRead obtains a set of cached rules from the given set of
// active failure association rules.
func cachedRulesFromFullRead(activeRules []*rules.FailureAssociationRule) ([]*CachedRule, error) {
	var result []*CachedRule
	for _, r := range activeRules {
		cr, err := NewCachedRule(r)
		if err != nil {
			return nil, errors.Annotate(err, "rule %s is invalid", r.RuleID).Err()
		}
		result = append(result, cr)
	}
	return result, nil
}

// cachedRulesFromDelta applies deltas to an existing list of rules,
// to obtain an updated set of rules.
func cachedRulesFromDelta(existing []*CachedRule, delta []*rules.FailureAssociationRule) ([]*CachedRule, error) {
	ruleByID := make(map[string]*CachedRule)
	for _, r := range existing {
		ruleByID[r.Rule.RuleID] = r
	}
	for _, d := range delta {
		if d.IsActive {
			cr, err := NewCachedRule(d)
			if err != nil {
				return nil, errors.Annotate(err, "rule %s is invalid", d.RuleID).Err()
			}
			ruleByID[d.RuleID] = cr
		} else {
			// Delete the rule, if it exists.
			delete(ruleByID, d.RuleID)
		}
	}
	var results []*CachedRule
	for _, r := range ruleByID {
		results = append(results, r)
	}
	return results, nil
}

// sortByDescendingPredicateLastUpdated sorts the given rules in descending
// predicate last-updated time order. If two rules have the same
// PredicateLastUpdated time, they are sorted in RuleID order.
func sortByDescendingPredicateLastUpdated(rules []*CachedRule) []*CachedRule {
	result := make([]*CachedRule, len(rules))
	copy(result, rules)
	sort.Slice(result, func(i, j int) bool {
		if result[i].Rule.PredicateLastUpdated.Equal(result[j].Rule.PredicateLastUpdated) {
			return result[i].Rule.RuleID < result[j].Rule.RuleID
		}
		return result[i].Rule.PredicateLastUpdated.After(result[j].Rule.PredicateLastUpdated)
	})
	return result
}

// rulesByID creates a mapping from rule ID to rules for the given list
// of failure association rules.
func rulesByID(rules []*CachedRule) map[string]*CachedRule {
	result := make(map[string]*CachedRule)
	for _, r := range rules {
		result[r.Rule.RuleID] = r
	}
	return result
}
