package ingest

// ingest implements the flow from Kafka to ClickHouse. It consists of
// two parts, which are:
//  1. The ingester consumes from Kafka, transforms the data to
//     the structure that will be stored in clickHouse, extracts
//     the schema from the batch of data it buffers and generate SQL.
//  2. The batcher is responsible to accumulate the data processed by
//     the ingester and to batch them to ClickHouse.
//
// The flow is then:
// kafka -> ingester -> batcher -> clickhouse.