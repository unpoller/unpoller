#!/bin/bash

# This script creates a local installation of unifi-poller.
# Recommend using Makefile to invoke: make install
# Supports Linux (systemd only) and macOS.

BINARY=unifi-poller

echo "Installing unifi-poller. If you get errors, you may need sudo."

# Install binary.
GOBIN=/usr/local/bin go install -ldflags "-w -s" ./...

# Making config folders and installing man page.
mkdir -p /usr/local/etc/${BINARY} /usr/local/share/man/man1
mv *.1.gz /usr/local/share/man/man1

# Installing config file, man page and launch agent or systemd unit file.
if [ ! -f /usr/local/etc/${BINARY}/up.conf ]; then
  cp examples/up.conf.example /usr/local/etc/${BINARY}/up.conf
fi
if [ -d ~/Library/LaunchAgents ]; then
  cp init/launchd/com.github.davidnewhall.${BINARY}.plist ~/Library/LaunchAgents
fi
if [ -d /etc/systemd/system ]; then
  cp init/systemd/${BINARY}.service /etc/systemd/system
fi

# Making systemd happy by telling it to reload.
if [ -x /bin/systemctl ]; then
  /bin/systemctl --system daemon-reload
fi

echo "Installation Complete. Edit the config file @ /usr/local/etc/${BINARY}/up.conf"
echo "Then start the daemon with:"
if [ -d ~/Library/LaunchAgents ]; then
  echo "   launchctl load ~/Library/LaunchAgents/com.github.davidnewhall.${BINARY}.plist"
fi
if [ -d /etc/systemd/system ]; then
  echo "   sudo /bin/systemctl start ${BINARY}"
fi
echo "Examine the log file at: /usr/local/var/log/${BINARY}.log (logs may go elsewhere on linux, check syslog)"
