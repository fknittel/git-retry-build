# Copyright 2017 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

from common import constants
from common import exceptions
from gae_libs import appengine_util
from gae_libs.http import auth_util


def IsPrivilegedUser(user_email, is_admin):
  """Returns True if the given email account is authorized for access."""
  return is_admin or (user_email and user_email.endswith('@google.com'))


def IsAllowedClientId(client_id):
  """Returns True if the given client id is allowed."""
  return client_id in constants.ALLOWED_CLIENT_IDS


def CanTriggerNewAnalysis(user_email, is_admin):
  """Returns True if the given email account could trigger a new analysis."""
  if not appengine_util.IsStaging():
    allowed_app_accounts = constants.ALLOWED_APP_ACCOUNTS
  else:
    allowed_app_accounts = constants.ALLOWED_STAGING_APP_ACCOUNTS
  return IsPrivilegedUser(user_email,
                          is_admin) or (user_email in allowed_app_accounts)


def ValidateOauthUserForNewAnalysis():
  """Validates whether the oauth user is authorized to trigger a new analysis.

  Returns:
    A tuple (user_email, is_admin).
    user_email (str): The email address of the oauth user.
    is_admin (bool): True if the oauth user is an Admin.

  Raises:
    common.exceptions.UnauthorizedException if the user has no permission.
  """
  user_email = auth_util.GetOauthUserEmail()

  if not user_email:
    raise exceptions.UnauthorizedException('Unknown user.')

  # For Google service accounts, no need to allow client ids for them,
  # since email address uniquely identifies credentials used.
  # At some point someone might want to use Findit API from a GCE project
  # (*@developer.gserviceaccount.com accounts) or from some script that use
  # service account keys (*@*.iam.gserviceaccount.com accounts).
  if not user_email.endswith('@appspot.gserviceaccount.com'):
    client_id = auth_util.GetOauthClientId()
    if not IsAllowedClientId(client_id):
      raise exceptions.UnauthorizedException(
          'Unknown client id %s.' % client_id)

  is_admin = auth_util.IsCurrentOauthUserAdmin()
  if not CanTriggerNewAnalysis(user_email, is_admin):
    raise exceptions.UnauthorizedException('Unknown email %s' % user_email)

  return user_email, is_admin
