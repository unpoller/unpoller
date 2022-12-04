#!/bin/sh

# Step 1, decide if we should use systemd or init/upstart
use_systemctl="True"
systemd_version=0
if ! command -V systemctl >/dev/null 2>&1; then
  use_systemctl="False"
else
    systemd_version=$(systemctl --version | head -1 | sed 's/systemd //g')
fi

remove() {
    printf "\033[32m Post Remove of a normal remove\033[0m\n"
    echo "Remove" > /tmp/postremove-proof
}

purge() {
    printf "\033[32m Post Remove purge, deb only\033[0m\n"
    echo "Purge" > /tmp/postremove-proof
}

upgrade() {
    printf "\033[32m Post Remove of an upgrade\033[0m\n"
    echo "Upgrade" > /tmp/postremove-proof
}

echo "$@"

action="$1"

case "$action" in
  "0" | "remove")
    remove
    ;;
  "1" | "upgrade")
    upgrade
    ;;
  "purge")
    purge
    ;;
  *)
    printf "\033[32m Alpine\033[0m"
    remove
    ;;
esac
