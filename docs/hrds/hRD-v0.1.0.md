# hRD-v0.1.0

Here is a visual representation of the proposed architecture for this release.

![Architecture](/docs/hrds/assets/hRD-v0.1.0-architecture.jpg)

This release should include an end-to-end experience with the following characteristics:

1. A set of **collectors** for the languages Python, TS/JS, Go and Rust that are fully tested which support the **transport** Kafka..
2. A **core** composed of the MongoDB **ingester**, the MongoDB **sinker** and the MongoDB **watcher**, fully tested and benchmarked. The kernel is expected to work with single workers for each component.
3. A fully working and tested implementation of the **exporter** optimised to work with the observer (using HTTP and web-sockets).
4. A version of the **observer** that supports monitoring, visualising and managing the system.

### Tasks

- [ ]  Clearly define the collector interface and data formats
- [ ]  Correct the Python collector
- [ ]  Correct the JS/TS collector
- [ ]  Write the Go collector
- [ ]  Write the Rust collector
- [ ]  Refactor the MongoDB ingester
- [ ]  Test the MongoDB ingester
- [ ]  Benchmark the MongoDB ingester
- [ ]  Decouple the MongoDB sinker from the ingester
- [ ]  Refactor the core to well define the kernel
- [ ]  Refactor the core to correctly define the exporter and its implementation
- [ ]  Conceptualise the watching process
    - [ ]  Define and implement hlogevent
    - [ ]  Test the hlogevent implementation
    - [ ]  Define and implement the MongoDB watcher
    - [ ]  Test the MongoDB watcher
    - [ ]  Benchmark the MongoDB watcher
- [ ]  Finish the observer to support the release proposal.