from dataclasses import dataclass
from enum import Enum
from typing import Any


class LogLevel(Enum):
    DEBUG = "DEBUG"
    INFO = "INFO"
    WARN = "WARN"
    ERROR = "ERROR"
    FATAL = "FATAL"


@dataclass(frozen=True)
class Config:
    client_id: str
    kafka_server: str
    kafka_username: str = ""
    kafka_password: str = ""
    channel_id: str = "default"
    default_level: LogLevel = LogLevel.DEBUG


@dataclass(frozen=True)
class Log:
    log_id: str
    sender_id: str
    timestamp: int
    level: str
    message: str
    data: Any

    def to_dict(self):
        return {k: v for k, v in self.__dict__.items()}
