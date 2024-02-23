"use strict";
/** Implementation of Hlog client for TS */
Object.defineProperty(exports, "__esModule", { value: true });
const types_1 = require("./types");
const kafkajs_1 = require("kafkajs");
const uuid_1 = require("uuid");
class Hlog {
    constructor(config) {
        this.client_id = config.client_id;
        this.default_level = config.default_level || types_1.LogLevel.DEBUG;
        this.kafka_topic = config.channel_id;
        this.kafka_config = {
            brokers: [config.kafka_server],
            clientId: config.client_id,
        };
        this.producer = new kafkajs_1.Kafka(this.kafka_config).producer();
    }
    async publish(message, data, level) {
        const msg = {
            log_id: (0, uuid_1.v4)(),
            sender_id: this.client_id,
            timestamp: Date.now(),
            level: level,
            message: message,
            data: data
        };
        await this.producer.connect();
        await this.producer.send({
            topic: this.kafka_topic,
            messages: [
                { value: JSON.stringify(msg) }
            ]
        });
        // FIXME: implement the actual behavior of this
        return true;
    }
    async debug(message, data) {
        return this.publish(message, data, types_1.LogLevel.DEBUG);
    }
    async info(message, data) {
        return this.publish(message, data, types_1.LogLevel.INFO);
    }
    async warn(message, data) {
        return this.publish(message, data, types_1.LogLevel.WARN);
    }
    async error(message, data) {
        return this.publish(message, data, types_1.LogLevel.ERROR);
    }
    async fatal(message, data) {
        return this.publish(message, data, types_1.LogLevel.FATAL);
    }
}
const hlog = new Hlog({
    client_id: "test",
    kafka_server: "0.0.0.0:65007",
    channel_id: "subnet",
});
const res = hlog.debug("hello world!", { foo: "bar" });
console.log(res);
exports.default = Hlog;
