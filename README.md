![Web look](https://github.com/hyperbolicresearch/hlog/blob/dev/assets/github_img.png)

> This project is under massive development and can undergo changes in the conceptual framework and implementation. We are just conducting experimentations, just thinking loud until otherwise said.

hlog is an observability platform focused on logs and metrics aggregation built to be the base layer of a whole computing ecosystem. To achieve performance, it's designed with several key characteristics:

1. `Extensibility` meaning that you can build on top of it by leveraging the exposed API.
2. `Performance` by being fault-tolerant, highly available and scalable.
3. `Flexibility` by enforcing the composability of the different components.
4. `Compatibility` by integrating with existing tools.
5. `Cloud-native` through containerization, support for configuration settings, etc...

## Components

The hlog stack is composed by three components:

- `collector` is the interface that interacts with applications and infrastructures and sending the collected logs and metrics to the system.
- `hlog` is the main server which processes and stores the logs and metrics, and also exposing the API for interaction with them.
- `observer` as the user interface.
