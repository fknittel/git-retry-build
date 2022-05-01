// Copyright 2022 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package rpc

import (
	"context"
	"fmt"
	"time"

	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/grpc/appstatus"
	"go.chromium.org/luci/server/auth"
	"go.chromium.org/luci/server/span"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"infra/appengine/weetbix/internal/bugs"
	"infra/appengine/weetbix/internal/clustering"
	"infra/appengine/weetbix/internal/clustering/rules"
	"infra/appengine/weetbix/internal/config/compiledcfg"
	configpb "infra/appengine/weetbix/internal/config/proto"
	pb "infra/appengine/weetbix/proto/v1"
)

// Rules implements pb.RulesServer.
type rulesServer struct {
}

// NewRulesSever returns a new pb.RulesServer.
func NewRulesSever() pb.RulesServer {
	return &pb.DecoratedRules{
		Prelude:  checkAllowedPrelude,
		Service:  &rulesServer{},
		Postlude: gRPCifyAndLogPostlude,
	}
}

// Retrieves a rule.
func (*rulesServer) Get(ctx context.Context, req *pb.GetRuleRequest) (*pb.Rule, error) {
	project, ruleID, err := parseRuleName(req.Name)
	if err != nil {
		return nil, invalidArgumentError(err)
	}

	cfg, err := readProjectConfig(ctx, project)
	if err != nil {
		return nil, err
	}

	r, err := rules.Read(span.Single(ctx), project, ruleID)
	if err != nil {
		if err == rules.NotExistsErr {
			return nil, appstatus.Error(codes.NotFound, "rule does not exist")
		}
		// This will result in an internal error being reported to the caller.
		return nil, errors.Annotate(err, "reading rule %s", ruleID).Err()
	}
	return createRulePB(r, cfg.Config), nil
}

// Lists rules.
func (*rulesServer) List(ctx context.Context, req *pb.ListRulesRequest) (*pb.ListRulesResponse, error) {
	project, err := parseProjectName(req.Parent)
	if err != nil {
		return nil, invalidArgumentError(err)
	}

	cfg, err := readProjectConfig(ctx, project)
	if err != nil {
		return nil, err
	}

	// TODO: Update to read all rules (not just active), and implement pagination.
	rs, err := rules.ReadActive(span.Single(ctx), project)
	if err != nil {
		// GRPCifyAndLog will log this, and report an internal error.
		return nil, errors.Annotate(err, "reading rules").Err()
	}

	rpbs := make([]*pb.Rule, 0, len(rs))
	for _, r := range rs {
		rpbs = append(rpbs, createRulePB(r, cfg.Config))
	}
	response := &pb.ListRulesResponse{
		Rules: rpbs,
	}
	return response, nil
}

