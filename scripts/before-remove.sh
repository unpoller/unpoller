#!/bin/bash

# This file is used by rpm and deb packages. FPM use.
# Edit this file as needed for your application.
# This file is only installed if FORMULA is set to service.

if [ "$1" = "upgrade" ] || [ "$1" = "1" ] ; then
  exit 0
fi

if [ -x "/bin/systemctl" ]; then
  /bin/systemctl stop {{BINARY}}
  /bin/systemctl disable {{BINARY}}
fi
