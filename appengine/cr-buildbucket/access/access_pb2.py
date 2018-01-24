# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: access.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
from google.protobuf import descriptor_pb2
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import duration_pb2 as google_dot_protobuf_dot_duration__pb2
from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='access.proto',
  package='access',
  syntax='proto3',
  serialized_pb=_b('\n\x0c\x61\x63\x63\x65ss.proto\x12\x06\x61\x63\x63\x65ss\x1a\x1egoogle/protobuf/duration.proto\x1a\x1bgoogle/protobuf/empty.proto\"\xc3\x04\n\x13\x44\x65scriptionResponse\x12\x42\n\tresources\x18\x01 \x03(\x0b\x32/.access.DescriptionResponse.ResourceDescription\x1a\xe7\x03\n\x13ResourceDescription\x12\x0c\n\x04kind\x18\x01 \x01(\t\x12\x0f\n\x07\x63omment\x18\x02 \x01(\t\x12M\n\x07\x61\x63tions\x18\x03 \x03(\x0b\x32<.access.DescriptionResponse.ResourceDescription.ActionsEntry\x12I\n\x05roles\x18\x04 \x03(\x0b\x32:.access.DescriptionResponse.ResourceDescription.RolesEntry\x1a\x19\n\x06\x41\x63tion\x12\x0f\n\x07\x63omment\x18\x01 \x01(\t\x1a\x30\n\x04Role\x12\x17\n\x0f\x61llowed_actions\x18\x01 \x03(\t\x12\x0f\n\x07\x63omment\x18\x02 \x01(\t\x1a\x66\n\x0c\x41\x63tionsEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\x45\n\x05value\x18\x02 \x01(\x0b\x32\x36.access.DescriptionResponse.ResourceDescription.Action:\x02\x38\x01\x1a\x62\n\nRolesEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\x43\n\x05value\x18\x02 \x01(\x0b\x32\x34.access.DescriptionResponse.ResourceDescription.Role:\x02\x38\x01\"F\n\x17PermittedActionsRequest\x12\x15\n\rresource_kind\x18\x01 \x01(\t\x12\x14\n\x0cresource_ids\x18\x02 \x03(\t\"\xa4\x02\n\x18PermittedActionsResponse\x12\x42\n\tpermitted\x18\x01 \x03(\x0b\x32/.access.PermittedActionsResponse.PermittedEntry\x12\x34\n\x11validity_duration\x18\x02 \x01(\x0b\x32\x19.google.protobuf.Duration\x1a&\n\x13ResourcePermissions\x12\x0f\n\x07\x61\x63tions\x18\x01 \x03(\t\x1a\x66\n\x0ePermittedEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\x43\n\x05value\x18\x02 \x01(\x0b\x32\x34.access.PermittedActionsResponse.ResourcePermissions:\x02\x38\x01\x32\xa7\x01\n\x06\x41\x63\x63\x65ss\x12W\n\x10PermittedActions\x12\x1f.access.PermittedActionsRequest\x1a .access.PermittedActionsResponse\"\x00\x12\x44\n\x0b\x44\x65scription\x12\x16.google.protobuf.Empty\x1a\x1b.access.DescriptionResponse\"\x00\x62\x06proto3')
  ,
  dependencies=[google_dot_protobuf_dot_duration__pb2.DESCRIPTOR,google_dot_protobuf_dot_empty__pb2.DESCRIPTOR,])
_sym_db.RegisterFileDescriptor(DESCRIPTOR)




_DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ACTION = _descriptor.Descriptor(
  name='Action',
  full_name='access.DescriptionResponse.ResourceDescription.Action',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='comment', full_name='access.DescriptionResponse.ResourceDescription.Action.comment', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=386,
  serialized_end=411,
)

_DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ROLE = _descriptor.Descriptor(
  name='Role',
  full_name='access.DescriptionResponse.ResourceDescription.Role',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='allowed_actions', full_name='access.DescriptionResponse.ResourceDescription.Role.allowed_actions', index=0,
      number=1, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='comment', full_name='access.DescriptionResponse.ResourceDescription.Role.comment', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=413,
  serialized_end=461,
)

_DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ACTIONSENTRY = _descriptor.Descriptor(
  name='ActionsEntry',
  full_name='access.DescriptionResponse.ResourceDescription.ActionsEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='access.DescriptionResponse.ResourceDescription.ActionsEntry.key', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='value', full_name='access.DescriptionResponse.ResourceDescription.ActionsEntry.value', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=_descriptor._ParseOptions(descriptor_pb2.MessageOptions(), _b('8\001')),
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=463,
  serialized_end=565,
)

_DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ROLESENTRY = _descriptor.Descriptor(
  name='RolesEntry',
  full_name='access.DescriptionResponse.ResourceDescription.RolesEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='access.DescriptionResponse.ResourceDescription.RolesEntry.key', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='value', full_name='access.DescriptionResponse.ResourceDescription.RolesEntry.value', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=_descriptor._ParseOptions(descriptor_pb2.MessageOptions(), _b('8\001')),
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=567,
  serialized_end=665,
)

_DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION = _descriptor.Descriptor(
  name='ResourceDescription',
  full_name='access.DescriptionResponse.ResourceDescription',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='kind', full_name='access.DescriptionResponse.ResourceDescription.kind', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='comment', full_name='access.DescriptionResponse.ResourceDescription.comment', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='actions', full_name='access.DescriptionResponse.ResourceDescription.actions', index=2,
      number=3, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='roles', full_name='access.DescriptionResponse.ResourceDescription.roles', index=3,
      number=4, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[_DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ACTION, _DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ROLE, _DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ACTIONSENTRY, _DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ROLESENTRY, ],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=178,
  serialized_end=665,
)

_DESCRIPTIONRESPONSE = _descriptor.Descriptor(
  name='DescriptionResponse',
  full_name='access.DescriptionResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='resources', full_name='access.DescriptionResponse.resources', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[_DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION, ],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=86,
  serialized_end=665,
)


_PERMITTEDACTIONSREQUEST = _descriptor.Descriptor(
  name='PermittedActionsRequest',
  full_name='access.PermittedActionsRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='resource_kind', full_name='access.PermittedActionsRequest.resource_kind', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='resource_ids', full_name='access.PermittedActionsRequest.resource_ids', index=1,
      number=2, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=667,
  serialized_end=737,
)


_PERMITTEDACTIONSRESPONSE_RESOURCEPERMISSIONS = _descriptor.Descriptor(
  name='ResourcePermissions',
  full_name='access.PermittedActionsResponse.ResourcePermissions',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='actions', full_name='access.PermittedActionsResponse.ResourcePermissions.actions', index=0,
      number=1, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=890,
  serialized_end=928,
)

_PERMITTEDACTIONSRESPONSE_PERMITTEDENTRY = _descriptor.Descriptor(
  name='PermittedEntry',
  full_name='access.PermittedActionsResponse.PermittedEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='access.PermittedActionsResponse.PermittedEntry.key', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='value', full_name='access.PermittedActionsResponse.PermittedEntry.value', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=_descriptor._ParseOptions(descriptor_pb2.MessageOptions(), _b('8\001')),
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=930,
  serialized_end=1032,
)

_PERMITTEDACTIONSRESPONSE = _descriptor.Descriptor(
  name='PermittedActionsResponse',
  full_name='access.PermittedActionsResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='permitted', full_name='access.PermittedActionsResponse.permitted', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='validity_duration', full_name='access.PermittedActionsResponse.validity_duration', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[_PERMITTEDACTIONSRESPONSE_RESOURCEPERMISSIONS, _PERMITTEDACTIONSRESPONSE_PERMITTEDENTRY, ],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=740,
  serialized_end=1032,
)

_DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ACTION.containing_type = _DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION
_DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ROLE.containing_type = _DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION
_DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ACTIONSENTRY.fields_by_name['value'].message_type = _DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ACTION
_DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ACTIONSENTRY.containing_type = _DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION
_DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ROLESENTRY.fields_by_name['value'].message_type = _DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ROLE
_DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ROLESENTRY.containing_type = _DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION
_DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION.fields_by_name['actions'].message_type = _DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ACTIONSENTRY
_DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION.fields_by_name['roles'].message_type = _DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ROLESENTRY
_DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION.containing_type = _DESCRIPTIONRESPONSE
_DESCRIPTIONRESPONSE.fields_by_name['resources'].message_type = _DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION
_PERMITTEDACTIONSRESPONSE_RESOURCEPERMISSIONS.containing_type = _PERMITTEDACTIONSRESPONSE
_PERMITTEDACTIONSRESPONSE_PERMITTEDENTRY.fields_by_name['value'].message_type = _PERMITTEDACTIONSRESPONSE_RESOURCEPERMISSIONS
_PERMITTEDACTIONSRESPONSE_PERMITTEDENTRY.containing_type = _PERMITTEDACTIONSRESPONSE
_PERMITTEDACTIONSRESPONSE.fields_by_name['permitted'].message_type = _PERMITTEDACTIONSRESPONSE_PERMITTEDENTRY
_PERMITTEDACTIONSRESPONSE.fields_by_name['validity_duration'].message_type = google_dot_protobuf_dot_duration__pb2._DURATION
DESCRIPTOR.message_types_by_name['DescriptionResponse'] = _DESCRIPTIONRESPONSE
DESCRIPTOR.message_types_by_name['PermittedActionsRequest'] = _PERMITTEDACTIONSREQUEST
DESCRIPTOR.message_types_by_name['PermittedActionsResponse'] = _PERMITTEDACTIONSRESPONSE

