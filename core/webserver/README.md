# UniFi Poller: `webserver` plugin

Built-In Web Server Go Library for UniFi Poller.

Enabling the web server is optional. It provides a window into the running data.
The web server may be secured with a simple password. SSL is also optional.

See the [Web Server Wiki](https://github.com/unifi-poller/unifi-poller/wiki/Web-Server)
for more information about how it works.

Other plugins must import this library to make use of it. While this library is
labeled as a plugin, it's pretty much required since everything imports it.
That said, it is still disabled by default, and won't store any data unless it's
enabled.

_This needs a better godoc and examples._

## Overview

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
