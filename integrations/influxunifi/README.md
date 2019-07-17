<img width="320px" src="https://raw.githubusercontent.com/wiki/davidnewhall/unifi-poller/images/unifi-poller-logo.png">


[![discord](https://badgen.net/badge/icon/Discord?color=0011ff&label&icon=https://simpleicons.now.sh/discord/eee "Captain's Discord")](https://discord.gg/DyVsMyt)
[![twitter](https://badgen.net/twitter/follow/TwitchCaptain?icon=https://simpleicons.now.sh/twitter/0099ff&label=TwitchCaptain&color=0116ff "TwitchCaptain @ Twitter")](https://twitter.com/TwitchCaptain)
[![grafana](https://badgen.net/https/code.golift.io/bd/grafana/dashboard-downloads/10414,10415,10416,10417,10418?icon=https://simpleicons.now.sh/grafana/ED7F38&color=0011ff "Grafana Dashboard Downloads")](http://grafana.com/dashboards?search=unifi-poller)
[![pulls](https://badgen.net/docker/pulls/golift/unifi-poller?icon=https://simpleicons.now.sh/docker/38B1ED&label=pulls&color=0011ff "Docker Pulls")](https://hub.docker.com/r/golift/unifi-poller)
[![DLs](https://img.shields.io/github/downloads/davidnewhall/unifi-poller/total.svg?logo=github&color=0116ff "GitHub Downloads")](https://www.somsubhra.com/github-release-stats/?username=davidnewhall&repository=unifi-poller)

[![unifi](https://badgen.net/badge/UniFi/5.10.x,5.11.x,UAP,USG,USW?list=|&icon=https://golift.io/svg/ubiquiti_color.svg&color=0099ee "UniFi Products Supported")](https://github.com/golift/unifi)
[![builer](https://badgen.net/badge/go/Application%20Builder?label=&icon=https://golift.io/svg/go.svg&color=0099ee "Go Application Builder")](https://github.com/golift/application-builder)
[![stars](https://badgen.net/github/stars/davidnewhall/unifi-poller?icon=https://simpleicons.now.sh/macys/fab&label=&color=0099ee "GitHub Stars")](https://github.com/davidnewhall/unifi-poller)
[![travis](https://badgen.net/travis/davidnewhall/unifi-poller?icon=travis&label=build "Travis Build")](https://travis-ci.org/davidnewhall/unifi-poller)

Collect your UniFi controller data and send it to an InfluxDB instance.
[Grafana Dashboards](http://grafana.com/dashboards?search=unifi-poller) included.
Updated 2019.

## Description

[Ubiquiti](https://www.ui.com) makes networking devices like switches, gateways
(routers) and wireless access points. They have a line of equipment named
[UniFi](https://www.ui.com/products/#unifi) that uses a
[controller](https://www.ui.com/download/unifi/) to keep stats and simplify network
device configuration. This controller can be installed on Windows, macOS and Linux.
Ubiquiti also provides a dedicated hardware device called a
[CloudKey](https://www.ui.com/unifi/unifi-cloud-key/) that runs the controller software.

UniFi Poller is a small Golang application that runs on Windows, macOS, Linux or
Docker. It polls a UniFi controller every 30 seconds for measurements and stores
the data in an Influx database. A small setup with 2 access points, 1 switch, 1
gateway and 40 clients produces over 3000 fields (metrics).

This application requires your controller to be running all the time. If you run
a UniFi controller, there's no excuse not to install
[Influx](https://github.com/davidnewhall/unifi-poller/wiki/InfluxDB),
[Grafana](https://github.com/davidnewhall/unifi-poller/wiki/Grafana) and this app.
You'll have a plethora of data at your fingertips and the ability to craft custom
graphs to slice the data any way you choose. Good luck!

## Installation

[See the Wiki!](https://github.com/davidnewhall/unifi-poller/wiki/Installation)
We have a special place for [Docker Users](https://github.com/davidnewhall/unifi-poller/wiki/Docker).

# Backstory

Okay, so here's the deal. I found a simple piece of code on GitHub that
sorta did what I needed; we all know that story. I wanted more data, so
I added more data collection. I believe I've completely rewritten every
piece of original code, except the copyright/license file and that's fine
with me. I probably wouldn't have made it this far if
[Garrett](https://github.com/dewski/unifi) hadn't written the original
code I started with. Many props my man.

The original code pulled only the client data. This app now pulls data
for clients, access points, security gateways and switches. I currently
own two UAP-AC-PROs, one USG-3 and one US-24-250W. If your devices differ
this app may miss some data. I'm willing to help and make it better.
Open an [Issue](https://github.com/davidnewhall/unifi-poller/issues) and
we'll figure out how to get things working for you.

# What's this data good for?

I've been trying to get my UAP data into Grafana. Sure, google search that.
You'll find [this](https://community.ubnt.com/t5/UniFi-Wireless/Grafana-dashboard-for-UniFi-APs-now-available/td-p/1833532).
And that's all you'll find. What if you don't want to deal with SNMP?
Well, here you go. I've replicated 90% of what you see on those SNMP-powered
dashboards with this Go app running on the same mac as my UniFi controller.
All without enabling SNMP nor trying to understand those OIDs. Mad props
to [waterside](https://community.ubnt.com/t5/user/viewprofilepage/user-id/303058)
for making this dashboard; it gave me a fantastic start to making my own.

I've also created [another forum post](https://community.ui.com/questions/Unifi-Poller-Store-Unifi-Controller-Metrics-in-InfluxDB-without-SNMP/58a0ea34-d2b3-41cd-93bb-d95d3896d1a1) you may use to get additional help.

# Development

The "What now..." section below used to be a lot larger. I've received a lot of
support, feedback and assistance from the community. Many thanks! This app is
extremely stable with a tiny memory and cpu footprint. I imagine one day we'll
figure out how to make it run on a CloudKey device directly; once I have one
personally that will be my goal. In addition to stability, this app provides
an intuitive installation and configuration process. Maintenance is a breeze too.

I'm not a software engineer, I'm a a firm believer in operational excellence above
all else. To that end, this app shall remain easy, intuitive and highly adaptable.
I'm totally open to add more configuration options if someone raises a need or concern.

You can control this app with puppet, chef, saltstack, homebrew or a simple bash
script if you needed to. It's available for macOS, Linux and Docker. It comes with
a systemd service unit that allows you automatically start it up on most Linux
hosts. It works just fine on [Windows](https://github.com/davidnewhall/unifi-poller/wiki/Windows) too.

The unifi data extraction is provided as an [external library](https://godoc.org/github.com/golift/unifi),
and you can import that code directly without futzing with this application. That
means, if you wanted to do something like make telegraf collect your data instead
of UniFi Poller you can achieve that with a little bit of Go code. You could write
a small app that acts as a telegraf input plugin using the [unifi](https://github.com/golift/unifi)
library to grab the data from your controller. As a bonus, all of the code in UniFi Poller is
[also a library](https://godoc.org/github.com/davidnewhall/unifi-poller/unifipoller)
and can be used in other projects.

# What now...

### Are there other devices that need to be included?

I have: switch, router, access point. Three total, and the type structs are
likely missing data for variants of these devices. e.g. Some UAPs have more
radios, I probably didn't properly account for that. Some gateways have more
ports, some switches have 10Gb, etc. These are things I do not have data on
to write code for. If you have these devices, and want them graphed, open an
Issue and lets discuss.

### Radios, Frequencies, Interfaces, vAPs

My access points only seem to have two radios, one interface and vAP per radio.
I'm not sure if the graphs, as-is, provide enough insight into APs with other
configurations. Help me figure that out?

# What's it look like?

Here's a picture of the Client dashboard.
![UniFi Clients Dashboard Image](https://grafana.com/api/dashboards/10418/images/6554/image)

Here's a picture of the USG dashboard.
![USG Dashboard Image](https://grafana.com/api/dashboards/10416/images/6552/image)

Here's a picture of the UAP dashboard. This only shows one device, but you can
select multiple to put specific stats side-by-side.
![UAP Dashboard Image](https://grafana.com/api/dashboards/10415/images/6551/image)

The USW / Switch Dashboard is pretty big with one data-filled section per selected port.
You can drill down into specific sites, switches, and ports. Compare ports in different
sites side-by-side. So easy! This screenshot barely does it justice.
![USW Dashboard Image](https://grafana.com/api/dashboards/10417/images/6553/image)


## Copyright & License
<img style="float: right;" align="right" width="200px" src="https://raw.githubusercontent.com/wiki/davidnewhall/unifi-poller/images/unifi-poller-logo.png">

-   Copyright © 2016 Garrett Bjerkhoel.
-   Copyright © 2018-2019 David Newhall II.
-   See [LICENSE](LICENSE) for license information.
