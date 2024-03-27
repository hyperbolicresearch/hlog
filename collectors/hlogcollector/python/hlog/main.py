# pyright: reportMissingImports=false

"""
The python client package for hlog
"""

import json
import uuid
from datetime import datetime
from typing import Any
from confluent_kafka import Producer
from .types import Config, Log, LogLevel


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
            level=level.value,
            message=message,
            data=data,
        )
        self.producer.produce(
            self.kafka_topic,
            value=json.dumps(msg.to_dict()),
        )
        self.producer.poll(1)
        # FIXME implement the actual behavior based on the asynchronous
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
