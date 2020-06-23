# webserver

Built-In Web Server Go Library for UniFi Poller. **INCOMPLETE**

Enabling the web server is optional. It provides a window into the running data.
The web server may be secured with a simple password. SSL is also optional.

## Interface

-   Recent logs from poller are visible.
-   Uptime and Version are displayed across the top.

### Controllers

-   The web server interface allows you to see the configuration for each controller.
-   You may select a controller, and then select a site on the controller.
-   Some meta data about each controller is displayed, such as sites, clients and devices.
-   Example config: [up.json.example](https://github.com/unifi-poller/unifi-poller/blob/master/examples/up.json.example)

### Input Plugins

-   You may view input plugin configuration. Currently only UniFi.
-   The example config above shows input plugin data.

### Output Plugins

-   You may view output plugin configuration. Currently Prometheus and InfluxDB.
-   The example config above shows output plugin data.

### Sites

Each controller has 1 or more sites. Most people only have 1, but some enterprises
run this software and have many more. Each site has devices like switches (`USW`), access
points (`UAP`) and routers (`USG`/`UDM`). We'll have counts for each device type. Each device has a name.
Each device has a count of clients (access points have clients). We'll want to expose
this, but it's not in a useful format yet. It'll look something like what you see below,
but keep the visualization expandable. We may add "model" and "serial number" for each device.
There is a handful of meta data per device. Some users have hundreds of devices.
```
{
  "site_name_here": {
    "clients": 22,
    "UAP": [
      {
        "name": "ap1-room",
        "clients": 6
      },
      {
        "name": "ap2-bran",
        "clients": 6
      }
    ],
    "USW": [
      {
        "name": "sw1-cube",
        "model": "US-500w-P",
        "serial": "xyz637sjs999",
        "clients": 7
      },
      {
        "name": "sw2-trap",
        "clients": 3
      }
    ],
    "USG": [
      {
        "name": "gw1-role",
        "clients": 22
      }
    ]
  }
}