// Creates a new rule.
func (*rulesServer) Create(ctx context.Context, req *pb.CreateRuleRequest) (*pb.Rule, error) {
	project, err := parseProjectName(req.Parent)
	if err != nil {
		return nil, invalidArgumentError(err)
	}

	cfg, err := readProjectConfig(ctx, project)
	if err != nil {
		return nil, err
	}

	ruleID, err := rules.GenerateID()
	if err != nil {
		return nil, errors.Annotate(err, "generating Rule ID").Err()
	}
	user := auth.CurrentUser(ctx).Email

	r := &rules.FailureAssociationRule{
		Project:        project,
		RuleID:         ruleID,
		RuleDefinition: req.Rule.GetRuleDefinition(),
		BugID: bugs.BugID{
			System: req.Rule.Bug.GetSystem(),
			ID:     req.Rule.Bug.GetId(),
		},
		IsActive:      req.Rule.GetIsActive(),
		IsManagingBug: req.Rule.GetIsManagingBug(),
		SourceCluster: clustering.ClusterID{
			Algorithm: req.Rule.SourceCluster.GetAlgorithm(),
			ID:        req.Rule.SourceCluster.GetId(),
		},
	}

	if err := validateBugAgainstConfig(cfg, r.BugID); err != nil {
		return nil, invalidArgumentError(err)
	}

	commitTime, err := span.ReadWriteTransaction(ctx, func(ctx context.Context) error {
		// Verify the bug is not used by another rule in this project.
		bugRules, err := rules.ReadByBug(ctx, r.BugID)
		if err != nil {
			return err
		}
		for _, otherRule := range bugRules {
			if otherRule.IsManagingBug {
				// Avoid conflicts by silently making the bug not managed
				// by this rule if there is another rule managing it.
				// Note: this validation implicitly discloses the existence
				// of rules in projects other than those the user may have
				// access to.
				r.IsManagingBug = false
			}
			if otherRule.Project == r.Project {
				return invalidArgumentError(fmt.Errorf("bug already used by a rule in the same project (%s/%s)", otherRule.Project, otherRule.RuleID))
			}
		}

		err = rules.Create(ctx, r, user)
		if err != nil {
			return invalidArgumentError(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	r.CreationTime = commitTime.In(time.UTC)
	r.CreationUser = user
	r.LastUpdated = commitTime.In(time.UTC)
	r.LastUpdatedUser = user
	r.PredicateLastUpdated = commitTime.In(time.UTC)

	return createRulePB(r, cfg.Config), nil
}

// Updates a rule.
func (*rulesServer) Update(ctx context.Context, req *pb.UpdateRuleRequest) (*pb.Rule, error) {
	project, ruleID, err := parseRuleName(req.Rule.GetName())
	if err != nil {
		return nil, invalidArgumentError(err)
	}

	cfg, err := readProjectConfig(ctx, project)
	if err != nil {
		return nil, err
	}

	user := auth.CurrentUser(ctx).Email

	var predicateUpdated bool
	var updatedRule *rules.FailureAssociationRule
	f := func(ctx context.Context) error {
		rule, err := rules.Read(ctx, project, ruleID)
		if err != nil {
			if err == rules.NotExistsErr {
				return appstatus.Error(codes.NotFound, "rule does not exist")
			}
			// This will result in an internal error being reported to the
			// caller.
			return errors.Annotate(err, "read rule").Err()
		}
		if req.Etag != "" && ruleETag(rule) != req.Etag {
			// Attach a codes.Aborted appstatus to a vanilla error to avoid
			// ReadWriteTransaction interpreting this case for a scenario
			// in which it should retry the transaction.
			err := errors.New("etag mismatch")
			return appstatus.Attach(err, status.New(codes.Aborted, "the rule was modified since it was last read; the update was not applied."))
		}
		updatePredicate := false
		updatingBug := false
		updatingManaged := false
		for _, path := range req.UpdateMask.Paths {
			// Only limited fields may be modified by the client.
			switch path {
			case "rule_definition":
				rule.RuleDefinition = req.Rule.RuleDefinition
				updatePredicate = true
			case "bug":
				bugID := bugs.BugID{
					System: req.Rule.Bug.GetSystem(),
					ID:     req.Rule.Bug.GetId(),
				}
				if err := validateBugAgainstConfig(cfg, bugID); err != nil {
					return invalidArgumentError(err)
				}

				updatingBug = true // Triggers validation.
				rule.BugID = bugID
			case "is_active":
				rule.IsActive = req.Rule.IsActive
				updatePredicate = true
			case "is_managing_bug":
				updatingManaged = true // Triggers validation.
				rule.IsManagingBug = req.Rule.IsManagingBug
			default:
				return invalidArgumentError(fmt.Errorf("unsupported field mask: %s", path))
			}
		}

		if updatingBug || updatingManaged {
			// Verify the new bug is not used by another rule in the
			// same project, and that there are not multiple rules
			// managing the same bug.
			bugRules, err := rules.ReadByBug(ctx, rule.BugID)
			if err != nil {
				// This will result in an internal error being reported
				// to the caller.
				return err
			}
			for _, otherRule := range bugRules {
				if otherRule.Project == project && otherRule.RuleID != ruleID {
					return invalidArgumentError(fmt.Errorf("bug already used by a rule in the same project (%s/%s)", otherRule.Project, otherRule.RuleID))
				}
			}
			for _, otherRule := range bugRules {
				if otherRule.Project != project && otherRule.IsManagingBug {
					if updatingManaged && rule.IsManagingBug {
						// The caller explicitly requested an update of
						// IsManagingBug to true, but we cannot do this.
						return invalidArgumentError(fmt.Errorf("bug already managed by a rule in another project (%s/%s)", otherRule.Project, otherRule.RuleID))
					}
					// If only changing the bug, avoid conflicts by silently
					// making the bug not managed by this rule if there is
					// another rule managing it.
					// Note: this validation implicitly discloses the existence
					// of rules in projects other than those the user may have
					// access to.
					rule.IsManagingBug = false
				}
			}
		}

		if err := rules.Update(ctx, rule, updatePredicate, user); err != nil {
			return invalidArgumentError(err)
		}
		updatedRule = rule
		predicateUpdated = updatePredicate
		return nil
	}
	commitTime, err := span.ReadWriteTransaction(ctx, f)
	if err != nil {
		return nil, err
	}
	updatedRule.LastUpdated = commitTime.In(time.UTC)
	updatedRule.LastUpdatedUser = user
	if predicateUpdated {
		updatedRule.PredicateLastUpdated = commitTime.In(time.UTC)
	}

	return createRulePB(updatedRule, cfg.Config), nil
}

// LookupBug looks up the rule associated with the given bug.
func (*rulesServer) LookupBug(ctx context.Context, req *pb.LookupBugRequest) (*pb.LookupBugResponse, error) {
	bug := bugs.BugID{
		System: req.System,
		ID:     req.Id,
	}
	if err := bug.Validate(); err != nil {
		return nil, invalidArgumentError(err)
	}
	rules, err := rules.ReadByBug(span.Single(ctx), bug)
	if err != nil {
		// This will result in an internal error being reported to the caller.
		return nil, errors.Annotate(err, "reading rule by bug %s:%s", bug.System, bug.ID).Err()
	}
	ruleNames := make([]string, 0, len(rules))
	for _, rule := range rules {
		ruleNames = append(ruleNames, ruleName(rule.Project, rule.RuleID))
	}
	return &pb.LookupBugResponse{
		Rules: ruleNames,
	}, nil
}

func createRulePB(r *rules.FailureAssociationRule, cfg *configpb.ProjectConfig) *pb.Rule {
	return &pb.Rule{
		Name:           ruleName(r.Project, r.RuleID),
		Project:        r.Project,
		RuleId:         r.RuleID,
		RuleDefinition: r.RuleDefinition,
		Bug:            createAssociatedBugPB(r.BugID, cfg),
		IsActive:       r.IsActive,
		IsManagingBug:  r.IsManagingBug,
		SourceCluster: &pb.ClusterId{
			Algorithm: r.SourceCluster.Algorithm,
			Id:        r.SourceCluster.ID,
		},
		CreateTime:              timestamppb.New(r.CreationTime),
		CreateUser:              r.CreationUser,
		LastUpdateTime:          timestamppb.New(r.LastUpdated),
		LastUpdateUser:          r.LastUpdatedUser,
		PredicateLastUpdateTime: timestamppb.New(r.PredicateLastUpdated),
		Etag:                    ruleETag(r),
	}
}

func ruleETag(rule *rules.FailureAssociationRule) string {
	return fmt.Sprintf(`W/"%s"`, rule.LastUpdated.UTC().Format(time.RFC3339Nano))
}

// validateBugAgainstConfig validates the specified bug is consistent with
// the project configuration.
func validateBugAgainstConfig(cfg *compiledcfg.ProjectConfig, bug bugs.BugID) error {
	switch bug.System {
	case bugs.MonorailSystem:
		project, _, err := bug.MonorailProjectAndID()
		if err != nil {
			return err
		}
		if project != cfg.Config.Monorail.Project {
			return fmt.Errorf("bug not in expected monorail project (%s)", cfg.Config.Monorail.Project)
		}
	case bugs.BuganizerSystem:
		// Buganizer bugs are permitted for all Weetbix projects.
	default:
		return fmt.Errorf("unsupported bug system: %s", bug.System)
	}
	return nil
}
