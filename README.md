# Hlog

Hlog let you aggregate in one place logs from an unlimited number of processes almost instantly, perfom monitoring on these aggregated data in real time, detect patterns in the process and run callbacks when specific events happen.


## Features

Hlog stands out with many unique features that it supports natively, amoung them:

- `Multiple steps logging` are transaction-oriented logs.
- `Actions` are callback functions that can be called when tracked events happen.
- `Advanced queries` on your logs and metadata about them.

## Components

The Hlog software comprises the following elements:

- The `core` which is the engine that keeps everything together.
- The storage layer which is made of `Clickhouse` and `Postgres`.
- The transport layer which comprises `HTTP` and `Kafka`.
- The observation layer which comprises a CLI that is read-only, a web application that allows for querying.
- The clients which are connectors (in `Go`, `Python`, and `TypeScript`/`JavaScript`) between applications and the log aggregator.
