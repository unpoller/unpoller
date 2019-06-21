_This doc is far from complete._

# Build Pipeline

Lets talk about how the software gets built for our users before we talk about
making changes to it.

## TravisCI

This repo is tested, built and deployed by [Travis-CI](https://travis-ci.org/davidnewhall/unifi-poller).

The [.travis.yml](.travis.yml) file in this repo coordinates the entire process.
As long as this document is kept up to date, this is what the travis file does:

-   Creates a go-capable build environment on a Linux host, some debian variant.
-   Install ruby-devel to get rubygems.
-   Installs other build tools including rpm and fpm from rubygems.
-   Starts docker, builds the docker container and runs it.
-   Tests that the Docker container ran and produced expected output.
-   Makes a release. `make release`: This does a lot of things, controlled by the [Makefile](Makefile).
    -   Runs go tests and go linters.
    -   Compiles the application binaries for Linux and macOS.
    -   Compiles a man page that goes into the packages.
    -   Creates rpm and deb packages using fpm.
    -   Puts the packages, gzipped binaries and a file containing the SHA256s of each asset into a release folder.

After the release is built and Docker image tested:
-   Deploys the release assets to the tagged release on GitHub using an encrypted GitHub Token (api key).
-   Runs [another script](scripts/formula-deploy.sh) to create and upload a Homebrew formula to [golift/homebrew-mugs](https://github.com/golift/homebrew-mugs).
    -   Uses an encrypted SSH key to upload the updated formula to the repo.
-   Travis does nothing else with Docker; it just makes sure the thing compiles and runs.

### Homebrew

it's a mac thing.

[Homebrew](https://brew.sh) is all I use at home. Please don't break the homebrew
formula stuff; it took a lot of pain to get it just right. I am very interested
in how it works for you.

### Docker

Docker is built automatically by Docker Cloud using the Dockerfile in the path
[init/docker/Dockerfile](init/docker/Dockerfile). Other than the Dockerfile, all
the configuration is done in the Cloud service under my personal account `golift`.

If you have need to change the Dockerfile, please clearly explain what problem your
changes are solving, and how it has been tested and validated. As far as I'm
concerned this file should never need to change again, but I'm not a Docker expert;
you're welcome to prove me wrong.

# Contributing

Make a pull request and tell me what you're fixing. Pretty simple. If I need to
I'll add more "rules." For now I'm happy to have help. Thank you!

## Wiki

**If you see typos, errors, omissions, etc, please fix them.**

At this point, the wiki is pretty solid. Please keep your edits brief and without
too much opinion. If you want to provide a way to do something, please also provide
any alternatives you're aware of. If you're not sure, just open an issue and we can
hash it out. I'm reasonable.

## Unifi Library

If you're trying to fix something in the unifi data collection (ie. you got an
unmarshal error, or you want to add something I didn't include) then you
should look at the [unifi library](https://github.com/golift/unifi). All the
data collection and export code lives there. Contributions and Issues are welcome
on that code base as well.
