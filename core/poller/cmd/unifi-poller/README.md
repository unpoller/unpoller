unifi-poller(1) -- Utility to poll Unifi Metrics and drop them into InfluxDB
===

## SYNOPSIS

`unifi-poller -c /usr/local/etc/unifi-poller.conf`

## DESCRIPTION

* This application polls a Unifi Controller API for Client and Device Metrics.
* The metrics are then stored in an InfluxDB instance.

## OPTIONS

`unifi-poller [-c <config file>] [-h] [-v]`

    -c, --config <file_path>
        Provide a configuration file (instead of the default).

    -v, --version
        Display version and exit.

    -h, --help
        Display usage and exit.


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
* previously: https://github.com/dewski/unifi
