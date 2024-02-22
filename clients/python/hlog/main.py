# pyright: reportMissingImports=false

"""

"""

from dataclasses import dataclass
from datetime import datetime
from enum import Enum
import json
from typing import Any, Dict, TypeVar
import uuid
from confluent_kafka import Producer


def serialize_value(value):
    return json.dumps(value).encode('utf-8')


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

V = TypeVar('V')

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


class Hlog:
    def __init__(self, configs: Config) -> None:
        self.client_id = configs.client_id
        self.default_level = configs.default_level
        self.kafka_topic = configs.channel_id
        self.kafka_config = {
            "bootstrap.servers": configs.kafka_server,
            "client.id": configs.client_id,
        }
        self.producer = Producer(self.kafka_config)

    def _publish(
        self,
        message: str,
        data: Any = None,
        level=LogLevel.DEBUG
    ) -> bool:
        msg = Log(
            log_id=str(uuid.uuid4()),
            sender_id=self.client_id,
            timestamp=int(datetime.now().replace(tzinfo=None).timestamp()),
            level=str(level),
            message=message,
            data=data,
        )
        self.producer.produce(
            self.kafka_topic,
            value=json.dumps(msg.to_dict()),
        )
        self.producer.poll(1)
        # FIXME implement the actual behavior based on the asyncronous
        # confirmation of the log.
        return True

    def debug(self, message: str, data: Any = None) -> bool:
        return self._publish(message, data, LogLevel.DEBUG)
    
    def info(self, message: str, data: Any = None) -> bool:
        return self._publish(message, data, LogLevel.INFO)
    
    def warn(self, message: str, data: Any = None) -> bool:
        return self._publish(message, data, LogLevel.WARN)
    
    def error(self, message: str, data: Any = None) -> bool:
        return self._publish(message, data, LogLevel.ERROR)
    
    def fatal(self, message: str, data: Any = None) -> bool:
        return self._publish(message, data, LogLevel.FATAL)
    
    def __repr__(self) -> str:
        return f"<Hlog client_id={self.client_id}>"
