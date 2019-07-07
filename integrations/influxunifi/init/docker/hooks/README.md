# Docker Build Hooks

The files in this folder are used by Docker Cloud to automate image builds.

If you want to build, maintain and push multi-architecture Docker images, you may
follow the example provided here. All of the hooks are generic, and will work with
any build. Two environment variables must be passed in from Docker Cloud config.
`BUILDS` must be set to the builds you're trying to perform. This repo is currently
set to `linux:armhf:arm: linux:arm64:arm64:armv8 linux:amd64:amd64: linux:i386:386:`.
The format is `os:name:arch:variant`. `os` and `name` are passed into the Dockerfile.
`os`, `arch` and `variant` are passed into `docker manifest annotate`. This does not
yet work with an OS other than `linux`.

Keep the build simple; see screenshot. This only supports one build tag, but it creates many more.

![UniFi Poller Docker Cloud Build Rules](https://raw.githubusercontent.com/wiki/davidnewhall/unifi-poller/images/unifi-poller-build-rules.png "UniFi Poller Docker Cloud Build Rules")

The fancy source tag is `/^v((\d+\.\d+)(?:\.\d+)?)$/` and it allows you to capture
the minor version without patch-level in `{\2}`. I no longer use `{\2}` in my build.
[See how it works here](https://regex101.com/r/fzt6ki/1).
