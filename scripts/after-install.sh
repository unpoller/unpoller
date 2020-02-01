#!/bin/sh

# This file is used by txz, deb and rpm packages.
# FPM adds this as the after-install script to all packages.

if [ -x "/bin/systemctl" ]; then
  # Reload and restart - this starts the application as user nobody.
  /bin/systemctl daemon-reload
  /bin/systemctl enable unifi-poller
  /bin/systemctl restart unifi-poller
elif [ -x /usr/sbin/service ]; then
  # Do not start or restart on freebsd. That's "bad practice."
  /usr/sbin/service unifi-poller enabled || /usr/sbin/service unifi-poller enable
fi
