unifi-poller(1) -- Utility to poll UniFi Controller Metrics and store them in InfluxDB
===

SYNOPSIS
---
`unifi-poller -c /etc/unifi-poller.conf`

This daemon polls a UniFi controller at a short interval and stores the collected
measurements in an Influx Database. The measurements and metrics collected belong
to every available site, device and client found on the controller. Including
UniFi Security Gateways, Access Points, Switches and possibly more.

Dashboards for Grafana are available.
Find them at [Grafana.com](https://grafana.com/dashboards?search=unifi-poller).

DESCRIPTION
---
Unifi-Poller is a small Golang application that runs on Windows, macOS, Linux or
Docker. It polls a UniFi controller every 30 seconds for measurements and stores
the data in an Influx database. See the example configuration file for more
examples and default configurations.

*   See the example configuration file for more examples and default configurations.

OPTIONS
---
`unifi-poller [-c <config-file>] [-j <filter>] [-h] [-v]`

    -c, --config <config-file>
        Provide a configuration file (instead of the default).

    -v, --version
        Display version and exit.

    -j, --dumpjson <filter>
        This is a debug option; use this when you are missing data in your graphs,
        and/or you want to inspect the raw data coming from the controller. The
        filter accepts three options: devices, clients, other. This will print a
        lot of information. Recommend piping it into a file and/or into jq for
        better visualization. This requires a valid config file that contains
        working authentication details for a UniFi Controller. This only dumps
        data for sites listed in the config file. The application exits after
        printing the JSON payload; it does not daemonize or report to InfluxDB
        with this option. The `other` option is special. This allows you request
        any api path. It must be enclosed in quotes with the word other. Example:
           unifi-poller -j "other /stat/admins"

    -h, --help
        Display usage and exit.

CONFIGURATION
---
*   Config File Default Location: `/etc/unifi-poller/up.conf`
*   Config File Default Format: `TOML`
*   Possible formats: `XML`, `JSON`, `TOML`, `YAML`

The config file can be written in four different syntax formats. The application
decides which one to used based on the file's name. If it contains `.xml` it will
be parsed as XML. The same goes for `.json` and `.yaml`. If the filename contains
none of these strings, then it is parsed as the default format, TOML. This option
is provided so the application can be easily adapted to any environment.

`Config File Parameters`

    sites          default: ["all"]
        This list of strings should represent the names of sites on the UniFi
        controller that will be polled for data. Pass `all` in the list to
        poll all sites. On startup, the application prints out all site names
        found in the controller; they're cryptic, but they have the human-name
        next to them. The cryptic names go into the config file `sites` list.
        The controller's first site is not cryptic and is named `default`.

    interval       default: 30s
        How often to poll the controller for updated client and device data.
        The UniFi Controller only updates traffic stats about every 30 seconds.

    debug          default: false
        This turns on time stamps and line numbers in logs, outputs a few extra
        lines of information while processing.

    quiet          default: false  
        Setting this to true will turn off per-device and per-interval logs. Only
        errors will be logged. Using this with debug=true adds line numbers to
        any error logs.

    max_errors     default: 0
        If you restart the UniFI controller, the poller will lose access until
        it is restarted. Specifying a number greater than -1 for max_errors will
        cause the poller to exit when it reaches the error count specified.
        This problematic condition can be triggered by InfluxDB having issues
        too. Generally only 1 error per interval is created, but if more than one
        backend is having issues > 1 error could be generated per interval. Once
        the poller exits, it is expected that something will restart it
        automatically so it gets back in line; something is usually systemd,
        docker or launchd. The default setting of 0 will cause an exit after
        just 1 error. Recommended values are 0-5.

    influx_url     default: http://127.0.0.1:8086
        This is the URL where the Influx web server is available.

    influx_user    default: unifi
        Username used to authenticate with InfluxDB.

    influx_pass    default: unifi
        Password used to authenticate with InfluxDB.

    influx_db      default: unifi
        Custom database created in InfluxDB to use with this application.
        On first setup, log into InfluxDB and create access:
        $ influx -host localhost -port 8086
        CREATE DATABASE unifi
        CREATE USER unifi WITH PASSWORD 'unifi' WITH ALL PRIVILEGES
        GRANT ALL ON unifi TO unifi

    unifi_url      default: https://127.0.0.1:8443
        This is the URL where the UniFi Controller is available.

    unifi_user     default: influxdb
        Username used to authenticate with UniFi controller. This should be a
        special service account created on the control with read-only access.

    unifi_user     no default   ENV: UNIFI_PASSWORD
        Password used to authenticate with UniFi controller. This can also be
        set in an environment variable instead of a configuration file.

    verify_ssl     default: false
        If your UniFi controller has a valid SSL certificate, you can enable
        this option to validate it. Otherwise, any SSL certificate is valid.

GO DURATION
---
This application uses the Go Time Durations for a polling interval.
The format is an integer followed by a time unit. You may append
multiple time units to add them together. A few valid time units are:

    ms   (millisecond)
    s    (second)
    m    (minute)

Example Use: `35s`, `1m`, `1m30s`

AUTHOR
---
*   Garrett Bjerkhoel (original code) ~ 2016
*   David Newhall II (rewritten) ~ 4/20/2018
*   David Newhall II (still going) ~ 6/7/2019

LOCATION
---
*   UniFi Poller: [https://github.com/davidnewhall/unifi-poller](https://github.com/davidnewhall/unifi-poller)
*   UniFi Library: [https://github.com/golift/unifi](https://github.com/golift/unifi)
