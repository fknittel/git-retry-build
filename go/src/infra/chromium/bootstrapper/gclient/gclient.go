// Copyright 2022 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package gclient is a package that enables performing gclient operations required by the chromium
// bootstrapper.
package gclient

import (
	"context"
	stderrors "errors"
	"infra/chromium/bootstrapper/cipd"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/common/logging"
)

type Client struct {
	gclientPath string
}

const (
	depotToolsPackage        = "infra/recipe_bundles/chromium.googlesource.com/chromium/tools/depot_tools"
	depotToolsPackageVersion = "refs/heads/main"
)

// NewClient returns a new gclient client that uses the gclient binary at gclientPath.
func NewClient(ctx context.Context, cipdClient *cipd.Client) (*Client, error) {
	logging.Infof(ctx, "resolving CIPD package %s@%s", depotToolsPackage, depotToolsPackageVersion)
	depotToolsInstance, err := cipdClient.ResolveVersion(ctx, depotToolsPackage, depotToolsPackageVersion)
	if err != nil {
		return nil, errors.Annotate(err, "failed to resolve depot tools package version").Err()
	}
	logging.Infof(ctx, "downloading CIPD package %s@%s", depotToolsInstance.PackageName, depotToolsInstance.InstanceID)
	packagePath, err := cipdClient.DownloadPackage(ctx, depotToolsInstance.PackageName, depotToolsInstance.InstanceID, "depot-tools")
	if err != nil {
		return nil, errors.Annotate(err, "failed to download/install depot tools").Err()
	}
	gclientPath := filepath.Join(packagePath, "depot_tools", "gclient")
	return &Client{gclientPath}, nil
}

// NewClientForTesting returns a new gclient client that uses the gclient binary found on the
// machine's path.
func NewClientForTesting() (*Client, error) {
	gclientPath, err := exec.LookPath("gclient")
	if err != nil {
		return nil, errors.Annotate(err, "gclient not on $PATH, please install depot_tools").Err()
	}
	return &Client{gclientPath}, nil
}

func (c *Client) GetDep(ctx context.Context, depsContents, depPath string) (string, error) {
	d, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}

	f := path.Join(d, "DEPS")
	if err := ioutil.WriteFile(f, []byte(depsContents), 0644); err != nil {
		return "", err
	}

	// --deps-file: The DEPS file to get dependency from
	// -r: get revision information about the dep at the given path
	cmd := exec.CommandContext(ctx, c.gclientPath, "getdep", "--deps-file", f, "-r", depPath)
	// Set DEPOT_TOOLS_UPDATE environment variable to 0 to prevent gclient from attempting to
	// update depot tools; just use the recipe bundle as-is (the recipe bundle also doesn't
	// contain the necessary update_depot_tools script)
	cmd.Env = append(os.Environ(), "DEPOT_TOOLS_UPDATE=0")
	output, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if stderrors.As(err, &exitErr) {
			return "", errors.Annotate(err, "gclient failed with output:\n%s", exitErr.Stderr).Err()
		}
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
