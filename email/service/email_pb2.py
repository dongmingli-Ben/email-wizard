# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: email.proto
"""Generated protocol buffer code."""
from google.protobuf.internal import builder as _builder
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x0b\x65mail.proto\x12\x05\x65mail\"/\n\x0c\x45mailRequest\x12\x0e\n\x06\x63onfig\x18\x01 \x01(\t\x12\x0f\n\x07n_mails\x18\x02 \x01(\x05\"\x1d\n\nEmailReply\x12\x0f\n\x07message\x18\x01 \x01(\t2D\n\x0b\x45mailHelper\x12\x35\n\tGetEmails\x12\x13.email.EmailRequest\x1a\x11.email.EmailReply\"\x00\x42#\n\rio.grpc.emailB\nEmailProtoP\x01\xa2\x02\x03HLWb\x06proto3')

_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, globals())
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'email_pb2', globals())
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'\n\rio.grpc.emailB\nEmailProtoP\001\242\002\003HLW'
  _EMAILREQUEST._serialized_start=22
  _EMAILREQUEST._serialized_end=69
  _EMAILREPLY._serialized_start=71
  _EMAILREPLY._serialized_end=100
  _EMAILHELPER._serialized_start=102
  _EMAILHELPER._serialized_end=170
# @@protoc_insertion_point(module_scope)
