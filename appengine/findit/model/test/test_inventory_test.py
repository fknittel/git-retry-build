# Copyright 2019 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

from datetime import datetime

from parameterized import parameterized

from google.appengine.ext import ndb

from gae_libs.testcase import TestCase
from model.test_inventory import LuciTest


class TestInventoryTest(TestCase):

  @parameterized.expand([(99,), ([tuple()],), (set([99]),),
                         ([('congif1'), (99)],)])
  def testDisabledTestVariantsPropertyValidate(self, disabled_test_variants):
    with self.assertRaises(TypeError):
      LuciTest(
          key=LuciTest.CreateKey('luci_project', 'normalized_step',
                                 'normalized_test'),
          disabled_test_variants=disabled_test_variants,
          last_updated_time=datetime(2019, 6, 28, 0, 0, 0))

  @parameterized.expand([([], set()),
                         (
                             [('config1', 'config2'), ('config1',)],
                             set([('config1', 'config2'), ('config1',)]),
                         ), (
                             set(),
                             set(),
                         ),
                         (
                             set([('config1', 'config2'), ('config1',)]),
                             set([('config1', 'config2'), ('config1',)]),
                         )])
  def testDisabledTestVariantsPropertySerialization(self, input_test_variants,
                                                    output_test_variants):
    test = LuciTest(
        key=LuciTest.CreateKey('luci_project', 'normalized_step',
                               'normalized_test'),
        disabled_test_variants=input_test_variants,
        last_updated_time=datetime(2019, 6, 28, 0, 0, 0))
    test.put()
    test = test.key.get()
    self.assertEqual(output_test_variants, test.disabled_test_variants)

  def testLuciTestCreateKey(self):
    actual_key = LuciTest.CreateKey('luci_project', 'normalized_step',
                                    'normalized_test')
    expected_key = ndb.Key('LuciTest',
                           'luci_project@normalized_step@normalized_test')
    self.assertEqual(expected_key, actual_key)

  def testLuciTestDisabled(self):
    test = LuciTest(
        key=LuciTest.CreateKey('luci_project', 'normalized_step',
                               'normalized_test'),
        disabled_test_variants=set(),
        last_updated_time=datetime(2019, 6, 28, 0, 0, 0))
    self.assertFalse(test.disabled)

    test.disabled_test_variants.add(('disabled_test_variant',))
    self.assertTrue(test.disabled)
