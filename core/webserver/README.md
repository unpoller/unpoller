# webserver

Built-In Web Server Go Library for UniFi Poller. **INCOMPLETE**

Enabling the web server is optional. It provides a window into the running data.
The web server may be secured with a simple password. SSL is also optional.

## Interface

-   Recent logs from poller are visible.
-   Uptime and Version are displayed across the top.

### Controllers

-   The web server interface allows you to see the configuration for each controller.
-   Some meta data about each controller is displayed, such as sites, clients and devices.
-   Example config: [up.json.example](https://github.com/unifi-poller/unifi-poller/blob/master/examples/up.json.example)

### Input Plugins

-   You may view input plugin configuration. Currently only UniFi.
-   The example config above shows input plugin data.

### Output Plugins

-   You may view output plugin configuration. Currently Prometheus and InfluxDB.
-   The example config above shows output plugin data.
