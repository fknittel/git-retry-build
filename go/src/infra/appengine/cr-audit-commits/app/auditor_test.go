// Copyright 2017 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"runtime"
	"testing"
	"time"

	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"go.chromium.org/luci/common/proto"
	"go.chromium.org/luci/common/proto/git"
	gitilespb "go.chromium.org/luci/common/proto/gitiles"
	"go.chromium.org/luci/common/proto/gitiles/mock_gitiles"
	"go.chromium.org/luci/gae/impl/memory"
	ds "go.chromium.org/luci/gae/service/datastore"
	"go.chromium.org/luci/server/router"

	"infra/appengine/cr-audit-commits/app/config"
	"infra/appengine/cr-audit-commits/app/rules"
	"infra/monorail"
)

type errorRule struct{}

// GetName returns the name of the rule.
func (rule errorRule) GetName() string {
	return "Dummy Rule"
}

// Run return errors if the commit hasn't been audited.
func (rule errorRule) Run(c context.Context, ap *rules.AuditParams, rc *rules.RelevantCommit, cs *rules.Clients) (*rules.RuleResult, error) {
	if rc.Status == rules.AuditScheduled {
		return nil, fmt.Errorf("error rule")
	}
	return &rules.RuleResult{
		RuleName:         "Dummy rule",
		RuleResultStatus: rules.RuleFailed,
		Message:          "",
		MetaData:         "",
	}, fmt.Errorf("error rule")
}

type dummyNotifier struct{}

func (d dummyNotifier) Notify(ctx context.Context, cfg *rules.RefConfig, rc *rules.RelevantCommit, cs *rules.Clients, state string) (string, error) {
	return "NotificationSent", nil
}

type strmap map[string]string

