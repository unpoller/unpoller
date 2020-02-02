## Docker Cloud Builds

This folder contains the files that build our Docker image. The image
is built by Docker Hub "automatically" using the [Dockerfile](Dockerfile)
and [hooks/](hooks/) in this folder.

## Docker Compose

The other files in this folder can be used locally to spin up
a full set of applications (minus the UniFi controller) to get
UniFi Poller up and running. Including InfluxDB, Grafana, and
Chronograph. This last app is useful to inspect the data stored
in InfluxDB by UniFi Poller.

##### HOWTO
**Learn more about how and when to use these *Docker Compose* files in the
[Docker Wiki](https://github.com/unifi-poller/unifi-poller/wiki/Docker).**
