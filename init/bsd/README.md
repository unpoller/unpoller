Generic FreeBSD rc.d service file lives here.

-   Marshal template like so (example in [Makefile](../../Makefile)):
```shell
    sed -e "s/{{BINARY}}/app-name/g" \
        -e "s/{{BINARYU}}/app_name/g" \
        -e "s/{{CONFIG_FILE}}/app-name.conf/g" \
        freebsd.rc.d > /usr/local/etc/rc.d/app-name
```
