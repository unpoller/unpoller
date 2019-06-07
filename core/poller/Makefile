BINARY:=unifi-poller
URL=https://github.com/davidnewhall/unifi-poller
MAINT="david at sleepers dot pro"
DESC="This daemon polls a Unifi controller at a short interval and stores the collected metric data in an Influx Database."
PACKAGE:=./cmd/$(BINARY)
VERSION:=$(shell git tag -l --merged | tail -n1 | tr -d v)
ITERATION:=$(shell git rev-list --all --count)

all: man build

# Prepare a release. Called in Travis CI.
release: clean test $(BINARY)-$(VERSION)-$(ITERATION).x86_64.rpm $(BINARY)_$(VERSION)-$(ITERATION)_amd64.deb $(BINARY)-$(VERSION).pkg
	# Prepareing a release!
	mkdir -p release
	gzip -9 $(BINARY).linux
	gzip -9 $(BINARY).macos
	mv $(BINARY)-$(VERSION)-$(ITERATION).x86_64.rpm $(BINARY)_$(VERSION)-$(ITERATION)_amd64.deb \
		$(BINARY)-$(VERSION).pkg $(BINARY).macos.gz $(BINARY).linux.gz release/

# Delete all build assets.
clean:
	# Cleaning up.
	rm -f $(BINARY){.macos,.linux,.1,}{,.gz}
	rm -f $(BINARY){_,-}*.{deb,rpm,pkg}
	rm -rf package_build_* release

# Build a man page from a markdown file using ronn.
man: $(BINARY).1.gz
$(BINARY).1.gz:
	# Building man page.
	@ronn --version > /dev/null || (echo "Ronn missing. Install ronn: $(URL)/wiki/Ronn" && false)
	ronn < "$(PACKAGE)/README.md" | gzip -9 > "$(BINARY).1.gz"

# Binaries

build: $(BINARY)
$(BINARY):
	go build -o $(BINARY) -ldflags "-w -s -X main.Version=$(VERSION)" $(PACKAGE)

linux: $(BINARY).linux
$(BINARY).linux:
	# Building linux binary.
	GOOS=linux go build -o $(BINARY).linux -ldflags "-w -s -X main.Version=$(VERSION)" $(PACKAGE)

macos: $(BINARY).macos
$(BINARY).macos:
	# Building darwin binary.
	GOOS=darwin go build -o $(BINARY).macos -ldflags "-w -s -X main.Version=$(VERSION)" $(PACKAGE)

# Packages

rpm: clean $(BINARY)-$(VERSION)-$(ITERATION).x86_64.rpm
$(BINARY)-$(VERSION)-$(ITERATION).x86_64.rpm: check_fpm package_build_linux
	@echo "Building 'rpm' package for $(BINARY) version '$(VERSION)-$(ITERATION)'."
	fpm -s dir -t rpm \
		--name $(BINARY) \
		--version $(VERSION) \
		--iteration $(ITERATION) \
		--after-install scripts/after-install.sh \
		--before-remove scripts/before-remove.sh \
		--license MIT \
		--url $(URL) \
		--maintainer $(MAINT) \
		--description $(DESC) \
		--chdir package_build_linux

deb: clean $(BINARY)_$(VERSION)-$(ITERATION)_amd64.deb
$(BINARY)_$(VERSION)-$(ITERATION)_amd64.deb: check_fpm package_build_linux
	@echo "Building 'deb' package for $(BINARY) version '$(VERSION)-$(ITERATION)'."
	fpm -s dir -t deb \
		--name $(BINARY) \
		--version $(VERSION) \
		--iteration $(ITERATION) \
		--after-install scripts/after-install.sh \
		--before-remove scripts/before-remove.sh \
		--license MIT \
		--url $(URL) \
		--maintainer $(MAINT) \
		--description $(DESC) \
		--chdir package_build_linux

osxpkg: clean $(BINARY)-$(VERSION).pkg
$(BINARY)-$(VERSION).pkg: check_fpm package_build_osx
	@echo "Building 'osx' package for $(BINARY) version '$(VERSION)-$(ITERATION)'."
	fpm -s dir -t osxpkg \
		--name $(BINARY) \
		--version $(VERSION) \
		--iteration $(ITERATION) \
		--after-install scripts/after-install.sh \
		--osxpkg-identifier-prefix com.github.davidnewhall \
		--license MIT \
		--url $(URL) \
		--maintainer $(MAINT) \
		--description $(DESC) \
		--chdir package_build_osx

