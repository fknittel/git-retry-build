#!/bin/bash
# Copyright 2018 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

set -e
set -x
set -o pipefail

PREFIX="$1"
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

# Make automake relocatable
git apply $SCRIPT_DIR/0001-relocatable.patch

./configure --enable-static --disable-shared \
  --prefix "$PREFIX" \
  --host "$CROSS_TRIPLE"
make install -j $(nproc)
