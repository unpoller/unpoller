<img width="320px" src="https://raw.githubusercontent.com/wiki/davidnewhall/unifi-poller/images/unifi-poller-logo.png">

[![discord](https://badgen.net/badge/icon/Discord?color=0011ff&label&icon=https://simpleicons.now.sh/discord/eee "Ubiquiti Discord")](https://discord.gg/KnyKYt2)
[![twitter](https://badgen.net/twitter/follow/TwitchCaptain?icon=https://simpleicons.now.sh/twitter/0099ff&label=TwitchCaptain&color=0116ff "TwitchCaptain @ Twitter")](https://twitter.com/TwitchCaptain)
[![grafana](https://badgen.net/https/golift.io/bd/grafana/dashboard-downloads/10414,10415,10416,10417,10418,11311,11312,11313,11314,11315?icon=https://simpleicons.now.sh/grafana/ED7F38&color=0011ff "Grafana Dashboard Downloads")](http://grafana.com/dashboards?search=unifi-poller)
[![pulls](https://badgen.net/docker/pulls/golift/unifi-poller?icon=https://simpleicons.now.sh/docker/38B1ED&label=pulls&color=0011ff "Docker Pulls")](https://hub.docker.com/r/golift/unifi-poller)
[![DLs](https://img.shields.io/github/downloads/davidnewhall/unifi-poller/total.svg?logo=github&color=0116ff "GitHub Downloads")](https://www.somsubhra.com/github-release-stats/?username=davidnewhall&repository=unifi-poller)

[![unifi](https://badgen.net/badge/UniFi/5.11.x,5.12.x,UAP,USG,USW,UDM?list=|&icon=https://docs.golift.io/svg/ubiquiti_color.svg&color=0099ee "UniFi Products Supported")](https://github.com/golift/unifi)
[![builer](https://badgen.net/badge/go/Application%20Builder?label=&icon=https://docs.golift.io/svg/go.svg&color=0099ee "Go Application Builder")](https://github.com/golift/application-builder)
[![stars](https://badgen.net/github/stars/davidnewhall/unifi-poller?icon=https://simpleicons.now.sh/macys/fab&label=&color=0099ee "GitHub Stars")](https://github.com/davidnewhall/unifi-poller)
[![travis](https://badgen.net/travis/davidnewhall/unifi-poller?icon=travis&label=build "Travis Build")](https://travis-ci.org/davidnewhall/unifi-poller)

Collect your UniFi controller data and report it to an InfluxDB instance,
or export it for Prometheus collection. Prometheus support is
[new](https://github.com/davidnewhall/unifi-poller/issues/88), and much
of the documentation still needs to be updated; 12/2/2019.
[Ten Grafana Dashboards](http://grafana.com/dashboards?search=unifi-poller)
included; with screenshots. Five for InfluxDB and five for Prometheus.

## Installation
[See the Wiki!](https://github.com/davidnewhall/unifi-poller/wiki/Installation)
We have a special place for [Docker Users](https://github.com/davidnewhall/unifi-poller/wiki/Docker).
I'm willing to help if you have troubles.
Open an [Issue](https://github.com/davidnewhall/unifi-poller/issues) and
we'll figure out how to get things working for you. You can also get help in
the #unifi-poller channel on the [Ubiquiti Discord server](https://discord.gg/KnyKYt2).
I've also [provided a forum post](https://community.ui.com/questions/Unifi-Poller-Store-Unifi-Controller-Metrics-in-InfluxDB-without-SNMP/58a0ea34-d2b3-41cd-93bb-d95d3896d1a1) you may use to get additional help.

## Description
[Ubiquiti](https://www.ui.com) makes networking devices like switches, gateways
(routers) and wireless access points. They have a line of equipment named
[UniFi](https://www.ui.com/products/#unifi) that uses a
[controller](https://www.ui.com/download/unifi/) to keep stats and simplify network
device configuration. This controller can be installed on Windows, macOS and Linux.
Ubiquiti also provides a dedicated hardware device called a
[CloudKey](https://www.ui.com/unifi/unifi-cloud-key/) that runs the controller software. More recently they've developed the Dream Machine; it's still in
beta / early access, but UniFi Poller can collect its data!

UniFi Poller is a small Golang application that runs on Windows, macOS, Linux or
Docker. In Influx-mode it polls a UniFi controller every 30 seconds for
measurements and exports the data to an Influx database. In Prometheus mode the
poller opens a web port and accepts Prometheus polling. It converts the UniFi
Controller API data into Prometheus exports on the fly.

This application requires your controller to be running all the time. If you run
a UniFi controller, there's no excuse not to install
[Influx](https://github.com/davidnewhall/unifi-poller/wiki/InfluxDB) or
[Prometheus](https://prometheus.io),
[Grafana](https://github.com/davidnewhall/unifi-poller/wiki/Grafana) and this app.
You'll have a plethora of data at your fingertips and the ability to craft custom
graphs to slice the data any way you choose. Good luck!

## Backstory
I found a simple piece of code on GitHub that sorta did what I needed;
we all know that story. I wanted more data, so I added more data collection.
I believe I've completely rewritten every piece of original code, except the
copyright/license file and that's fine with me. I probably wouldn't have made
it this far if [Garrett](https://github.com/dewski/unifi) hadn't written the
original code I started with. Many props my man.

The original code pulled only the client data. This app now pulls data
for clients, access points, security gateways, dream machines and switches. I
used to own two UAP-AC-PROs, one USG-3 and one US-24-250W, but have since upgraded
a few devices. Many other users have also provided feedback to improve this app,
and we have reports of it working on nearly every switch, AP and gateway.

## What's this data good for?
I've been trying to get my UAP data into Grafana. Sure, google search that.
You'll find [this](https://community.ubnt.com/t5/UniFi-Wireless/Grafana-dashboard-for-UniFi-APs-now-available/td-p/1833532). What if you don't want to deal with SNMP?
Well, here you go. I've replicated 400% of what you see on those SNMP-powered
dashboards with this Go app running on the same mac as my UniFi controller.
All without enabling SNMP nor trying to understand those OIDs. Mad props
to [waterside](https://community.ubnt.com/t5/user/viewprofilepage/user-id/303058)
for making this dashboard; it gave me a fantastic start to making my own dashboards.

## Operation
You can control this app with puppet, chef, saltstack, homebrew or a simple bash
script if you needed to. Packages are available for macOS, Linux and Docker.
It comes with a systemd service unit that allows you automatically start it up on most Linux hosts.
It works just fine on [Windows](https://github.com/davidnewhall/unifi-poller/wiki/Windows) too.
Most people prefer Docker, and this app is right at home in that environment.

## Development
The UniFi data extraction is provided as an [external library](https://godoc.org/golift.io/unifi),
and you can import that code directly without futzing with this application. That
means, if you wanted to do something like make telegraf collect your data instead
of UniFi Poller you can achieve that with a little bit of Go code. You could write
a small app that acts as a telegraf input plugin using the [unifi](https://github.com/golift/unifi)
library to grab the data from your controller. As a bonus, all of the code in UniFi Poller is
[in libraries](https://godoc.org/github.com/davidnewhall/unifi-poller/pkg)
and can be used in other projects.

## What's it look like?

There are five total dashboards available. Below you'll find screenshots of a few.

##### Client Dashboard (InfluxDB)
![UniFi Clients Dashboard Image](https://grafana.com/api/dashboards/10418/images/6660/image)

##### USG Dashboard (InfluxDB)
![USG Dashboard Image](https://grafana.com/api/dashboards/10416/images/6663/image)

##### UAP Dashboard (InfluxDB)
![UAP Dashboard Image](https://grafana.com/api/dashboards/10415/images/6662/image)

##### USW / Switch Dashboard (InfluxDB)
You can drill down into specific sites, switches, and ports. Compare ports in different
sites side-by-side. So easy! This screenshot barely does it justice.
![USW Dashboard Image](https://grafana.com/api/dashboards/10417/images/6664/image)

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
<img style="float: right;" align="right" width="200px" src="https://raw.githubusercontent.com/wiki/davidnewhall/unifi-poller/images/unifi-poller-logo.png">

-   Copyright © 2016 Garrett Bjerkhoel.
-   Copyright © 2018-2019 David Newhall II.
-   See [LICENSE](LICENSE) for license information.
