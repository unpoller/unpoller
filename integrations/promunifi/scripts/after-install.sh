#!/bin/bash

# This file is used by deb, rpm and osx packages.
# FPM adds this as the after-install script.

if [ "$(uname -s)" = "Darwin" ]; then
  # Copy the config file into place if it does not exist.
  if [ ! -f /usr/local/etc/unifi-poller/up.conf ] && [ -f /usr/local/etc/unifi-poller/up.conf.example ]; then
    cp /usr/local/etc/unifi-poller/up.conf.example /usr/local/etc/unifi-poller/up.conf
  fi

  # Allow admins to change the configuration and write logs.
  chgrp -R admin /usr/local/etc/unifi-poller
  chmod -R g+wr /usr/local/etc/unifi-poller

  # Make sure admins can write logs.
  chgrp admin /usr/local/var/log
  chmod g=rwx /usr/local/var/log

  # This starts it as root. no no no .... not sure how to fix that.
  # launchctl load /Library/LaunchAgents/com.github.davidnewhall.unifi-poller.plist

elif [ -x "/bin/systemctl" ]; then
  /bin/systemctl daemon-reload
  /bin/systemctl enable unifi-poller
  /bin/systemctl restart unifi-poller
fi
