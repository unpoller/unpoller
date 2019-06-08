#!/bin/bash

# This file is used by deb, rpm and osx packages.
# FPM adds this as the after-install script.

if [ "$(uname -s)" = "Darwin" ]; then
  # Copy the config file into place if it does not exist.
  if [ ! -f /usr/local/etc/unifi-poller/up.conf ] && [ -f /usr/local/etc/unifi-poller/up.conf.example ]; then
    cp /usr/local/etc/unifi-poller/up.conf.example /usr/local/etc/unifi-poller/up.conf
  fi

  # Allow admins to change the configuration and delete the docs.
  chgrp -R admin /usr/local/etc/unifi-poller /usr/local/share/doc/unifi-poller
  chmod -R g+wr /usr/local/etc/unifi-poller  /usr/local/share/doc/unifi-poller

  # Make sure admins can delete logs.
  chown -R nobody:admin /usr/local/var/log/unifi-poller
  chmod 0775 /usr/local/var/log/unifi-poller
  chmod -R g+rw /usr/local/var/log/unifi-poller

  # Restart the service - this starts the application as user nobody.
  launchctl unload /Library/LaunchAgents/com.github.davidnewhall.unifi-poller.plist
  launchctl load /Library/LaunchAgents/com.github.davidnewhall.unifi-poller.plist

elif [ -x "/bin/systemctl" ]; then
  # Reload and restart - this starts the application as user nobody.
  /bin/systemctl daemon-reload
  /bin/systemctl enable unifi-poller
  /bin/systemctl restart unifi-poller
fi
