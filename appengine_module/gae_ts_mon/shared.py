# Copyright 2016 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import contextlib
import os

from google.appengine.api import modules
from google.appengine.api import namespace_manager
from google.appengine.ext import ndb

from infra_libs.ts_mon.common import metrics

REGION = 'appengine'
PRODXMON_ENDPOINT = 'https://prodxmon-wbl-pa.googleapis.com/v1:insert'
PRODXMON_SERVICE_ACCOUNT_EMAIL = (
    'app-engine-metric-publishers@'
    'prodx-mon-chrome-infra.google.com.iam.gserviceaccount.com')
INSTANCE_NAMESPACE = 'ts_mon_instance_namespace'
# Duration of inactivity to consider an instance dead.
INSTANCE_EXPIRE_SEC = 30 * 60
INSTANCE_EXPECTED_TO_HAVE_TASK_NUM_SEC = 5 * 60
INTERNAL_CALLBACK_NAME = '__gae_ts_mon_callback'
# The cron request path for the tasknum assignment handler.
CRON_REQUEST_PATH_TASKNUM_ASSIGNER = '/internal/cron/ts_mon/send'


appengine_default_version = metrics.StringMetric(
    'appengine/default_version',
    'Name of the version currently marked as default.',
    None)
started_counter = metrics.CounterMetric(
    'appengine/instances/started',
    'Count the number of GAE instance initializations.',
    None)
shutdown_counter = metrics.CounterMetric(
    'appengine/instances/shutdown',
    'Count the number of GAE instance shutdowns.',
    None)
expired_counter = metrics.CounterMetric(
    'appengine/instances/expired',
    'Count the number of GAE instance expirations due to inactivity.',
    None)


class Instance(ndb.Model):
  """Used to map instances to small integers.

  Each instance "owns" an entity with the key <instance-id>.<version>.<module>.
  `task_num` is a mapping assigned by a cron job to the instance; -1=undefined.
  """

  @classmethod
  def _use_memcache(cls, _):
    return False

  @classmethod
  def _get_kind(cls):
    return 'TSMonInstance'

  task_num = ndb.IntegerProperty(default=-1)
  last_updated = ndb.DateTimeProperty(auto_now_add=True)


def instance_key_id():
  if is_python3_env():
    instance_id = os.getenv('GAE_INSTANCE', '')
    version_name = os.getenv('GAE_VERSION', '')
    service_name = os.getenv('GAE_SERVICE', '')
  else:
    instance_id = modules.get_current_instance_id()
    version_name = modules.get_current_version_name()
    service_name = modules.get_current_module_name()
  return '{}.{}.{}'.format(instance_id, version_name, service_name)


@contextlib.contextmanager
def instance_namespace_context():
  previous_namespace = namespace_manager.get_namespace()
  try:
    namespace_manager.set_namespace(INSTANCE_NAMESPACE)
    yield
  finally:
    namespace_manager.set_namespace(previous_namespace)


def get_instance_entity():
  with instance_namespace_context():
    return Instance.get_or_insert(instance_key_id())


def is_python3_env():
  # GAEv1's python2 env does not expose this variable.
  return os.getenv('GAE_RUNTIME', '').startswith('python3')
