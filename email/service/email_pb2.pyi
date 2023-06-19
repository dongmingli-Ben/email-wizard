from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class EmailReply(_message.Message):
    __slots__ = ["message"]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    message: str
    def __init__(self, message: _Optional[str] = ...) -> None: ...

class EmailRequest(_message.Message):
    __slots__ = ["config", "n_mails"]
    CONFIG_FIELD_NUMBER: _ClassVar[int]
    N_MAILS_FIELD_NUMBER: _ClassVar[int]
    config: str
    n_mails: int
    def __init__(self, config: _Optional[str] = ..., n_mails: _Optional[int] = ...) -> None: ...
