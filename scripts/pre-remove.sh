#!/bin/bash

# Step 1, decide if we should use systemd or init/upstart
use_systemctl="True"
systemd_version=0
if ! command -V systemctl >/dev/null 2>&1; then
  use_systemctl="False"
else
    systemd_version=$(systemctl --version | head -1 | sed 's/systemd //g')
fi

if [ "$1" = "upgrade" ] || [ "$1" = "1" ] ; then
  exit 0
fi

if [ "${use_systemctl}" = "False" ]; then
  service unpoller stop
  service unpoller disable
else
  systemctl stop unpoller
  systemctl disable unpoller
fi
