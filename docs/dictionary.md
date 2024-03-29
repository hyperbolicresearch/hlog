# Incomplete and possibly ambiguous dictionary of hlog-specific concepts

> *To ensure a good communication and a standardised approach to the system, we define the following with a degree of precision that is variable and subject to evolution with time.*

### collector

`collector` is the interface which collects data from the data source. Examples / implementations of collectors can be (by example) the Go or the Python collectors.

### core

`core` is the main part of the system which consists of a set of programs responsible for the ingestion, the processing and the storage of the data. The core has the following parts:

1. The **kernel**
2. The **ingesters**
3. **The sinkers**
4. The **watchers**

### exporter

`exporter` is the interface that exposes the data, their processing and there transformation to the external world as well as the hlogevents. An implementation of the exporter can be a web API.

### hlog

`hlog` is a highly flexible, extensible and design-wise minimalistic data aggregation framework built to be the skeleton of a diverse computing ecosystem.

### hlogcli

`hlogcli` is a CLI accessory component that exposes the capabilities for monitoring, querying, visualising and managing the system.

### hlogevent

`hlogevent` is a data structure that is emitted when a pattern takes place.

### ingester

`ingester` is the interface the consumes the data collected by the collectors.

### kernel

`kernel` is the central manager of the core. It orchestrates the work of the different core components and ensures an efficient utilisation of the resources.

### logger

`logger` is the writer of the system. 

> In the current way to do things, the logger makes the explicit writes ; a work that is supposed to be handled by the exposer, but as it seems to me, the logger can be used by the exposer to do the work, but more efficiently, by example to prevent back-pressure.
> 

### observable

`observable` are trackable entities that can be watched. They can be defined by simple database queries, aggregations, transactions, sequences of queries, …

### observer

`observer` is a web-based UI (accessory component) that provides the capabilities for monitoring, querying, visualising and managing the system.

### pattern

`pattern` is a condition that can be verified on an observable.

### plugin

`plugin` is not a part of hlog properly speaking, but is a separate service built on top of the API exposed by the exporter.

### sinker

`sinker` is a component of the hlog’s core that is responsible to sink the data to a database. Examples / implementations of sinkers can be (by example) the MongoDB sinker or the ClickHouse batcher.

### storage

`storage` refers to where the data is stored.

### transport

`transport` refers the mean by which the collectors send data to the system.

### watcher

`watcher` is a core component that monitors the observables and emits hlogevents when patterns occurs.