DescriptionResponse = _reflection.GeneratedProtocolMessageType('DescriptionResponse', (_message.Message,), dict(

  ResourceDescription = _reflection.GeneratedProtocolMessageType('ResourceDescription', (_message.Message,), dict(

    Action = _reflection.GeneratedProtocolMessageType('Action', (_message.Message,), dict(
      DESCRIPTOR = _DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ACTION,
      __module__ = 'access_pb2'
      # @@protoc_insertion_point(class_scope:access.DescriptionResponse.ResourceDescription.Action)
      ))
    ,

    Role = _reflection.GeneratedProtocolMessageType('Role', (_message.Message,), dict(
      DESCRIPTOR = _DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ROLE,
      __module__ = 'access_pb2'
      # @@protoc_insertion_point(class_scope:access.DescriptionResponse.ResourceDescription.Role)
      ))
    ,

    ActionsEntry = _reflection.GeneratedProtocolMessageType('ActionsEntry', (_message.Message,), dict(
      DESCRIPTOR = _DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ACTIONSENTRY,
      __module__ = 'access_pb2'
      # @@protoc_insertion_point(class_scope:access.DescriptionResponse.ResourceDescription.ActionsEntry)
      ))
    ,

    RolesEntry = _reflection.GeneratedProtocolMessageType('RolesEntry', (_message.Message,), dict(
      DESCRIPTOR = _DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ROLESENTRY,
      __module__ = 'access_pb2'
      # @@protoc_insertion_point(class_scope:access.DescriptionResponse.ResourceDescription.RolesEntry)
      ))
    ,
    DESCRIPTOR = _DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION,
    __module__ = 'access_pb2'
    # @@protoc_insertion_point(class_scope:access.DescriptionResponse.ResourceDescription)
    ))
  ,
  DESCRIPTOR = _DESCRIPTIONRESPONSE,
  __module__ = 'access_pb2'
  # @@protoc_insertion_point(class_scope:access.DescriptionResponse)
  ))
_sym_db.RegisterMessage(DescriptionResponse)
_sym_db.RegisterMessage(DescriptionResponse.ResourceDescription)
_sym_db.RegisterMessage(DescriptionResponse.ResourceDescription.Action)
_sym_db.RegisterMessage(DescriptionResponse.ResourceDescription.Role)
_sym_db.RegisterMessage(DescriptionResponse.ResourceDescription.ActionsEntry)
_sym_db.RegisterMessage(DescriptionResponse.ResourceDescription.RolesEntry)

PermittedActionsRequest = _reflection.GeneratedProtocolMessageType('PermittedActionsRequest', (_message.Message,), dict(
  DESCRIPTOR = _PERMITTEDACTIONSREQUEST,
  __module__ = 'access_pb2'
  # @@protoc_insertion_point(class_scope:access.PermittedActionsRequest)
  ))
_sym_db.RegisterMessage(PermittedActionsRequest)

PermittedActionsResponse = _reflection.GeneratedProtocolMessageType('PermittedActionsResponse', (_message.Message,), dict(

  ResourcePermissions = _reflection.GeneratedProtocolMessageType('ResourcePermissions', (_message.Message,), dict(
    DESCRIPTOR = _PERMITTEDACTIONSRESPONSE_RESOURCEPERMISSIONS,
    __module__ = 'access_pb2'
    # @@protoc_insertion_point(class_scope:access.PermittedActionsResponse.ResourcePermissions)
    ))
  ,

  PermittedEntry = _reflection.GeneratedProtocolMessageType('PermittedEntry', (_message.Message,), dict(
    DESCRIPTOR = _PERMITTEDACTIONSRESPONSE_PERMITTEDENTRY,
    __module__ = 'access_pb2'
    # @@protoc_insertion_point(class_scope:access.PermittedActionsResponse.PermittedEntry)
    ))
  ,
  DESCRIPTOR = _PERMITTEDACTIONSRESPONSE,
  __module__ = 'access_pb2'
  # @@protoc_insertion_point(class_scope:access.PermittedActionsResponse)
  ))
_sym_db.RegisterMessage(PermittedActionsResponse)
_sym_db.RegisterMessage(PermittedActionsResponse.ResourcePermissions)
_sym_db.RegisterMessage(PermittedActionsResponse.PermittedEntry)


_DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ACTIONSENTRY.has_options = True
_DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ACTIONSENTRY._options = _descriptor._ParseOptions(descriptor_pb2.MessageOptions(), _b('8\001'))
_DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ROLESENTRY.has_options = True
_DESCRIPTIONRESPONSE_RESOURCEDESCRIPTION_ROLESENTRY._options = _descriptor._ParseOptions(descriptor_pb2.MessageOptions(), _b('8\001'))
_PERMITTEDACTIONSRESPONSE_PERMITTEDENTRY.has_options = True
_PERMITTEDACTIONSRESPONSE_PERMITTEDENTRY._options = _descriptor._ParseOptions(descriptor_pb2.MessageOptions(), _b('8\001'))
# @@protoc_insertion_point(module_scope)