func TestTaskAuditor(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows: flaky")
	}

	Convey("CommitScanner handler test", t, func() {
		ctx := memory.Use(context.Background())

		auditorPath := "/_task/auditor"

		withTestingContext := func(c *router.Context, next router.Handler) {
			c.Context = ctx
			ds.GetTestable(ctx).CatchupIndexes()
			next(c)
		}

		r := router.New()
		r.GET(auditorPath, router.NewMiddlewareChain(withTestingContext), taskAuditor)
		srv := httptest.NewServer(r)
		client := &http.Client{}
		auditorTestClients = &rules.Clients{}
		Convey("Unknown Ref", func() {
			resp, err := client.Get(srv.URL + auditorPath + "?refUrl=unknown")
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, 400)

		})
		Convey("Dummy Repo", func() {
			// Mock global configuration.
			expectedRuleMap := map[string]*rules.RefConfig{
				"dummy-repo": {
					BaseRepoURL:    "https://dummy.googlesource.com/dummy.git",
					GerritURL:      "https://dummy-review.googlesource.com",
					BranchName:     "refs/heads/master",
					StartingCommit: "000000",
					Rules: map[string]rules.AccountRules{"rules": {
						Account: "dummy@test.com",
						Rules: []rules.Rule{
							rules.DummyRule{
								Name: "Dummy rule",
								Result: &rules.RuleResult{
									RuleName:         "Dummy rule",
									RuleResultStatus: rules.RulePassed,
									Message:          "",
									MetaData:         "",
								},
							},
						},
						Notification: dummyNotifier{},
					}},
				},
			}
			configGetOld := configGet
			configGet = func(context.Context) map[string]*rules.RefConfig {
				return expectedRuleMap
			}
			defer func() {
				configGet = configGetOld
			}()

			ctl := gomock.NewController(t)
			defer ctl.Finish()

			escapedRepoURL := url.QueryEscape("https://dummy.googlesource.com/dummy.git/+/refs/heads/master")
			gitilesMockClient := mock_gitiles.NewMockGitilesClient(ctl)
			auditorTestClients.GitilesFactory = func(host string, httpClient *http.Client) (gitilespb.GitilesClient, error) {
				return gitilesMockClient, nil
			}
			Convey("Test scanning", func() {
				ds.Put(ctx, &rules.RepoState{
					RepoURL:            "https://dummy.googlesource.com/dummy.git/+/refs/heads/master",
					ConfigName:         "dummy-repo",
					LastKnownCommit:    "123456",
					LastRelevantCommit: "999999",
					LastUpdatedTime:    time.Now().UTC(),
				})

				Convey("No revisions", func() {
					gitilesMockClient.EXPECT().Refs(gomock.Any(), proto.MatcherEqual(&gitilespb.RefsRequest{
						Project:  "dummy",
						RefsPath: "refs/heads/master",
					})).Return(&gitilespb.RefsResponse{
						Revisions: strmap{"refs/heads/master/refs/heads/master": "123456"},
					}, nil)
					gitilesMockClient.EXPECT().Log(gomock.Any(), proto.MatcherEqual(&gitilespb.LogRequest{
						Project:            "dummy",
						Committish:         "123456",
						ExcludeAncestorsOf: "123456",
						PageSize:           int32(config.MaxCommitsPerRefUpdate),
					})).Return(&gitilespb.LogResponse{
						Log: []*git.Commit{},
					}, nil)
					resp, err := client.Get(srv.URL + auditorPath + "?refUrl=" + escapedRepoURL)
					So(err, ShouldBeNil)
					So(resp.StatusCode, ShouldEqual, 200)
					rs := &rules.RepoState{RepoURL: "https://dummy.googlesource.com/dummy.git/+/refs/heads/master"}
					err = ds.Get(ctx, rs)
					So(err, ShouldBeNil)
					So(rs.LastKnownCommit, ShouldEqual, "123456")
					So(rs.LastRelevantCommit, ShouldEqual, "999999")
					So(rs.Paused, ShouldEqual, false)
				})
				Convey("No interesting revisions", func() {
					gitilesMockClient.EXPECT().Refs(gomock.Any(), proto.MatcherEqual(&gitilespb.RefsRequest{
						Project:  "dummy",
						RefsPath: "refs/heads/master",
					})).Return(&gitilespb.RefsResponse{
						Revisions: strmap{"refs/heads/master/refs/heads/master": "abcdef000123123"},
					}, nil)
					gitilesMockClient.EXPECT().Log(gomock.Any(), proto.MatcherEqual(&gitilespb.LogRequest{
						Project:            "dummy",
						Committish:         "abcdef000123123",
						ExcludeAncestorsOf: "123456",
						PageSize:           int32(config.MaxCommitsPerRefUpdate),
					})).Return(&gitilespb.LogResponse{
						Log: []*git.Commit{{Id: "abcdef000123123"}},
					}, nil)
					resp, err := client.Get(srv.URL + auditorPath + "?refUrl=" + escapedRepoURL)
					So(err, ShouldBeNil)
					So(resp.StatusCode, ShouldEqual, 200)
					rs := &rules.RepoState{RepoURL: "https://dummy.googlesource.com/dummy.git/+/refs/heads/master"}
					err = ds.Get(ctx, rs)
					So(err, ShouldBeNil)
					So(rs.LastKnownCommit, ShouldEqual, "abcdef000123123")
					So(rs.LastRelevantCommit, ShouldEqual, "999999")
					So(rs.Paused, ShouldEqual, false)
				})
				Convey("Interesting revisions", func() {
					gitilesMockClient.EXPECT().Refs(gomock.Any(), proto.MatcherEqual(&gitilespb.RefsRequest{
						Project:  "dummy",
						RefsPath: "refs/heads/master",
					})).Return(&gitilespb.RefsResponse{
						Revisions: strmap{"refs/heads/master/refs/heads/master": "deadbeef"},
					}, nil)
					gitilesMockClient.EXPECT().Log(gomock.Any(), proto.MatcherEqual(&gitilespb.LogRequest{
						Project:            "dummy",
						Committish:         "deadbeef",
						ExcludeAncestorsOf: "123456",
						PageSize:           int32(config.MaxCommitsPerRefUpdate),
					})).Return(&gitilespb.LogResponse{
						Log: []*git.Commit{
							{Id: "deadbeef"},
							{
								Id: "c001c0de",
								Author: &git.Commit_User{
									Email: "dummy@test.com",
									Time:  rules.MustGitilesTime("Sun Sep 03 00:56:34 2017"),
								},
								Committer: &git.Commit_User{
									Email: "dummy@test.com",
									Time:  rules.MustGitilesTime("Sun Sep 03 00:56:34 2017"),
								},
							},
						},
					}, nil)
					resp, err := client.Get(srv.URL + auditorPath + "?refUrl=" + escapedRepoURL)
					So(err, ShouldBeNil)
					So(resp.StatusCode, ShouldEqual, 200)
					rs := &rules.RepoState{RepoURL: "https://dummy.googlesource.com/dummy.git/+/refs/heads/master"}
					err = ds.Get(ctx, rs)
					So(err, ShouldBeNil)
					So(rs.LastKnownCommit, ShouldEqual, "deadbeef")
					So(rs.LastRelevantCommit, ShouldEqual, "c001c0de")
					So(rs.Paused, ShouldEqual, false)
					rc := &rules.RelevantCommit{
						RepoStateKey: ds.KeyForObj(ctx, rs),
						CommitHash:   "c001c0de",
					}
					err = ds.Get(ctx, rc)
					So(err, ShouldBeNil)
					So(rc.PreviousRelevantCommit, ShouldEqual, "999999")
				})
				Convey("Force push", func() {
					gitilesMockClient.EXPECT().Refs(gomock.Any(), proto.MatcherEqual(&gitilespb.RefsRequest{
						Project:  "dummy",
						RefsPath: "refs/heads/master",
					})).Return(&gitilespb.RefsResponse{
						Revisions: strmap{"refs/heads/master/refs/heads/master": "abcdef000123123"},
					}, nil)
					gitilesMockClient.EXPECT().Log(gomock.Any(), proto.MatcherEqual(&gitilespb.LogRequest{
						Project:            "dummy",
						Committish:         "abcdef000123123",
						ExcludeAncestorsOf: "123456",
						PageSize:           int32(config.MaxCommitsPerRefUpdate),
					})).Return(nil, grpc.Errorf(codes.NotFound, "not found"))
					gitilesMockClient.EXPECT().Log(gomock.Any(), proto.MatcherEqual(&gitilespb.LogRequest{
						Project:    "dummy",
						Committish: "abcdef000123123",
						PageSize:   1,
					})).Return(&gitilespb.LogResponse{
						Log: []*git.Commit{{Id: "abcdef000123123"}},
					}, nil)
					gitilesMockClient.EXPECT().Log(gomock.Any(), proto.MatcherEqual(&gitilespb.LogRequest{
						Project:    "dummy",
						Committish: "123456",
						PageSize:   1,
					})).Return(nil, grpc.Errorf(codes.NotFound, "not found"))
					auditorTestClients.Monorail = rules.MockMonorailClient{
						Ii: &monorail.InsertIssueResponse{
							Issue: &monorail.Issue{
								Id: 12345,
							},
						},
					}

					resp, err := client.Get(srv.URL + auditorPath + "?refUrl=" + escapedRepoURL)
					So(err, ShouldBeNil)
					So(resp.StatusCode, ShouldEqual, 409)
					rs := &rules.RepoState{RepoURL: "https://dummy.googlesource.com/dummy.git/+/refs/heads/master"}
					err = ds.Get(ctx, rs)
					So(err, ShouldBeNil)
					So(rs.LastKnownCommit, ShouldEqual, "123456")
					So(rs.Paused, ShouldEqual, true)
				})
				Convey("Temporary new head not found error", func() {
					gitilesMockClient.EXPECT().Refs(gomock.Any(), proto.MatcherEqual(&gitilespb.RefsRequest{
						Project:  "dummy",
						RefsPath: "refs/heads/master",
					})).Return(&gitilespb.RefsResponse{
						Revisions: strmap{"refs/heads/master/refs/heads/master": "abcdef000123123"},
					}, nil)
					gitilesMockClient.EXPECT().Log(gomock.Any(), proto.MatcherEqual(&gitilespb.LogRequest{
						Project:            "dummy",
						Committish:         "abcdef000123123",
						ExcludeAncestorsOf: "123456",
						PageSize:           int32(config.MaxCommitsPerRefUpdate),
					})).Return(nil, grpc.Errorf(codes.NotFound, "not found"))
					gitilesMockClient.EXPECT().Log(gomock.Any(), proto.MatcherEqual(&gitilespb.LogRequest{
						Project:    "dummy",
						Committish: "abcdef000123123",
						PageSize:   1,
					})).Return(nil, grpc.Errorf(codes.NotFound, "not found"))
					resp, err := client.Get(srv.URL + auditorPath + "?refUrl=" + escapedRepoURL)
					So(err, ShouldBeNil)
					So(resp.StatusCode, ShouldEqual, 502)
					rs := &rules.RepoState{RepoURL: "https://dummy.googlesource.com/dummy.git/+/refs/heads/master"}
					err = ds.Get(ctx, rs)
					So(err, ShouldBeNil)
					So(rs.LastKnownCommit, ShouldEqual, "123456")
					So(rs.Paused, ShouldEqual, false)
				})
				Convey("Existed relevant commits", func() {
					gitilesMockClient.EXPECT().Refs(gomock.Any(), &gitilespb.RefsRequest{
						Project:  "dummy",
						RefsPath: "refs/heads/master",
					}).Return(&gitilespb.RefsResponse{
						Revisions: strmap{"refs/heads/master/refs/heads/master": "deadbeef"},
					}, nil)
					gitilesMockClient.EXPECT().Log(gomock.Any(), &gitilespb.LogRequest{
						Project:            "dummy",
						Committish:         "deadbeef",
						ExcludeAncestorsOf: "123456",
						PageSize:           int32(config.MaxCommitsPerRefUpdate),
					}).Return(&gitilespb.LogResponse{
						Log: []*git.Commit{
							{Id: "deadbeef"},
							{
								Id: "000002",
								Author: &git.Commit_User{
									Email: "dummy@test.com",
									Time:  rules.MustGitilesTime("Sun Sep 03 00:56:34 2017"),
								},
								Committer: &git.Commit_User{
									Email: "dummy@test.com",
									Time:  rules.MustGitilesTime("Sun Sep 03 00:56:34 2017"),
								},
							},
							{
								Id: "000001",
								Author: &git.Commit_User{
									Email: "dummy@test.com",
									Time:  rules.MustGitilesTime("Sun Sep 03 00:10:34 2017"),
								},
								Committer: &git.Commit_User{
									Email: "dummy@test.com",
									Time:  rules.MustGitilesTime("Sun Sep 03 00:10:34 2017"),
								},
							},
							{
								Id: "000000",
								Author: &git.Commit_User{
									Email: "dummy@test.com",
									Time:  rules.MustGitilesTime("Sun Sep 03 00:00:34 2017"),
								},
								Committer: &git.Commit_User{
									Email: "dummy@test.com",
									Time:  rules.MustGitilesTime("Sun Sep 03 00:00:34 2017"),
								},
							},
						},
					}, nil)
					rs := &rules.RepoState{RepoURL: "https://dummy.googlesource.com/dummy.git/+/refs/heads/master"}
					rsk := ds.KeyForObj(ctx, rs)
					ds.Put(ctx, []*rules.RelevantCommit{
						{
							RepoStateKey:           rsk,
							CommitHash:             "999999",
							Status:                 rules.AuditCompleted,
							PreviousRelevantCommit: "999998",
						},
						{
							RepoStateKey:           rsk,
							CommitHash:             "000000",
							Status:                 rules.AuditCompletedWithActionRequired,
							PreviousRelevantCommit: "999999",
						},
						{
							RepoStateKey:           rsk,
							CommitHash:             "000001",
							Status:                 rules.AuditCompletedWithActionRequired,
							PreviousRelevantCommit: "000000",
						},
					})

					resp, err := client.Get(srv.URL + auditorPath + "?refUrl=" + escapedRepoURL)
					So(err, ShouldBeNil)
					So(resp.StatusCode, ShouldEqual, 200)
					err = ds.Get(ctx, rs)
					So(err, ShouldBeNil)
					So(rs.LastKnownCommit, ShouldEqual, "deadbeef")
					So(rs.LastRelevantCommit, ShouldEqual, "000002")

					ds.GetTestable(ctx).CatchupIndexes()
					var rcs []*rules.RelevantCommit
					err = ds.GetAll(ctx, ds.NewQuery("RelevantCommit").Ancestor(rsk), &rcs)

					So(err, ShouldBeNil)
					So(rcs[0].PreviousRelevantCommit, ShouldEqual, "999999")
					So(rcs[0].Status, ShouldEqual, rules.AuditCompletedWithActionRequired)
					So(rcs[1].PreviousRelevantCommit, ShouldEqual, "000000")
					So(rcs[1].Status, ShouldEqual, rules.AuditCompletedWithActionRequired)
					So(rcs[2].PreviousRelevantCommit, ShouldEqual, "000001")
					So(rcs[2].Status, ShouldEqual, rules.AuditCompleted)
				})
			})
			Convey("Test auditing", func() {
				repoState := &rules.RepoState{
					ConfigName:         "dummy-repo",
					RepoURL:            "https://dummy.googlesource.com/dummy.git/+/refs/heads/master",
					LastKnownCommit:    "222222",
					LastRelevantCommit: "222222",
					LastUpdatedTime:    time.Now().UTC(),
				}
				err := ds.Put(ctx, repoState)
				rsk := ds.KeyForObj(ctx, repoState)

				So(err, ShouldBeNil)
				gitilesMockClient.EXPECT().Refs(gomock.Any(), proto.MatcherEqual(&gitilespb.RefsRequest{
					Project:  "dummy",
					RefsPath: "refs/heads/master",
				})).Return(&gitilespb.RefsResponse{
					Revisions: strmap{"refs/heads/master/refs/heads/master": "222222"},
				}, nil)
				gitilesMockClient.EXPECT().Log(gomock.Any(), proto.MatcherEqual(&gitilespb.LogRequest{
					Project:            "dummy",
					Committish:         "222222",
					ExcludeAncestorsOf: "222222",
					PageSize:           int32(config.MaxCommitsPerRefUpdate),
				})).Return(&gitilespb.LogResponse{
					Log: []*git.Commit{},
				}, nil)

				Convey("No commits", func() {
					resp, err := client.Get(srv.URL + auditorPath + "?refUrl=" + escapedRepoURL)
					So(err, ShouldBeNil)
					So(resp.StatusCode, ShouldEqual, 200)
				})
				Convey("With commits", func() {
					for i := 0; i < 10; i++ {
						rc := &rules.RelevantCommit{
							RepoStateKey:  rsk,
							CommitHash:    fmt.Sprintf("%02d%02d%02d", i, i, i),
							Status:        rules.AuditScheduled,
							AuthorAccount: "dummy@test.com",
						}
						err := ds.Put(ctx, rc)
						So(err, ShouldBeNil)
					}
					Convey("All pass", func() {
						resp, err := client.Get(srv.URL + auditorPath + "?refUrl=" + escapedRepoURL)
						So(err, ShouldBeNil)
						So(resp.StatusCode, ShouldEqual, 200)
						for i := 0; i < 10; i++ {
							rc := &rules.RelevantCommit{
								RepoStateKey: rsk,
								CommitHash:   fmt.Sprintf("%02d%02d%02d", i, i, i),
							}
							err := ds.Get(ctx, rc)
							So(err, ShouldBeNil)
							So(rc.Status, ShouldEqual, rules.AuditCompleted)
						}
					})
					Convey("Some fail", func() {
						dummyRuleTmp := expectedRuleMap["dummy-repo"].Rules["rules"].Rules[0].(rules.DummyRule)
						dummyRuleTmp.Result.RuleResultStatus = rules.RuleFailed
						resp, err := client.Get(srv.URL + auditorPath + "?refUrl=" + escapedRepoURL)
						So(err, ShouldBeNil)
						So(resp.StatusCode, ShouldEqual, 200)
						for i := 0; i < 10; i++ {
							rc := &rules.RelevantCommit{
								RepoStateKey: rsk,
								CommitHash:   fmt.Sprintf("%02d%02d%02d", i, i, i),
							}
							err := ds.Get(ctx, rc)
							So(err, ShouldBeNil)
							So(rc.Status, ShouldEqual, rules.AuditCompletedWithActionRequired)
						}
					})
					Convey("Some error", func() {
						expectedRuleMap["dummy-repo"].Rules["rules"].Rules[0] = errorRule{}
						resp, err := client.Get(srv.URL + auditorPath + "?refUrl=" + escapedRepoURL)
						So(err, ShouldBeNil)
						So(resp.StatusCode, ShouldEqual, 200)
						for i := 0; i < 10; i++ {
							rc := &rules.RelevantCommit{
								RepoStateKey: rsk,
								CommitHash:   fmt.Sprintf("%02d%02d%02d", i, i, i),
							}
							err := ds.Get(ctx, rc)
							So(err, ShouldBeNil)
							So(rc.Status, ShouldEqual, rules.AuditScheduled)
							So(rc.Retries, ShouldEqual, 1)
						}
					})

				})
			})
			Convey("Test unpausing", func() {
				auditorTestClients.Monorail = rules.MockMonorailClient{
					Ii: &monorail.InsertIssueResponse{
						Issue: &monorail.Issue{
							Id: 12345,
						},
					},
				}
				Convey("Unpausing successfully", func() {
					gitilesMockClient.EXPECT().Refs(gomock.Any(), &gitilespb.RefsRequest{
						Project:  "dummy",
						RefsPath: "refs/heads/master",
					}).Return(&gitilespb.RefsResponse{
						Revisions: strmap{"refs/heads/master/refs/heads/master": "abcdef000123123"},
					}, nil)
					gitilesMockClient.EXPECT().Log(gomock.Any(), &gitilespb.LogRequest{
						Project:            "dummy",
						Committish:         "abcdef000123123",
						ExcludeAncestorsOf: "999999",
						PageSize:           int32(config.MaxCommitsPerRefUpdate),
					}).Return(&gitilespb.LogResponse{
						Log: []*git.Commit{},
					}, nil)

					ds.Put(ctx, &rules.RepoState{
						RepoURL:            "https://dummy.googlesource.com/dummy.git/+/refs/heads/master",
						ConfigName:         "dummy-repo",
						LastKnownCommit:    "123456",
						LastRelevantCommit: "999999",
						LastUpdatedTime:    time.Now().Add(-config.StuckScannerDuration).UTC(),
					})

					resp, err := client.Get(srv.URL + auditorPath + "?refUrl=" + escapedRepoURL)
					So(err, ShouldBeNil)
					So(resp.StatusCode, ShouldEqual, 409)
					rs := &rules.RepoState{RepoURL: "https://dummy.googlesource.com/dummy.git/+/refs/heads/master"}
					err = ds.Get(ctx, rs)
					So(err, ShouldBeNil)
					So(rs.LastKnownCommit, ShouldEqual, "123456")
					So(rs.Paused, ShouldEqual, true)

					expectedRuleMap["dummy-repo"].OverwriteLastKnownCommit = "999999"
					resp, err = client.Get(srv.URL + auditorPath + "?refUrl=" + escapedRepoURL)
					So(err, ShouldBeNil)
					So(resp.StatusCode, ShouldEqual, 200)
					rs = &rules.RepoState{RepoURL: "https://dummy.googlesource.com/dummy.git/+/refs/heads/master"}
					err = ds.Get(ctx, rs)
					So(err, ShouldBeNil)
					So(rs.LastKnownCommit, ShouldEqual, "999999")
					So(rs.Paused, ShouldEqual, false)
				})
				Convey("An OverwriteLastKnownCommit will not be used twice", func() {
					ds.Put(ctx, &rules.RepoState{
						RepoURL:                          "https://dummy.googlesource.com/dummy.git/+/refs/heads/master",
						ConfigName:                       "dummy-repo",
						LastKnownCommit:                  "123456",
						LastRelevantCommit:               "999999",
						LastUpdatedTime:                  time.Now().Add(-config.StuckScannerDuration).UTC(),
						AcceptedOverwriteLastKnownCommit: "999999",
					})
					resp, err := client.Get(srv.URL + auditorPath + "?refUrl=" + escapedRepoURL)
					So(err, ShouldBeNil)
					So(resp.StatusCode, ShouldEqual, 409)
					rs := &rules.RepoState{RepoURL: "https://dummy.googlesource.com/dummy.git/+/refs/heads/master"}
					err = ds.Get(ctx, rs)
					So(err, ShouldBeNil)
					So(rs.LastKnownCommit, ShouldEqual, "123456")
					So(rs.Paused, ShouldEqual, true)

					expectedRuleMap["dummy-repo"].OverwriteLastKnownCommit = "999999"
					resp, err = client.Get(srv.URL + auditorPath + "?refUrl=" + escapedRepoURL)
					So(err, ShouldBeNil)
					So(resp.StatusCode, ShouldEqual, 409)
					rs = &rules.RepoState{RepoURL: "https://dummy.googlesource.com/dummy.git/+/refs/heads/master"}
					err = ds.Get(ctx, rs)
					So(err, ShouldBeNil)
					So(rs.LastKnownCommit, ShouldEqual, "123456")
					So(rs.Paused, ShouldEqual, true)
				})
			})
			srv.Close()
		})
	})
}
