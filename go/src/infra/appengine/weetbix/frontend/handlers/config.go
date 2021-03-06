// Copyright 2022 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package handlers

import (
	"go.chromium.org/luci/server/router"
)

type projectConfig struct {
	Project  string    `json:"project"`
	Monorail *monorail `json:"monorail"`
}

type monorail struct {
	Project       string `json:"project"`
	DisplayPrefix string `json:"displayPrefix"`
}

// GetConfig serves a GET request for
// /api/projects/:project/config.
func (h *Handlers) GetConfig(ctx *router.Context) {
	projectID, cfg, ok := obtainProjectConfigOrError(ctx)
	if !ok {
		return
	}
	cfgpb := cfg.Config
	result := &projectConfig{
		Project: projectID,
		Monorail: &monorail{
			Project:       cfgpb.Monorail.Project,
			DisplayPrefix: cfgpb.Monorail.DisplayPrefix,
		},
	}
	respondWithJSON(ctx, result)
}
