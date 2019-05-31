#!/bin/bash

# This script removes a local installation of unifi-poller.
# Recommend using Makefile to invoke: make uninstall
# Supports Linux (systemd only) and macOS.

BINARY=unifi-poller

echo "Uninstall unifi-poller. You may need sudo on Linux. Do not use sudo on macOS."

# Stopping the daemon
if [ -x /bin/systemctl ]; then
  /bin/systemctl disable ${BINARY}
  /bin/systemctl stop ${BINARY}
fi

if [ -x /bin/launchctl ] && [ -f ~/Library/LaunchAgents/com.github.davidnewhall.${BINARY}.plist ]; then
  echo Unloading ~/Library/LaunchAgents/com.github.davidnewhall.${BINARY}.plist
  /bin/launchctl unload ~/Library/LaunchAgents/com.github.davidnewhall.${BINARY}.plist || true
fi

if [ -x /bin/launchctl ] && [ -f /Library/LaunchAgents/com.github.davidnewhall.${BINARY}.plist ]; then
  echo Unloading /Library/LaunchAgents/com.github.davidnewhall.${BINARY}.plist
  /bin/launchctl unload /Library/LaunchAgents/com.github.davidnewhall.${BINARY}.plist || true
  echo "Delete this file manually: sudo rm -f /Library/LaunchAgents/com.github.davidnewhall.${BINARY}.plist"
fi

# Deleting config file, binary, man page, launch agent or unit file.
rm -rf /usr/local/{etc,bin}/${BINARY} /usr/local/share/man/man1/${BINARY}.1.gz
rm -f ~/Library/LaunchAgents/com.github.davidnewhall.${BINARY}.plist
rm -f /etc/systemd/system/${BINARY}.service

# Making systemd happy by telling it to reload.
if [ -x /bin/systemctl ]; then
  /bin/systemctl --system daemon-reload
fi
