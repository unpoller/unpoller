#!/bin/sh

# This is a quick and dirty script to install the latest Linux package.
#
# Use it like this, pick curl or wget:  (sudo is optional)
# ----
#   curl -s https://raw.githubusercontent.com/unifi-poller/unifi-poller/master/scripts/install.sh | sudo sh
#   wget -qO- https://raw.githubusercontent.com/unifi-poller/unifi-poller/master/scripts/install.sh | sudo sh
# ----
#
# - If you're on RedHat/CentOS/Fedora, installs the latest rpm package.
# - If you're on Debian/Ubuntu/Knoppix, installs the latest deb package.
# - If you're on FreeBSD, installs the latest txz package.
#
# This is part of application-builder.
# https://github.com/golift/application-builder

REPO=unifi-poller/unifi-poller
BREW=golift/mugs/unifi-poller
LATEST=https://api.github.com/repos/${REPO}/releases/latest
ISSUES=https://github.com/${REPO}/issues/new
ARCH=$(uname -m)
OS=$(uname -s)
P=" ==>"

# Nothing else needs to be changed. Unless you're fixing things!
echo "<-------------------------------------------------->"

if [ "$OS" = "Darwin" ]; then
  echo "${P} On a mac? Use Homebrew:"
  echo "     brew install ${BREW}"
  exit
fi

# $ARCH is passed into egrep to find the right file.
if [ "$ARCH" = "x86_64" ] || [ "$ARCH" = "amd64" ]; then
  ARCH="x86_64|amd64"
elif [[ $ARCH = *386* ]] || [[ $ARCH = *686* ]]; then
  ARCH="i386"
elif [[ $ARCH = *arm64* ]] || [[ $ARCH = *armv8* ]] || [[ $ARCH = *aarch64* ]]; then
  ARCH="arm64"
elif [[ $ARCH = *armv6* ]] || [[ $ARCH = *armv7* ]]; then
  ARCH="armhf"
else
  echo "${P} [ERROR] Unknown Architecture: ${ARCH}"
  echo "${P} $(uname -a)"
  echo "${P} Please report this, along with the above OS details:"
  echo "     ${ISSUES}"
  exit 1
fi

if [ "$1" = "deb" ] || [ "$1" = "rpm" ] || [ "$1" = "txz" ]; then
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
   else
    pkg --version > /dev/null 2>&1
    if [ "$?" = "0" ]; then
      FILE=txz
    fi
   fi
  fi
fi

if [ "$FILE" = "" ]; then
  echo "${P} [ERROR] No pkg (freebsd), dpkg (debian) or rpm (redhat) package managers found; not sure what package to download!"
  echo "${P} $(uname -a)"
  echo "${P} If you feel this is a mistake, please report this along with the above OS details:"
  echo "     ${ISSUES}"
  exit 1
fi

# curl or wget?
curl --version > /dev/null 2>&1
if [ "$?" = "0" ]; then
  CMD="curl -sL"
else
  wget --version > /dev/null 2>&1
  if [ "$?" = "0" ]; then
    CMD="wget -qO-"
  fi
fi

if [ "$CMD" = "" ]; then
  echo "${P} [ERROR] Could not locate curl nor wget - please install one to download packages!"
  exit 1
fi

# Grab latest release file from github.
URL=$($CMD ${LATEST} | egrep "browser_download_url.*(${ARCH})\.${FILE}\"" | cut -d\" -f 4)

if [ "$?" != "0" ] || [ "$URL" = "" ]; then
  echo "${P} [ERROR] Missing latest release for '${FILE}' file ($OS/${ARCH}) at ${LATEST}"
  echo "${P} $(uname -a)"
  echo "${P} Please report error this, along with the above OS details:"
  echo "     ${ISSUES}"
  exit 1
fi

INSTALLER="rpm -Uvh"
if [ "$FILE" = "deb" ]; then
  INSTALLER="dpkg --force-confdef --force-confold --install"
elif [ "$FILE" = "txz" ]; then
  INSTALLER="pkg install"
fi

FILE=$(basename ${URL})
echo "${P} Downloading: ${URL}"
echo "${P} To Location: /tmp/${FILE}"
$CMD ${URL} > /tmp/${FILE}

# Install it.
if [ "$(id -u)" = "0" ]; then
  echo "${P} Downloaded. Installing the package!"
  echo "${P} Executing: ${INSTALLER} /tmp/${FILE}"
  $INSTALLER /tmp/${FILE}
  echo "<-------------------------------------------------->"
else
  echo "${P} Downloaded. Install the package like this:"
  echo "     sudo $INSTALLER /tmp/${FILE}"
fi
