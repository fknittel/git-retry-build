// Copyright 2016 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package som implements HTTP server that handles requests to default module.
package som

import (
	"encoding/json"
	"net/http"

	"github.com/luci/luci-go/common/logging"
	"github.com/luci/luci-go/common/tsmon/metric"
	"github.com/luci/luci-go/server/auth/xsrf"
	"github.com/luci/luci-go/server/router"
)

var (
	jsErrors = metric.NewCounter("frontend/js_errors",
		"Number of uncaught javascript errors.", nil)
)

type eCatcherReq struct {
	Errors    map[string]int64 `json:"errors"`
	XSRFToken string           `json:"xsrf_token"`
}

// PostClientMonHandler handles uncaught javascript errors reported by the client.
func PostClientMonHandler(ctx *router.Context) {
	c, w, r := ctx.Context, ctx.Writer, ctx.Request

	req := &eCatcherReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		errStatus(c, w, http.StatusBadRequest, err.Error())
		return
	}

	if err := xsrf.Check(c, req.XSRFToken); err != nil {
		errStatus(c, w, http.StatusForbidden, err.Error())
		return
	}

	for _, errCount := range req.Errors {
		jsErrors.Add(c, errCount)
	}
	logging.Errorf(c, "clientmon report: %v", req.Errors)
}
