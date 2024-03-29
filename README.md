![Web look](https://github.com/hyperbolicresearch/hlog/blob/dev/assets/github_img.png)

> This project is under massive development and can undergo changes in the conceptual framework and implementation. We are just conducting experimentations, just thinking loud until otherwise said.

## Overview

**hlog** is a highly flexible, extensible and design-wise minimalistic data aggregation framework built to be the skeleton of a diverse computing ecosystem. It is designed to have different key characteristics, such as:

- `Extensibility`: As a minimalistic framework, it intends to be highly extensible by providing different levels of interfaces for the access and manipulation of the data and its processings.
- `Flexibility`: One of the promise of hlog is to enforce composability between its components, providing great abstractions so that it can avoid high coupling with any particular technology or tool.
- `Compatibility`: : hlog is not meant to be a brand new tool that will live in its own world. It should leverage existing community standards and evolve with them.
- `Performance`: By the ambiguous word performance, one should understand a system that is resilient (fault-tolerant), efficient in term of resource utilisation, and highly scalable.
- `Cloud-native`: It is thought upfront to be easily deployable in cloud environments.


## Core components

The stack is composed by three theoretical core components:

1. **collector**: A collector interfaces with the data source (applications, infrastructures, platforms, …), collect and send the data to the system.
2. **core**: The core is the main part of the system and consists of a set of programs responsible for the ingestion, the processing, and the storage of the data. It’s managed by the **kernel**.
3. **exporter**: An exporter is an interface that exposes the data, its processing and its transformation to the external world.

## Accessory components

Although useful, those are not part of the core components, but nonetheless, they constitute a set of tools that come handy:

1. **observer**: The observer is a web-based UI that gives access to the system and allow for monitoring, querying, visualising, experiencing and managing the system.
2. **hlogcli**: A console-based version of the observer.

## Docs
- [Glossary](/docs/dictionary.md)
- [hlog Release Descriptions (hRDs)](/docs/hrds/)
- [hlog Improvements Proposals (hIPs)](/docs/hips/)
