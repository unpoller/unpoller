#!/bin/sh


OS="$(uname -s)"

if [ "${OS}" = "Linux" ]; then
    # Make a user and group for this app, but only if it does not already exist.
    id unpoller >/dev/null 2>&1  || \
        useradd --system --user-group --no-create-home --home-dir /tmp --shell /bin/false unpoller
elif [ "${OS}" = "OpenBSD" ]; then
    id unpoller >/dev/null 2>&1  || \
        useradd  -g =uid -d /tmp -s /bin/false unpoller
elif [ "${OS}" = "FreeBSD" ]; then
    id unpoller >/dev/null 2>&1  || \
        pw useradd unpoller -d /tmp -w no -s /bin/false
else
    echo "Unknown OS: ${OS}, please add system user unpoller manually."
fi
