unpoller(1) -- Utility to poll UniFi Controller Metrics and store them in InfluxDB
===

SYNOPSIS
---
`unpoller -c /etc/unpoller/up.conf`

This daemon polls a UniFi controller at a short interval and stores the collected
measurements in an Influx Database. The measurements and metrics collected belong
to every available site, device and client found on the controller. Including
UniFi Security Gateways, Access Points, Switches and possibly more.

Dashboards for Grafana are available.
Find them at [Grafana.com](https://grafana.com/dashboards?search=unpoller).

DESCRIPTION
---
UniFi Poller is a small Golang application that runs on Windows, macOS, Linux or
Docker. It polls a UniFi controller every 30 seconds for measurements and stores
the data in an Influx database. See the example configuration file for more
examples and default configurations.

*   See the example configuration file for more examples and default configurations.

OPTIONS
---
`unpoller [-c <config-file>,[config-file]] [-j <filter>] [-e <pass>] [-h] [-v]`

    -c, --config <config-file>,[config-file]
        Provide a configuration file (instead of the default). You may provide
        multiple file paths separated by commas. The first file found is used.

    -e, --encrypt <password|->
        Encrypts a password with bcrypt. This is useful to make an account
        password for the built-in web server. You may provide the password
        on the CLI, or use "-" as the pass to be prompted for the password
        without echo.

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
           unpoller -j "other /stat/admins"

    -h, --help
        Display usage and exit.

CONFIGURATION
---
*   Config File Default Location:
    *   Mac/BSD: `/usr/local/etc/unpoller/up.conf`
    *   Windows: `C:\\ProgramData\\unpoller\\up.conf`
    *   Others:  `/etc/unpoller/up.conf`
*   Config File Default Format: `TOML`
*   Possible formats: `JSON`, `TOML`, `YAML`

The config file can be written in four different syntax formats. The application
decides which one to use based on the file's name. If it contains `.json` it will
be parsed as JSON. The same goes for `.yaml`. If the filename contains neither
of these strings, then it is parsed as the default format, TOML. This feature
is provided so the application can be easily adapted to any environment. `XML`
used to work but no longer does.

`Config File Parameters`

Configuration file (up.conf) parameters are documented in the wiki.

*   [https://github.com/unpoller/unpoller/wiki/Configuration](https://github.com/unpoller/unpoller/wiki/Configuration)

`Shell Environment Parameters`

This application can be fully configured using shell environment variables.
Find documentation for this feature on the Docker Wiki page, and the above Configuration wiki.

*   [https://github.com/unpoller/unpoller/wiki/Docker](https://github.com/unpoller/unpoller/wiki/Docker)

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
*   David Newhall II (still going) ~ 2/1/2020

LOCATION
---
*   UniFi Poller: [https://golift.io/unpoller/](https://golift.io/unpoller/)
*   UniFi Library: [https://github.com/unpoller/unifi](https://github.com/unpoller/unifi)
*   Grafana Dashboards: [https://grafana.com/dashboards?search=unpoller](https://grafana.com/dashboards?search=unpoller)
