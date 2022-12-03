# poller

## UniFi Poller Core

This module ties the inputs together with the outputs.

Aggregates metrics on request. Provides CLI app and args parsing.

## Ideal

This library has no notion of "UniFi" or controllers, or Influx, or Prometheus.
This library simply provides an input interface and an output interface.
Each interface uses an `[]any` type, so any type of data can be used.
That is to say, you could write input and output plugins that work with, say,
Cisco gear, or any other network (or even non-network) data. The existing plugins
should provide ample example of how to use this library, but at some point the
godoc will improve.

## Features

- Automatically unmarshal's plugin config structs from config file and/or env variables.
- Initializes all "imported" plugins on startup.
- Provides input plugins a Logger, requires an interface for Metrics and Events retrieval.
- Provides Output plugins an interface to retrieve Metrics and Events, and a Logger.
- Provides automatic aggregation of Metrics and Events from multiple sources.
