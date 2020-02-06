<img width="320px" src="https://raw.githubusercontent.com/wiki/unifi-poller/unifi-poller/images/unifi-poller-logo.png">

[![discord](https://badgen.net/badge/icon/Discord?color=0011ff&label&icon=https://simpleicons.now.sh/discord/eee "Ubiquiti Discord")](https://discord.gg/KnyKYt2)
[![twitter](https://badgen.net/twitter/follow/TwitchCaptain?icon=https://simpleicons.now.sh/twitter/0099ff&label=TwitchCaptain&color=0116ff "TwitchCaptain @ Twitter")](https://twitter.com/TwitchCaptain)
[![grafana](https://badgen.net/https/golift.io/bd/grafana/dashboard-downloads/11310,10419,10414,10415,10416,10417,10418,11311,11312,11313,11314,11315?icon=https://simpleicons.now.sh/grafana/ED7F38&color=0011ff "Grafana Dashboard Downloads")](http://grafana.com/dashboards?search=unifi-poller)
[![pulls](https://badgen.net/docker/pulls/golift/unifi-poller?icon=https://simpleicons.now.sh/docker/38B1ED&label=pulls&color=0011ff "Docker Pulls")](https://hub.docker.com/r/golift/unifi-poller)
[![DLs](https://img.shields.io/github/downloads/unifi-poller/unifi-poller/total.svg?logo=github&color=0116ff "GitHub Downloads")](https://www.somsubhra.com/github-release-stats/?username=unifi-poller&repository=unifi-poller)

[![unifi](https://badgen.net/badge/UniFi/5.12.x,5.13.x,UAP,USG,USW,UDM?list=|&icon=https://docs.golift.io/svg/ubiquiti_color.svg&color=0099ee "UniFi Products Supported")](https://github.com/golift/unifi)
[![builer](https://badgen.net/badge/go/Application%20Builder?label=&icon=https://docs.golift.io/svg/go.svg&color=0099ee "Go Application Builder")](https://github.com/golift/application-builder)
[![stars](https://badgen.net/github/stars/unifi-poller/unifi-poller?icon=https://simpleicons.now.sh/macys/fab&label=&color=0099ee "GitHub Stars")](https://github.com/unifi-poller/unifi-poller)
[![travis](https://badgen.net/travis/unifi-poller/unifi-poller?icon=travis&label=build "Travis Build")](https://travis-ci.org/unifi-poller/unifi-poller)

Collect your UniFi controller data and report it to an InfluxDB instance,
or export it for Prometheus collection.
[Twelve Grafana Dashboards](http://grafana.com/dashboards?search=unifi-poller)
included; with screenshots. Six for InfluxDB and six for Prometheus.

## Installation

[See the Wiki!](https://github.com/unifi-poller/unifi-poller/wiki/Installation)
We have a special place for [Docker Users](https://github.com/unifi-poller/unifi-poller/wiki/Docker).
I'm willing to help if you have troubles.
Open an [Issue](https://github.com/unifi-poller/unifi-poller/issues) and
we'll figure out how to get things working for you. You can also get help in
the #unifi-poller channel on the [Ubiquiti Discord server](https://discord.gg/KnyKYt2). I've also
[provided a forum post](https://community.ui.com/questions/Unifi-Poller-Store-Unifi-Controller-Metrics-in-InfluxDB-without-SNMP/58a0ea34-d2b3-41cd-93bb-d95d3896d1a1)
you may use to get additional help.

## Description

[Ubiquiti](https://www.ui.com) makes networking devices like switches, gateways
(routers) and wireless access points. They have a line of equipment named
[UniFi](https://www.ui.com/products/#unifi) that uses a
[controller](https://www.ui.com/download/unifi/) to keep stats and simplify network
device configuration. This controller can be installed on Windows, macOS, FreeBSD,
Linux or Docker. Ubiquiti also provides a dedicated hardware device called a
[CloudKey](https://www.ui.com/unifi/unifi-cloud-key/) that runs the controller software.
More recently they've developed the Dream Machine; it's still in
beta / early access, but UniFi Poller can collect its data!

UniFi Poller is a small Golang application that runs on Windows, macOS, FreeBSD,
Linux or Docker. In Influx-mode it polls a UniFi controller every 30 seconds for
measurements and exports the data to an Influx database. In Prometheus mode the
poller opens a web port and accepts Prometheus polling. It converts the UniFi
Controller API data into Prometheus exports on the fly.

This application requires your controller to be running all the time. If you run
a UniFi controller, there's no excuse not to install
[Influx](https://github.com/unifi-poller/unifi-poller/wiki/InfluxDB) or
[Prometheus](https://prometheus.io),
[Grafana](https://github.com/unifi-poller/unifi-poller/wiki/Grafana) and this app.
You'll have a plethora of data at your fingertips and the ability to craft custom
graphs to slice the data any way you choose. Good luck!

## Backstory

I found a simple piece of code on GitHub that sorta did what I needed;
we all know that story. I wanted more data, so I added more data collection.
I probably wouldn't have made it this far if [Garrett](https://github.com/dewski/unifi)
hadn't written the original code I started with. Many props my man.
The original code pulled only the client data. This app now pulls data
for clients, access points, security gateways, dream machines and switches.

I've been trying to get my UAP data into Grafana. Sure, google search that.
You'll find [this](https://community.ubnt.com/t5/UniFi-Wireless/Grafana-dashboard-for-UniFi-APs-now-available/td-p/1833532).
What if you don't want to deal with SNMP?
Well, here you go. I've replicated 400% of what you see on those SNMP-powered
dashboards with this Go app running on the same mac as my UniFi controller.
All without enabling SNMP nor trying to understand those OIDs. Mad props
to [waterside](https://community.ubnt.com/t5/user/viewprofilepage/user-id/303058)
for making this dashboard; it gave me a fantastic start to making my own dashboards.

## Operation

You can control this app with puppet, chef, saltstack, homebrew or a simple bash
script if you needed to. Packages are available for macOS, Linux, FreeBSD and Docker.
It works just fine on [Windows](https://github.com/unifi-poller/unifi-poller/wiki/Windows) too.
Most people prefer Docker, and this app is right at home in that environment.

## What's it look like?

There are 12 total dashboards available; the 6 InfluxDB dashboards are very similar
to the 6 Prometheus dashboards. Below you'll find screenshots of the first four dashboards.

##### Client Dashboard (InfluxDB)
![UniFi Clients Dashboard Image](https://grafana.com/api/dashboards/10418/images/7540/image)

##### USG Dashboard (InfluxDB)
![USG Dashboard Image](https://grafana.com/api/dashboards/10416/images/7543/image)

##### UAP Dashboard (InfluxDB)
![UAP Dashboard Image](https://grafana.com/api/dashboards/10415/images/7542/image)

##### USW / Switch Dashboard (InfluxDB)
You can drill down into specific sites, switches, and ports. Compare ports in different
sites side-by-side. So easy! This screenshot barely does it justice.
![USW Dashboard Image](https://grafana.com/api/dashboards/10417/images/7544/image)

## Integrations

The following fine folks are providing their services, completely free! These service
integrations are used for things like storage, building, compiling, distribution and
documentation support. This project succeeds because of them. Thank you!

<p style="text-align: center;">
<a title="Jfrog Bintray" alt="Jfrog Bintray" href="https://bintray.com"><img src="https://docs.golift.io/integrations/bintray.png"/></a>
<a title="GitHub" alt="GitHub" href="https://GitHub.com"><img src="https://docs.golift.io/integrations/octocat.png"/></a>
<a title="Docker Cloud" alt="Docker" href="https://cloud.docker.com"><img src="https://docs.golift.io/integrations/docker.png"/></a>
<a title="Travis-CI" alt="Travis-CI" href="https://Travis-CI.com"><img src="https://docs.golift.io/integrations/travis-ci.png"/></a>
<a title="Homebrew" alt="Homebrew" href="https://brew.sh"><img src="https://docs.golift.io/integrations/homebrew.png"/></a>
<a title="Go Lift" alt="Go Lift" href="https://golift.io"><img src="https://docs.golift.io/integrations/golift.png"/></a>
<a title="Grafana" alt="Grafana" href="https://grafana.com"><img src="https://docs.golift.io/integrations/grafana.png"/></a>
</p>

## Copyright & License
<img style="float: right;" align="right" width="200px" src="https://raw.githubusercontent.com/wiki/unifi-poller/unifi-poller/images/unifi-poller-logo.png">

-   Copyright © 2018-2020 David Newhall II.
-   See [LICENSE](LICENSE) for license information.
