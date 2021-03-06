// Copyright 2017 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package rules

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"go.chromium.org/luci/common/proto"
	"go.chromium.org/luci/common/proto/git"
	gitilespb "go.chromium.org/luci/common/proto/gitiles"
	"go.chromium.org/luci/common/proto/gitiles/mock_gitiles"
)

func TestAutoRollRules(t *testing.T) {
	t.Parallel()
	Convey("AutoRoll rules work", t, func() {
		ctx := context.Background()
		rc := &RelevantCommit{
			CommitHash:       "b07c0de",
			Status:           AuditScheduled,
			CommitTime:       time.Date(2017, time.August, 25, 15, 0, 0, 0, time.UTC),
			CommitterAccount: "autoroller@sample.com",
			AuthorAccount:    "autoroller@sample.com",
			CommitMessage:    "Roll dep ABC..XYZ",
		}
		cfg := &RefConfig{
			BaseRepoURL: "https://a.googlesource.com/a.git",
			GerritURL:   "https://a-review.googlesource.com/",
			BranchName:  "master",
		}
		ap := &AuditParams{
			TriggeringAccount: "releasebot@sample.com",
			RepoCfg:           cfg,
		}
		testClients := &Clients{}

		Convey("Only modifies DEPS", func() {
			// Inject gitiles log response
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			gitilesMockClient := mock_gitiles.NewMockGitilesClient(ctl)
			testClients.GitilesFactory = func(host string, httpClient *http.Client) (gitilespb.GitilesClient, error) {
				return gitilesMockClient, nil
			}
			gitilesMockClient.EXPECT().Log(gomock.Any(), proto.MatcherEqual(&gitilespb.LogRequest{
				Project:    "a",
				Committish: "b07c0de",
				PageSize:   1,
				TreeDiff:   true,
			})).Return(&gitilespb.LogResponse{
				Log: []*git.Commit{
					{
						Id: "b07c0de",
						TreeDiff: []*git.Commit_TreeDiff{
							{
								Type:    git.Commit_TreeDiff_MODIFY,
								OldPath: "DEPS",
								NewPath: "DEPS",
							},
						},
					},
				},
			}, nil)
			// Run rule
			rr, _ := AutoRollRules(rc.CommitterAccount, []string{"DEPS"}, nil).Rules[0].Run(ctx, ap, rc, testClients)
			// Check result code
			So(rr.RuleResultStatus, ShouldEqual, RulePassed)
		})
		Convey("Introduces unexpected changes", func() {
			Convey("Modifies other file", func() {
				// Inject gitiles log response
				ctl := gomock.NewController(t)
				defer ctl.Finish()
				gitilesMockClient := mock_gitiles.NewMockGitilesClient(ctl)
				testClients.GitilesFactory = func(host string, httpClient *http.Client) (gitilespb.GitilesClient, error) {
					return gitilesMockClient, nil
				}
				gitilesMockClient.EXPECT().Log(gomock.Any(), proto.MatcherEqual(&gitilespb.LogRequest{
					Project:    "a",
					Committish: "b07c0de",
					PageSize:   1,
					TreeDiff:   true,
				})).Return(&gitilespb.LogResponse{
					Log: []*git.Commit{
						{
							Id: "b07c0de",
							TreeDiff: []*git.Commit_TreeDiff{
								{
									Type:    git.Commit_TreeDiff_MODIFY,
									OldPath: "DEPS",
									NewPath: "DEPS",
								},
								{
									Type:    git.Commit_TreeDiff_ADD,
									NewPath: "other/path",
								},
							},
						},
					},
				}, nil)
				// Run rule
				rr, _ := AutoRollRules(rc.CommitterAccount, []string{"DEPS"}, nil).Rules[0].Run(ctx, ap, rc, testClients)
				// Check result code
				So(rr.RuleResultStatus, ShouldEqual, RuleFailed)
			})
			Convey("Renames DEPS", func() {
				ctl := gomock.NewController(t)
				defer ctl.Finish()
				gitilesMockClient := mock_gitiles.NewMockGitilesClient(ctl)
				testClients.GitilesFactory = func(host string, httpClient *http.Client) (gitilespb.GitilesClient, error) {
					return gitilesMockClient, nil
				}
				gitilesMockClient.EXPECT().Log(gomock.Any(), proto.MatcherEqual(&gitilespb.LogRequest{
					Project:    "a",
					Committish: "b07c0de",
					PageSize:   1,
					TreeDiff:   true,
				})).Return(&gitilespb.LogResponse{
					Log: []*git.Commit{
						{
							Id: "b07c0de",
							TreeDiff: []*git.Commit_TreeDiff{
								{
									Type:    git.Commit_TreeDiff_RENAME,
									OldPath: "DEPS",
									NewPath: "DEPS.bak",
								},
							},
						},
					},
				}, nil)
				// Run rule
				rr, _ := AutoRollRules(rc.CommitterAccount, []string{"DEPS"}, nil).Rules[0].Run(ctx, ap, rc, testClients)
				// Check result code
				So(rr.RuleResultStatus, ShouldEqual, RuleFailed)
			})
		})
	})
}
