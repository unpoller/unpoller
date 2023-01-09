<img width="320px" src="https://unpoller.com/img/logo.png">

[![discord](https://badgen.net/badge/icon/Discord?color=0011ff&label&icon=https://simpleicons.now.sh/discord/eee "GoLift Discord")](https://golift.io/discord)
[![grafana](https://badgen.net/https/golift.io/bd/grafana/dashboard-downloads/11310,10419,10414,10415,10416,10417,10418,11311,11312,11313,11314,11315?icon=https://simpleicons.now.sh/grafana/ED7F38&color=0011ff "Grafana Dashboard Downloads")](http://grafana.com/dashboards?search=unifi-poller)
[![pulls](https://badgen.net/docker/pulls/golift/unifi-poller?icon=https://simpleicons.now.sh/docker/38B1ED&label=pulls&color=0011ff "Docker Pulls")](https://hub.docker.com/r/golift/unifi-poller)
[![stars](https://badgen.net/github/stars/unifi-poller/unifi-poller?icon=https://simpleicons.now.sh/macys/fab&label=&color=0099ee "GitHub Stars")](https://github.com/unpoller/unpoller)

[![unifi](https://badgen.net/badge/UniFi/5.12.x,5.13.x,UAP,USG,USW,UDM?list=|&icon=https://docs.golift.io/svg/ubiquiti_color.svg&color=0099ee "UniFi Products Supported")](https://github.com/golift/unifi)

Collect your UniFi controller data and report it to an InfluxDB instance,
or export it for Prometheus collection.
[Twelve Grafana Dashboards](http://grafana.com/dashboards?search=unifi-poller)
included; with screenshots. Six for InfluxDB and six for Prometheus.

## Installation

[See the Documentation!](https://unpoller.com)
We're willing to help if you have troubles.
Open an [Issue](https://github.com/unpoller/unpoller/issues) and
we'll figure out how to get things working for you. You can also get help in
the #unpoller channel on the [GoLift Discord server](https://golift.io/discord). There is also
[a forum post](https://community.ui.com/questions/Unifi-Poller-Store-Unifi-Controller-Metrics-in-InfluxDB-without-SNMP/58a0ea34-d2b3-41cd-93bb-d95d3896d1a1)
you may use to get additional help.

## Description

[Ubiquiti](https://www.ui.com) makes networking devices like switches, gateways
(routers) and wireless access points. They have a line of equipment named
[UniFi](https://www.ui.com/products/#unifi) that uses a
[controller](https://www.ui.com/download/unifi/) to keep stats and simplify network
device configuration. This controller can be installed on Windows, macOS, FreeBSD,
Linux or Docker. Ubiquiti also provides a dedicated hardware device called a
[CloudKey](https://www.ui.com/unifi/unifi-cloud-key/) that runs the controller software.
More recently they've developed the Dream Machine, and UnPoller can collect its data!

UnPoller is a small Golang application that runs on Windows, macOS, FreeBSD,
Linux or Docker. In Influx-mode it polls a UniFi controller every 30 seconds for
measurements and exports the data to an Influx database. In Prometheus mode the
poller opens a web port and accepts Prometheus polling. It converts the UniFi
Controller API data into Prometheus exports on the fly.

This application requires your controller to be running all the time. If you run
a UniFi controller, there's no excuse not to install
Influx or
[Prometheus](https://prometheus.io),
Grafana and this app.
You'll have a plethora of data at your fingertips and the ability to craft custom
graphs to slice the data any way you choose. Good luck!

Supported as of Poller v2.0.2, are [Loki](https://grafana.com/oss/loki/)
and the collection of UniFi events, alarms, anomalies and IDS data.
This data can be exported to Loki or InfluxDB, or both!

## Operation

You can control this app with puppet, chef, saltstack, homebrew or a simple bash
script if you needed to. Packages are available for macOS, Linux, FreeBSD and Docker.
It works just fine on Windows too.

## What does it look like?

There are 12 total dashboards available; the 6 InfluxDB dashboards are very similar
to the 6 Prometheus dashboards. On the [documentation website](https://unpoller.com)
you'll find screenshots of some of the dashboards.

## Integrations

The following fine folks are providing their services, completely free! These service
integrations are used for things like storage, building, compiling, distribution and
documentation support. This project succeeds because of them. Thank you!

<p style="text-align: center;">
<a title="PackageCloud" alt="PackageCloud" href="https://packagecloud.io"><img src="https://docs.golift.io/integrations/packagecloud.png"/></a>
<a title="GitHub" alt="GitHub" href="https://GitHub.com"><img src="https://docs.golift.io/integrations/octocat.png"/></a>
<a title="Docker Cloud" alt="Docker" href="https://cloud.docker.com"><img src="https://docs.golift.io/integrations/docker.png"/></a>
<a title="Homebrew" alt="Homebrew" href="https://brew.sh"><img src="https://docs.golift.io/integrations/homebrew.png"/></a>
<a title="Go Lift" alt="Go Lift" href="https://golift.io"><img src="https://docs.golift.io/integrations/golift.png"/></a>
<a title="Grafana" alt="Grafana" href="https://grafana.com"><img src="https://docs.golift.io/integrations/grafana.png"/></a>
</p>

## Copyright & License

<img style="float: right;" align="right" width="200px" src="https://unpoller.com/img/logo.png">

- Copyright © 2018-2020 David Newhall II.
- See [LICENSE](LICENSE) for license information.
