# Copyright 2021 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.
from . import add_windows_package
from . import add_windows_driver

from PB.recipes.infra.windows_image_builder import windows_image_builder as wib


class Customization(object):
  """ Base customization class. Provides support for pinning and executing
      recipes.
  """

  def __init__(self, image, cust, arch, scripts, configs, step, path,
               powershell, m_file, archive, source):
    """ __init__ copies common module objects to class references. These are
        commonly used for all customizations
        Args:
          image: wib.Image proto object
          cust: wib.Customization proto object
          arch: string representing architecture to build the image for
          scripts: path to the scripts resource dir
          step: module object for recipe_engine/step
          path: module object for recipe_engine/path
          powershell: module object for recipe_modules/powershell
          m_file: module object for recipe_engine/file
          archive: module object for recipe_engine/archive
          source: module object for Source from sources.py
    """
    # generate a copy of image
    self._image = wib.Image()
    self._image.CopyFrom(image)
    # remove all customizations from the image
    for _ in image.customizations:
      image.customizations.pop()
    # generate a copy of customization
    self._customization = wib.Customization()
    self._customization.CopyFrom(cust)
    # Add customization to the image
    self._image.customizations.append(self._customization)
    self._arch = arch
    self._scripts = scripts
    self._step = step
    self._path = path
    self._powershell = powershell
    self._source = source
    self._file = m_file
    self._archive = archive
    self._key = ''
    self._configs = configs
    self._name = ''

  def name(self):
    """ name returns the name of the customization object. This needs to be
        set by the inheriting class"""
    return self._name

  def customization(self):
    """customization returns wib.Customization object representing self"""
    return self._customization

  def image(self):
    """ image returns wib.Image object representing self"""
    return self._image

  def set_key(self, key):
    """ set_key is used to set the identification keys for the customization
        Args:
          key: string representing the unique key for this customization
    """
    self._key = key

  def execute_script(self, name, command, *args, **kwargs):
    """ Executes the windows powershell script
        Args:
          name: string representing step name
          command: string|path representing command to be run
          args: args to be passed on to the command
          kwargs: logs and ret_codes, logs ([]str) are list os paths to watch
                  and record logs from. ret_codes ([]int) is a list of ints,
                  these will be treated as success return codes upon execution
    """
    logs = kwargs['logs']
    ret_codes = kwargs['ret_codes']
    return self._powershell(
        name, command, logs=logs, ret_codes=ret_codes, args=list(args))
