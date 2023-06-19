from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class EmailContentRequest(_message.Message):
    __slots__ = ["additional_info", "email"]
    ADDITIONAL_INFO_FIELD_NUMBER: _ClassVar[int]
    EMAIL_FIELD_NUMBER: _ClassVar[int]
    additional_info: str
    email: str
    def __init__(self, email: _Optional[str] = ..., additional_info: _Optional[str] = ...) -> None: ...

class EmailParseReply(_message.Message):
    __slots__ = ["message"]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    message: str
    def __init__(self, message: _Optional[str] = ...) -> None: ...