# OSX packages use /usr/local because Apple doesn't allow writing many other places.
package_build_osx: man macos
	# Building package environment for macOS.
	mkdir -p $@/usr/local/bin $@/usr/local/etc/$(BINARY) $@/Library/LaunchAgents
	mkdir -p $@/usr/local/share/man/man1 $@/usr/local/share/doc/$(BINARY) $@/usr/local/var/log
	# Copying the binary, config file and man page into the env.
	cp $(BINARY).macos $@/usr/local/bin/$(BINARY)
	cp *.1.gz $@/usr/local/share/man/man1
	cp examples/*.conf.example $@/usr/local/etc/$(BINARY)/
	cp examples/* $@/usr/local/share/doc/$(BINARY)/
	cp init/launchd/com.github.davidnewhall.unifi-poller.plist $@/Library/LaunchAgents/

# Build an environment that can be packaged for linux.
package_build_linux: man linux
	# Building package environment for linux.
	mkdir -p $@/usr/bin $@/etc/$(BINARY) $@/lib/systemd/system
	mkdir -p $@/usr/share/man/man1 $@/usr/share/doc/$(BINARY)
	# Copying the binary, config file and man page into the env.
	cp $(BINARY).linux $@/usr/bin/$(BINARY)
	cp *.1.gz $@/usr/share/man/man1
	cp examples/*.conf.example $@/etc/$(BINARY)/
	cp examples/up.conf.example $@/etc/$(BINARY)/up.conf
	cp examples/* $@/usr/share/doc/$(BINARY)/
	# Fixing the paths in the systemd unit file before copying it into the emv.
	sed "s%ExecStart.*%ExecStart=/usr/bin/$(BINARY) --config=/etc/$(BINARY)/up.conf%" \
		init/systemd/unifi-poller.service > $@/lib/systemd/system/$(BINARY).service

check_fpm:
	@fpm --version > /dev/null || (echo "FPM missing. Install FPM: https://fpm.readthedocs.io/en/latest/installing.html" && false)

# Extras

# Run code tests and lint.
test: lint
	# Testing.
	go test -race -covermode=atomic $(PACKAGE)
lint:
	# Checking lint.
	golangci-lint run --enable-all -D gochecknoglobals

# Install locally into /usr/local. Not recommended.
install: man
	scripts/local_install.sh

# If you installed with `make install` run `make uninstall` before installing a binary package.
# This will remove the package install from macOS, it will not remove a package install from Linux.
uninstall:
	[ -x /bin/systemctl ] && /bin/systemctl disable $(BINARY) || true
	[ -x /bin/systemctl ] && /bin/systemctl stop $(BINARY) || true
	[ -x /bin/launchctl ] && [ -f ~/Library/LaunchAgents/com.github.davidnewhall.$(BINARY).plist ] \
		&& /bin/launchctl unload ~/Library/LaunchAgents/com.github.davidnewhall.$(BINARY).plist || true
	[ -x /bin/launchctl ] && [ -f /Library/LaunchAgents/com.github.davidnewhall.$(BINARY).plist ] \
		&& /bin/launchctl unload /Library/LaunchAgents/com.github.davidnewhall.$(BINARY).plist || true
	rm -rf /usr/local/{etc,bin}/$(BINARY) /usr/local/share/man/man1/$(BINARY).1.gz
	rm -f ~/Library/LaunchAgents/com.github.davidnewhall.$(BINARY).plist
	rm -f /etc/systemd/system/$(BINARY).service
	[ -x /bin/systemctl ] && /bin/systemctl --system daemon-reload || true
	@[ -x /bin/launchctl ] && [ -f /Library/LaunchAgents/com.github.davidnewhall.$(BINARY).plist ] \
		&& echo "  ==> Delete this file manually: sudo rm -f /Library/LaunchAgents/com.github.davidnewhall.$(BINARY).plist" || true

# Don't run this unless you're ready to debug untested vendored dependencies.
deps:
	dep ensure -update
