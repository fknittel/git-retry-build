// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

import {assert} from 'chai';
import sinon from 'sinon';

import * as example from 'shared/test/constants-hotlist.js';

import {MrHotlistSettingsPage} from './mr-hotlist-settings-page.js';

/** @type {MrHotlistSettingsPage} */
let element;

describe('mr-hotlist-settings-page', () => {
  beforeEach(() => {
    // @ts-ignore
    element = document.createElement('mr-hotlist-settings-page');
    document.body.appendChild(element);
  });

  afterEach(() => {
    document.body.removeChild(element);
  });

  it('initializes', async () => {
    assert.instanceOf(element, MrHotlistSettingsPage);
  });

  it('shows loading message with null hotlist', async () => {
    await element.updateComplete;
    assert.include(element.shadowRoot.innerHTML, 'Loading');
  });

  it('renders hotlist', async () => {
    sinon.stub(element, 'stateChanged');
    element._hotlist = example.HOTLIST;
    await element.updateComplete;
  });

  it('renders private hotlist', async () => {
    sinon.stub(element, 'stateChanged');
    element._hotlist = {...example.HOTLIST, hotlistPrivacy: 0};
    await element.updateComplete;
    assert.include(element.shadowRoot.innerHTML, 'Members only');
  });
});
