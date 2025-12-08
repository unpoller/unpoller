## Docker Compose

The files in this folder can be used locally to spin up
a full set of applications (minus the UniFi controller) to get
UniFi Poller up and running. Including InfluxDB, Grafana, and
Chronograph. This last app is useful to inspect the data stored
in InfluxDB by UniFi Poller.

##### HOWTO
**Learn more about how and when to use these *Docker Compose* files in the
[Docker Wiki](https://unpoller.com/docs/install/dockercompose).**

## Health Check

The UniFi Poller Docker image includes a built-in health check that validates
the configuration and checks plugin connectivity. The health check runs every
30 seconds and marks the container as unhealthy if configuration issues are
detected or if enabled outputs cannot be reached.

You can manually run the health check:
```bash
docker exec <container_name> /usr/bin/unpoller --health
```

The health check is automatically used by Docker and container orchestration
platforms (Kubernetes, Docker Swarm, etc.) to determine container health status.
