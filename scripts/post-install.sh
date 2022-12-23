#!/bin/sh

# Step 1, decide if we should use systemd or init/upstart
use_systemctl="True"
systemd_version=0
if ! command -V systemctl >/dev/null 2>&1; then
  use_systemctl="False"
else
    systemd_version=$(systemctl --version | head -1 | sed 's/systemd //g' | sed 's/ .*//')
fi

cleanup() {
    # This is where you remove files that were not needed on this platform / system
    if [ "${use_systemctl}" = "False" ]; then
        rm -f /etc/systemd/system/unpoller.service
    else
        rm -f /etc/chkconfig/unpoller
        rm -f /etc/init.d/unpoller
    fi
}

OS="$(uname -s)"

cleanInstall() {
    printf "\033[32m Post Install of an clean install\033[0m\n"
    # Step 3 (clean install), enable the service in the proper way for this platform

    if [ "${use_systemctl}" = "False" ]; then
        if command -V chkconfig >/dev/null 2>&1; then
          chkconfig --add unpoller
        fi

        service unpoller restart ||:
    else
        # rhel/centos7 cannot use ExecStartPre=+ to specify the pre start should be run as root
        # even if you want your service to run as non root.
        if [ "${systemd_version}" -lt 231 ]; then
            printf "\033[31m systemd version %s is less then 231, fixing the service file \033[0m\n" "${systemd_version}"
            sed -i "s/=+/=/g" /etc/systemd/system/unpoller.service
        fi
        printf "\033[32m Reload the service unit from disk\033[0m\n"
        systemctl daemon-reload ||:
        printf "\033[32m Unmask the service\033[0m\n"
        systemctl unmask unpoller ||:
        printf "\033[32m Set the preset flag for the service unit\033[0m\n"
        systemctl preset unpoller ||:
        printf "\033[32m Set the enabled flag for the service unit\033[0m\n"
        systemctl enable unpoller ||:
        systemctl restart unpoller ||:
    fi
}

upgrade() {
    printf "\033[32m Post Install of an upgrade\033[0m\n"
    # Step 3(upgrade), do what you need
    
}

# Step 2, check if this is a clean install or an upgrade
action="$1"
if  [ "$1" = "configure" ] && [ -z "$2" ]; then
  # Alpine linux does not pass args, and deb passes $1=configure
  action="install"
elif [ "$1" = "configure" ] && [ -n "$2" ]; then
    # deb passes $1=configure $2=<current version>
    action="upgrade"
fi

case "$action" in
  "1" | "install")
    cleanInstall
    ;;
  "2" | "upgrade")
    printf "\033[32m Post Install of an upgrade\033[0m\n"
    upgrade
    ;;
  *)
    # $1 == version being installed
    printf "\033[32m Alpine\033[0m"
    cleanInstall
    ;;
esac

# Step 4, clean up unused files, yes you get a warning when you remove the package, but that is ok.
cleanup
