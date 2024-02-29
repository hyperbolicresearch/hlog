package ingest

// ingest implements the flow from Kafka to storages. It consists of
// two parts, which are:
//  1. Ingest to MongoDB.
//
//     Ingest into MongDB is the default real-time ingestion pipeline.
//     As soon as a new log gets added to Kafka, we consume it and
//     proceed with MongoDB, which notifies us on inserting in order for
//     us to notify the clients. Like this, we have near instant notifications
//     and real-time log monitoring.
//
//  2. Ingest to ClickHouse.
//
//     Ingestion into ClickHouse is for analytical purposes only. It is
//     done periodically and the period is parameterizable. Upon inserting
//     (into MongoDB), we ingest the logs to Kafka where they get accumulated
//     for the <period> amount of time. Ingestion into ClickHouse happens
//     by batches which allows to update metrics and query results.
//
//  Flow
//  Client -> Kafka -> MongoDB -> Client & Kafka -> ClickHouse <- Client
