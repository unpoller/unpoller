#!/bin/bash

# This file is used by rpm and deb packages.  FPM use.

systemctl daemon-reload
systemctl restart unifi-poller
