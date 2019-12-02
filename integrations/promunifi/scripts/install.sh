#!/bin/bash

# This is a quick and drity script to install the latest Linux package.
#
# Use it like this:  (sudo is optional)
# ===
#   curl https://raw.githubusercontent.com/davidnewhall/unifi-poller/master/scripts/install.sh | sudo bash
# ===
# If you're on redhat, this installs the latest rpm. If you're on Debian, it installs the latest deb package.
#
# This is part of application-builder.
# https://github.com/golift/application-builder

REPO=davidnewhall/unifi-poller
LATEST=https://api.github.com/repos/${REPO}/releases/latest
ARCH=$(uname -m)

# $ARCH is passed into egrep to find the right file.
if [ "$ARCH" = "x86_64" ] || [ "$ARCH" = "amd64" ]; then
  ARCH="x86_64|amd64"
elif [[ $ARCH == *386* ]] || [[ $ARCH == *686* ]]; then
  ARCH="i386"
elif [[ $ARCH == *arm64* ]] || [[ $ARCH == *armv8* ]]; then
  ARCH="arm64"
elif [[ $ARCH == *armv6* ]] || [[ $ARCH == *armv7* ]]; then
  ARCH="armhf"
else
  echo "Unknown Architecture. Submit a pull request to fix this, please."
  echo ==> $ARCH
  exit 1
fi

if [ "$1" == "deb" ] || [ "$1" == "rpm" ]; then
  FILE=$1
else
  # If you have both, rpm wins.
  rpm --version > /dev/null 2>&1
  if [ "$?" = "0" ]; then
    FILE=rpm
  else
   dpkg --version > /dev/null 2>&1
   if [ "$?" = "0" ]; then
     FILE=deb
   fi
  fi
fi

if [ "$FILE" = "" ]; then
  echo "No dpkg or rpm package managers found!"
  exit 1
fi

# curl or wget?
curl --version > /dev/null 2>&1
if [ "$?" = "0" ]; then
  CMD="curl -L"
else
  wget --version > /dev/null 2>&1
  if [ "$?" = "0" ]; then
    CMD="wget -O-"
  fi
fi

if [ "$CMD" = "" ]; then
  echo "Need curl or wget - could not find either!"
  exit 1
fi

# Grab latest release file from github.
URL=$($CMD ${LATEST} | egrep "browser_download_url.*(${ARCH})\.${FILE}\"" | cut -d\" -f 4)

if [ "$?" != "0" ] || [ "$URL" = "" ]; then
  echo "Error locating latest release at ${LATEST}"
  exit 1
fi

INSTALLER="rpm -Uvh"
if [ "$FILE" = "deb" ]; then
  INSTALLER="dpkg --force-confdef --force-confold --install"
fi

FILE=$(basename ${URL})
echo "Downloading: ${URL} to /tmp/${FILE}"
$CMD ${URL} > /tmp/${FILE}

# Install it.
if [ "$(id -u)" = "0" ]; then
  echo "==================================="
  echo "Downloaded. Installing the package!"
  echo "Running: ${INSTALLER} /tmp/${FILE}"
  $INSTALLER /tmp/${FILE}
else
  echo "================================"
  echo "Downloaded. Install the package:"
  echo "sudo $INSTALLER /tmp/${FILE}"
fi
