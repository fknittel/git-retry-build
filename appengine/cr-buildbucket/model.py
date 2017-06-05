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
BUILD_TIMEOUT = datetime.timedelta(days=2)


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


class CanaryPreference(messages.Enum):
  # The build system will decide whether to use canary or not
  AUTO = 1
  # Use the production build infrastructure
  PROD = 2
  # Use the canary build infrastructure
  CANARY = 3


class PubSubCallback(ndb.Model):
  """Parameters for a callack push task."""
  topic = ndb.StringProperty(required=True, indexed=False)
  auth_token = ndb.StringProperty(indexed=False)
  user_data = ndb.StringProperty(indexed=False)


class BucketState(ndb.Model):
  """Persistent state of a single bucket."""
  # If True, no new bulids may be leased for this bucket.
  is_paused = ndb.BooleanProperty()


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
  """

  status = msgprop.EnumProperty(BuildStatus, default=BuildStatus.SCHEDULED)
  status_changed_time = ndb.DateTimeProperty(auto_now_add=True)
  update_time = ndb.DateTimeProperty(auto_now=True)

  # Creation time attributes.

  create_time = ndb.DateTimeProperty(auto_now_add=True)
  created_by = auth.IdentityProperty()
  # a generic way to distinguish builds.
  # Different buckets have different permissions.
  bucket = ndb.StringProperty(required=True)
  # a list of tags, where each tag is a string
  # with ":" symbol. The first occurrence of ":" splits tag name and tag
  # value. Contains only tags specified by the build request. Old Build
  # entities do not have this field.
  initial_tags = ndb.StringProperty(repeated=True, indexed=False)
  # superset of initial_tags. May contain auto-added tags.
  tags = ndb.StringProperty(repeated=True)
  # immutable arbitrary build parameters.
  parameters = ndb.JsonProperty()
  # PubSub message parameters for build status change notifications.
  pubsub_callback = ndb.StructuredProperty(PubSubCallback, indexed=False)
  # id of the original build that this build was derived from.
  retry_of = ndb.IntegerProperty()
  # Specifies whether canary of build infrastructure should be used for this
  # build.
  canary_preference = msgprop.EnumProperty(CanaryPreference, indexed=False)

  # Lease-time attributes.

  # current lease expiration date.
  # The moment the build is leased, |lease_expiration_date| is set to
  # (utcnow + lease_duration).
  lease_expiration_date = ndb.DateTimeProperty()
  # None if build is not leased, otherwise a random value.
  # Changes every time a build is leased. Can be used to verify that a client
  # is the leaseholder.
  lease_key = ndb.IntegerProperty(indexed=False)
  # True if the build is currently leased. Otherwise False
  is_leased = ndb.ComputedProperty(lambda self: self.lease_key is not None)
  leasee = auth.IdentityProperty()
  never_leased = ndb.BooleanProperty()

  # Start time attributes.

  # a URL to a build-system-specific build, viewable by a human.
  url = ndb.StringProperty(indexed=False)
  # when the build started. Unknown for old builds.
  start_time = ndb.DateTimeProperty()
  # True if canary build infrastructure is used to run this build.
  # It may be None only in SCHEDULED state. Otherwise it must be True or False.
  # If canary_preference is CANARY, this field value does not have to be True,
  # e.g. if the build infrastructure does not have a canary.
  canary = ndb.BooleanProperty()

  # Completion time attributes.

  complete_time = ndb.DateTimeProperty()
  result = msgprop.EnumProperty(BuildResult)
  result_details = ndb.JsonProperty()
  cancelation_reason = msgprop.EnumProperty(CancelationReason)
  failure_reason = msgprop.EnumProperty(FailureReason)

  # Swarming integration

  swarming_hostname = ndb.StringProperty()
  swarming_task_id = ndb.StringProperty()

  def _pre_put_hook(self):
    """Checks Build invariants before putting."""
    super(Build, self)._pre_put_hook()
    is_started = self.status == BuildStatus.STARTED
    is_completed = self.status == BuildStatus.COMPLETED
    is_canceled = self.result == BuildResult.CANCELED
    is_failure = self.result == BuildResult.FAILURE
    is_leased = self.lease_key is not None
    assert (self.result is not None) == is_completed
    assert (self.cancelation_reason is not None) == is_canceled
    assert (self.failure_reason is not None) == is_failure
    assert not (is_completed and is_leased)
    assert (self.lease_expiration_date is not None) == is_leased
    assert (self.leasee is not None) == is_leased
    # no cover due to a bug in coverage (https://stackoverflow.com/a/35325514)
    assert not self.tags or all(':' in t for t in self.tags)  # pragma: no cover
    assert self.create_time
    assert (self.complete_time is not None) == is_completed
    assert not is_started or self.start_time
    assert not self.start_time or self.start_time >= self.create_time
    assert not self.complete_time or self.complete_time >= self.create_time
    assert (not self.complete_time or not self.start_time or
            self.complete_time >= self.start_time)

  def regenerate_lease_key(self):
    """Changes lease key to a different random int."""
    while True:
      new_key = random.randint(0, 1 << 31)
      if new_key != self.lease_key:  # pragma: no branch
        self.lease_key = new_key
        break

  def clear_lease(self):  # pragma: no cover
    """Clears build's lease attributes."""
    self.lease_key = None
    self.lease_expiration_date = None
    self.leasee = None


class TagIndexEntry(ndb.Model):
  """A single entry in a TagIndex, references a build."""
  created_time = ndb.DateTimeProperty(auto_now_add=True)
  # ID of the build.
  build_id = ndb.IntegerProperty(indexed=False)
  # Bucket of the build.
  bucket = ndb.StringProperty(indexed=False)


class TagIndex(ndb.Model):
  """A custom index of builds by a tag.

  Entity key:
    Entity id is a build tag in the same "<key>:<value>" format that builds use.
    TagIndex has no parent.
  """

  MAX_ENTRY_COUNT = 1000

  # if incomplete, this TagIndex should not be used in search.
  # It is set to True if there are more than MAX_ENTRY_COUNT builds
  # for this tag.
  permanently_incomplete = ndb.BooleanProperty()

  # entries is a superset of all builds that have the tag equal to the id of
  # this entity. It may contain references to non-existent builds or builds that
  # do not actually have this tag; such builds must be ignored.
  #
  # It is sorted by build id in descending order.
  entries = ndb.LocalStructuredProperty(
      TagIndexEntry, repeated=True, indexed=False)


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
