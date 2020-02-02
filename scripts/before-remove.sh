#!/bin/sh

# This file is used by txz, rpm and deb packages. FPM use.

if [ "$1" = "upgrade" ] || [ "$1" = "1" ] ; then
  exit 0
fi

if [ -x "/bin/systemctl" ]; then
  /bin/systemctl stop unifi-poller
  /bin/systemctl disable unifi-poller
elif [ -x /usr/sbin/service ]; then
  /usr/sbin/service unifi-poller stop
  /usr/sbin/service unifi-poller disable
fi
