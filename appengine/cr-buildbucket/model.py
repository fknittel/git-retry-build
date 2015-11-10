# Copyright 2014 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import datetime
import random

from components import auth
from components import utils
from google.appengine.ext import ndb
from google.appengine.ext.ndb import msgprop
from protorpc import messages


BEGINING_OF_THE_WORLD = datetime.datetime(2010, 1, 1, 0, 0, 0, 0)


class BuildStatus(messages.Enum):
  # A build is created, can be leased by someone and started.
  SCHEDULED = 1
  # Someone has leased the build and marked it as started.
  STARTED = 2
  # A build is completed. See BuildResult for more details.
  COMPLETED = 3


class BuildResult(messages.Enum):
  # A build has completed successfully.
  SUCCESS = 1
  # A build has completed unsuccessfully.
  FAILURE = 2
  # A build was canceled.
  CANCELED = 3


class FailureReason(messages.Enum):
  # Build failed
  BUILD_FAILURE = 1
  # Something happened within buildbucket.
  BUILDBUCKET_FAILURE = 2
  # Something happened with build infrastructure, but not buildbucket.
  INFRA_FAILURE = 3
  # A build-system rejected a build because its definition is invalid.
  INVALID_BUILD_DEFINITION = 4


class CancelationReason(messages.Enum):
  # A build was canceled explicitly, probably by an API call.
  CANCELED_EXPLICITLY = 1
  # A build was canceled by buildbucket due to timeout.
  TIMEOUT = 2


class Callback(ndb.Model):
  """Parameters for a callack push task."""
  url = ndb.StringProperty(required=True, indexed=False)
  headers = ndb.JsonProperty()
  method = ndb.StringProperty(indexed=False)
  queue_name = ndb.StringProperty(indexed=False)


class Build(ndb.Model):
  """Describes a build.

  Build key:
    Build keys are autogenerated, monotonically decreasing integers.
    That is, when sorted by key, new builds are first.
    Build has no parent.

    Build id is a 64 bits integer represented as a string to the user.
    - 1 highest order bit is set to 0 to keep value positive.
    - 43 bits are 43 lower bits of bitwise-inverted time since
      BEGINING_OF_THE_WORLD at 1ms resolution.
      It is good for 2**43 / 365.3 / 24 / 60 / 60 / 1000 = 278 years
      or 2010 + 278 = year 2288.
    - 16 bits are set to a random value. Assuming an instance is internally
      consistent with itself, it can ensure to not reuse the same 16 bits in two
      consecutive requests and/or throttle itself to one request per
      millisecond. Using random value reduces to 2**-15 the probability of
      collision on exact same timestamp at 1ms resolution, so a maximum
      theoretical rate of 65536000 requests/sec but an effective rate in the
      range of ~64k qps without much transaction conflicts. We should be fine.
    - 4 bits are 0. This is to represent the 'version' of the entity
      schema.

    The idea is taken from Swarming TaskRequest entity:
    https://code.google.com/p/swarming/source/browse/appengine/swarming/server/task_request.py#329

  Attributes:
    status (BuildStatus): status of the build.
    bucket (string): a generic way to distinguish builds. Different buckets have
      different permissions.
    tags (list of string): a list of tags, where each tag is a string with ":"
      symbol. The first occurance of ":" splits tag name and tag value.
    parameters (dict): immutable arbitrary build parameters.
    callback (Callback): push task parameters for build status changes.
    lease_expiration_date (datetime): current lease expiration date.
      The moment the build is leased, |lease_expiration_date| is set to
      (utcnow + lease_duration).
    lease_key (int): None if build is not leased, otherwise a random value.
      Changes every time a build is leased. Can be used to verify that a client
      is the leaseholder.
    is_leased (bool): True if the build is currently leased. Otherwise False
    never_leased (bool): True if the build was never leased before.
    url (str): a URL to a build-system-specific build, viewable by a human.
    result (BuildResult): build result.
    cancelation_reason (CancelationReason): why the build was canceled.
  """

  status = msgprop.EnumProperty(BuildStatus, default=BuildStatus.SCHEDULED)
  status_changed_time = ndb.DateTimeProperty(auto_now_add=True)
  update_time = ndb.DateTimeProperty(auto_now=True)

  # Creation time attributes.
  create_time = ndb.DateTimeProperty(auto_now_add=True)
  created_by = auth.IdentityProperty()
  bucket = ndb.StringProperty(required=True)
  tags = ndb.StringProperty(repeated=True)
  parameters = ndb.JsonProperty()
  callback = ndb.StructuredProperty(Callback, indexed=False)

  # Lease-time attributes.
  lease_expiration_date = ndb.DateTimeProperty()
  lease_key = ndb.IntegerProperty(indexed=False)
  is_leased = ndb.ComputedProperty(lambda self: self.lease_key is not None)
  leasee = auth.IdentityProperty()
  never_leased = ndb.BooleanProperty()

  # Start time attributes.
  url = ndb.StringProperty(indexed=False)

  # Completion time attributes.
  complete_time = ndb.DateTimeProperty()
  result = msgprop.EnumProperty(BuildResult)
  result_details = ndb.JsonProperty()
  cancelation_reason = msgprop.EnumProperty(CancelationReason)
  failure_reason = msgprop.EnumProperty(FailureReason)

  def _pre_put_hook(self):
    """Checks Build invariants before putting."""
    super(Build, self)._pre_put_hook()
    is_completed = self.status == BuildStatus.COMPLETED
    assert (self.result is not None) == is_completed
    is_canceled = self.result == BuildResult.CANCELED
    is_failure = self.result == BuildResult.FAILURE
    assert (self.cancelation_reason is not None) == is_canceled
    assert (self.failure_reason is not None) == is_failure
    is_leased = self.lease_key is not None
    assert not (is_completed and is_leased)
    assert (self.lease_expiration_date is not None) == is_leased
    assert (self.leasee is not None) == is_leased
    assert not self.tags or all(':' in t for t in self.tags)

  def regenerate_lease_key(self):
    """Changes lease key to a different random int."""
    while True:
      new_key = random.randint(0, 1 << 31)
      if new_key != self.lease_key:  # pragma: no branch
        self.lease_key = new_key
        break


def new_build_id():
  """Returns a valid ndb.Key for a new Build.

  See model.Build's docstring, "Build key" section.
  """
  # Build ID bits: "0N{43}R{16}V{4}"
  # where N is now bits, R is random bits and V is version bits.
  utcnow = utils.utcnow()
  assert utcnow >= BEGINING_OF_THE_WORLD
  delta = utcnow - BEGINING_OF_THE_WORLD
  now = int(round(delta.total_seconds() * 1000.))
  inverted_now = ~now & ((1 << 43) - 1)
  suffix = random.getrandbits(16)
  return int((inverted_now << 20) | (suffix << 4))
