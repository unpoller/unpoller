#!/bin/bash

# This file is used by osxpkg packages.  FPM use.

# Copy the config file into place if it does not exist.
if [ ! -f /usr/local/etc/unifi-poller/up.conf ] && [ -f /usr/local/etc/unifi-poller/up.conf.example ]; then
  cp /usr/local/etc/unifi-poller/up.conf.example /usr/local/etc/unifi-poller/up.conf
fi

# Make sure admins can write logs.
chgrp admin /usr/local/var/log
chmod g=rwx /usr/local/var/log

# This starts it as root. no no no .... not sure how to fix that.
# launchctl load /Library/LaunchAgents/com.github.davidnewhall.unifi-poller.plist
