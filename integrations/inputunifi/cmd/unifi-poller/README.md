unifi-poller(1) -- Utility to poll Unifi Metrics and drop them into InfluxDB
===

## SYNOPSIS

`unifi-poller -c /usr/local/etc/unifi-poller.conf`

## DESCRIPTION

* This application polls a Unifi Controller API for Client and Device Metrics.
* The metrics are then stored in an InfluxDB instance.
* See the example configuration file for help configuring this application.

## OPTIONS

`unifi-poller [-c <config-file>] [-h] [-v]`

    -c, --config <config-file>
        Provide a configuration file (instead of the default).

    -v, --version
        Display version and exit.

    -h, --help
        Display usage and exit.

## CONFIGURATION

* Config File Default Location: /usr/local/etc/unifi-poller/up.conf

`Config File Parameters`

    `interval`       default: 30s
        How often to poll the controller for updated client and device data.
        The Unifi Controller only updates traffic stats about every 30 seconds.

    `debug`          default: false
        This turns on time stamps and line numbers in logs, outputs a few extra
        lines of information while processing.

    `quiet`          default: false  
        Setting this to true will turn off per-device and per-interval logs. Only
        errors will be logged. Using this with debug=true adds line numbers to
        any error logs.

    `influx_url`     default: http://127.0.0.1:8086
        This is the URL where the Influx web server is available.

    `influx_user`    default: unifi
        Username used to authenticate with InfluxDB. Many servers do not use auth.

    `influx_pass`    default: unifi
        Password used to authenticate with InfluxDB.

    `influx_db`      default: unifi
        Custom database created in InfluxDB to use with this application.

    `unifi_url`      default: https://127.0.0.1:8443
        This is the URL where the Unifi Controller is available.

    `unifi_user`     default: influxdb
        Username used to authenticate with Unifi controller. This should be a
        special service account created on the control with read-only access.

    `unifi_user`     no default   ENV: UNIFI_PASSWORD
        Password used to authenticate with Unifi controller. This can also be
        set in an environment variable instead of a configuration file.

    `verify_ssl`     default: false
        If your Unifi controller has a valid SSL certificate, you can enable
        this option to validate it. Otherwise, any SSL certificate is valid.

## GO DURATION

This application uses the Go Time Durations for a polling interval.
The format is an integer followed by a time unit. You may append
multiple time units to add them together. Some valid time units are:

     `us` (microsecond)
     `ns` (nanosecond)
     `ms` (millisecond)
     `s`  (second)
     `m`  (minute)
     `h`  (hour)

Example Use: `1m`, `5h`, `100ms`, `17s`, `1s45ms`, `1m3s`

## AUTHOR

* Garrett Bjerkhoel (original code) ~ 2016
* David Newhall II (rewritten) ~ 4/20/2018

## LOCATION

* https://github.com/davidnewhall/unifi-poller
* /usr/local/bin/unifi-poller
* config-file: /usr/local/etc/unifi-poller/up.conf
* previously: https://github.com/dewski/unifi
