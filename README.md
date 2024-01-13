# Hlog

Hlog let you aggregate in one place logs from an unlimited number of processes almost instantly, perfom monitoring on these aggregated data in real time, detect patterns in the process and run callbacks when specific events happen.


## Features

Hlog stands out with many unique features that it supports natively, amoung them:

- `Optimized transport` with a custom-made protocol on top of TCP.
- `Multiple steps logging` are transaction-oriented logs.
- `Actions` are callback functions that can be called when tracked events happen.
- `Advanced queries` on your logs and metadata about them.
- `Fast IO` enabled by technologies like Clickhouse.

## Usage

Use-cases are presented using the python programming language.

```python
# Simple one time logging
logger = hlog()
data = {"foo": "bar"}
options = {"tags": ["greeting"]}

logger.info("hello, from the hlog team", data, options)
```

```python
# Transaction-based logging
logger = hlog()
tx = logger()
tx.start()

tx.info("Wake up...", 1)
tx.info("Do hard work...", 2)
tx.info("Go sleep...", 3)

tx.end()
```