#!/bin/sh
#
# FreeBSD rc.d startup script for {{BINARY}}.
#
# PROVIDE: {{BINARY}}
# REQUIRE: networking syslog
# KEYWORD:

. /etc/rc.subr

name="{{BINARYU}}"
real_name="{{BINARY}}"
rcvar="{{BINARYU}}_enable"
{{BINARYU}}_command="/usr/local/bin/${real_name}"
{{BINARYU}}_user="{{BINARY}}"
{{BINARYU}}_config="/usr/local/etc/${real_name}/{{CONFIG_FILE}}"
pidfile="/var/run/${real_name}/pid"

# This runs `daemon` as the `{{BINARYU}}_user` user.
command="/usr/sbin/daemon"
command_args="-P ${pidfile} -r -t ${real_name} -T ${real_name} -l daemon ${{{BINARYU}}_command} -c ${{{BINARYU}}_config}"

load_rc_config ${name}
: ${{{BINARYU}}_enable:=no}

# Make a place for the pid file.
mkdir -p $(dirname ${pidfile})
chown -R ${{BINARYU}}_user $(dirname ${pidfile})

# Suck in optional exported override variables.
# ie. add something like the following to this file: export UP_POLLER_DEBUG=true
[ -f "/usr/local/etc/defaults/${real_name}" ] && . "/usr/local/etc/defaults/${real_name}"

# Go!
run_rc_command "$1"
