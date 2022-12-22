#!/bin/sh
#
# FreeBSD rc.d startup script for unpoller.
#
# PROVIDE: unpoller
# REQUIRE: networking syslog
# KEYWORD:

. /etc/rc.subr

name="unpoller"
real_name="unpoller"
rcvar="unpoller_enable"
unpoller_command="/usr/local/bin/${real_name}"
unpoller_user="unpoller"
unpoller_config="/usr/local/etc/${real_name}/up.conf"
pidfile="/var/run/${real_name}/pid"

# This runs `daemon` as the `unpoller_user` user.
command="/usr/sbin/daemon"
command_args="-P ${pidfile} -r -t ${real_name} -T ${real_name} -l daemon ${unpoller_command} -c ${unpoller_config}"

load_rc_config ${name}
: ${unpoller_enable:=no}

# Make a place for the pid file.
mkdir -p $(dirname ${pidfile})
chown -R $unpoller_user $(dirname ${pidfile})

# ensure log directory exists
mkdir -p /usr/local/var/log/${real_name}
chown -R $unpoller_user /usr/local/var/log/${real_name}

# Suck in optional exported override variables.
# ie. add something like the following to this file: export UP_POLLER_DEBUG=true
[ -f "/usr/local/etc/defaults/${real_name}" ] && . "/usr/local/etc/defaults/${real_name}"

# Go!
run_rc_command "$1"
