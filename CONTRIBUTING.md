_This doc is far from complete._

# Build Pipeline

Lets talk about how the software gets built for our users before we talk about
making changes to it.


## Github Actions

This repo is tested, built and deployed by [Github Actions](https://github.com/unpoller/unpoller/actions).

The [.github/](.github/) directory in this repo coordinates the entire process.
As long as this document is kept up to date, this is what github does:

-   Builds and Tests code changes
-   Lints code changes
-   On Release (through git tags) it uses goreleaser-pro to build and release:
-      Linux, Mac and Windows Binaries
-      Provides a packaged source copy
-      Builds Debian, RedHat packages
-      Builds Mac universal binary
-      Builds Windows executable
-      Builds numerous platform docker images and uploads them

After the release is built and Docker image tested:
-   Deploys the release assets to the tagged release on [GitHub releases](https://github.com/unpoller/unpoller/releases)

### Homebrew

it's a mac thing. [Homebrew](https://brew.sh)

### Docker

Docker is built automatically and uploaded to ghcr.io by the release github action.

# Contributing

Make a pull request and tell me what you're fixing. Pretty simple. If I need to
I'll add more "rules." For now I'm happy to have help. Thank you!

## Wiki

**If you see typos, errors, omissions, etc, please fix them.**

At this point, the wiki is pretty solid. Please keep your edits brief and without
too much opinion. If you want to provide a way to do something, please also provide
any alternatives you're aware of. If you're not sure, just open an issue and we can
hash it out. I'm reasonable.

## UniFi Libraries

The UniFi data extraction is provided as an [external library](https://godoc.org/github.com/unifi-poller/unifi),
and you can import that code directly without futzing with this application. That
means, if you wanted to do something like make telegraf collect your data instead
of UniFi Poller you can achieve that with a little bit of Go code. You could write
a small app that acts as a telegraf input plugin using the [unifi](https://github.com/unifi-poller/unifi)
library to grab the data from your controller.

This application is very dynamic and built using several package repos.
They are all in the [UniFi Poller GitHub Org](https://github.com/unifi-poller).
