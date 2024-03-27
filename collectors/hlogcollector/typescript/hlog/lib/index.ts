/** Implementation of Hlog client for TS */

import { Config, Log, LogLevel } from "./types";
import { KafkaConfig, Producer, Kafka } from "kafkajs"
import { v4 as uuidv4 } from "uuid";

class Hlog {
  client_id: string;
  default_level: LogLevel;
  kafka_topic: string;
  kafka_config: KafkaConfig
  producer: Producer

  constructor(config: Config) {
    this.client_id = config.client_id;
    this.default_level = config.default_level || LogLevel.DEBUG;
    this.kafka_topic = config.channel_id;
    this.kafka_config = {
      brokers: [config.kafka_server],
      clientId: config.client_id,
    };
    this.producer = new Kafka(this.kafka_config).producer();
  }
  
  async publish(message: string, data: any, level: LogLevel) {
    const msg: Log = {
      log_id: uuidv4(),
      sender_id: this.client_id,
      timestamp: Date.now(),
      level: level,
      message: message,
      data: data
    }
    await this.producer.connect()
    await this.producer.send({
      topic: this.kafka_topic,
      messages: [
        {value: JSON.stringify(msg)}
      ]
    })

    // FIXME: implement the actual behavior of this
    return true;
  }

  async debug(message: string, data: any) {
    return this.publish(message, data, LogLevel.DEBUG);
  }

  async info(message: string, data: any) {
    return this.publish(message, data, LogLevel.INFO);
  }

  async warn(message: string, data: any) {
    return this.publish(message, data, LogLevel.WARN);
  }

  async error(message: string, data: any) {
    return this.publish(message, data, LogLevel.ERROR);
  }

  async fatal(message: string, data: any) {
    return this.publish(message, data, LogLevel.FATAL);
  }
}

export default Hlog;
