// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package gs

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	sv "go.chromium.org/chromiumos/infra/proto/go/lab_platform"
	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/common/logging"

	"infra/cmd/stable_version2/internal/utils"
	svlib "infra/cros/stableversion"
	svdata "infra/cros/stableversion/proto"
)

// ParseOmahaStatus the omaha stable version strings.
func ParseOmahaStatus(ctx context.Context, data []byte) ([]*sv.StableCrosVersion, error) {
	res := make([]*sv.StableCrosVersion, 0)
	var omahaDatas svdata.OmahaDatas
	if err := utils.Unmarshaller.Unmarshal(bytes.NewReader(data), &omahaDatas); err != nil {
		return nil, err
	}
	m := make(map[string]string)
	for _, od := range omahaDatas.GetOmahaData() {
		b, v, err := parseOne(od, m)
		if err != nil {
			logging.Debugf(ctx, "fail to parse %s: %s", od.GetBoard(), err)
			continue
		}
		m[b] = v
	}

	for k, v := range m {
		res = append(res, utils.MakeCrOSSV(k, v))
	}
	return res, nil
}

// ParseMetadataResult is the information present in the Omaha metadata file.
type ParseMetadataResult struct {
	BuildTargets        []string
	ModelInBuildTargets map[string][]string
	FirmwareVersions    []*sv.StableFirmwareVersion
}

// ParseMetadata parses the build metadata.
func ParseMetadata(data []byte) (*ParseMetadataResult, error) {
	var bm svdata.BuildMetadata
	if err := utils.Unmarshaller.Unmarshal(bytes.NewReader(data), &bm); err != nil {
		return nil, err
	}
	var res []*sv.StableFirmwareVersion
	var buildTargets []string
	modelInBuildTargets := make(map[string][]string)
	if bm.GetUnibuild() {
		for buildTarget, bm := range bm.GetBoardMetadata() {
			buildTargets = append(buildTargets, buildTarget)
			for model, mm := range bm.GetModels() {
				modelInBuildTargets[buildTarget] = append(modelInBuildTargets[buildTarget], model)
				key := utils.MakeStableVersionKey(buildTarget, model)
				res = append(res, &sv.StableFirmwareVersion{
					Key:     key,
					Version: getFirmwareVersion(mm),
				})
			}
		}
	} else {
		for buildTarget, bm := range bm.GetBoardMetadata() {
			buildTargets = append(buildTargets, buildTarget)
			modelInBuildTargets[buildTarget] = append(modelInBuildTargets[buildTarget], buildTarget)
			key := utils.MakeStableVersionKey(buildTarget, buildTarget)
			res = append(res, &sv.StableFirmwareVersion{
				Key:     key,
				Version: bm.GetMainFirmwareVersion(),
			})
		}
	}
	return &ParseMetadataResult{
		buildTargets,
		modelInBuildTargets,
		res,
	}, nil
}

func parseOne(od *svdata.OmahaData, m map[string]string) (string, string, error) {
	if validEntry(od) {
		b, v := parseEntry(od)
		ov, ok := m[b]
		if !ok {
			return b, v, nil
		}
		cp, err := svlib.CompareCrOSVersions(v, ov)
		if err != nil {
			return "", "", err
		}
		if cp == 1 {
			return b, v, nil
		}
		return "", "", errors.New(fmt.Sprintf("%s smaller than existing versions %s", v, ov))
	}
	return "", "", errors.New("not in beta channel")
}

func validEntry(od *svdata.OmahaData) bool {
	return od.GetChannel() == "beta"
}

func parseEntry(od *svdata.OmahaData) (string, string) {
	b := correctBuildTarget(od.GetBoard().GetPublicCodename())
	v := fmt.Sprintf("R%d-%s", od.GetMilestone(), od.GetChromeOsVersion())
	return b, v
}

func correctBuildTarget(b string) string {
	return strings.Replace(b, "-", "_", -1)
}

func getFirmwareVersion(mm *svdata.ModelMetadata) string {
	if mm.GetReadwriteFirmwareVersion() != "" {
		return mm.GetReadwriteFirmwareVersion()
	}
	return mm.GetReadonlyFirmwareVersion()
}
