# This Makefile is written as generic as possible.
# Setting these variables and creating the necesarry paths in your GitHub repo will make this file work.
#
# github username
GHUSER=davidnewhall
# docker hub username
DHUSER=golift
MAINT=David Newhall II <david at sleepers dot pro>
DESC=Polls a UniFi controller and stores metrics in InfluxDB
GOLANGCI_LINT_ARGS=--enable-all -D gochecknoglobals
BINARY:=$(shell basename $(shell pwd))
URL:=https://github.com/$(GHUSER)/$(BINARY)
CONFIG_FILE=up.conf

# These don't generally need to be changed.

# md2roff turns markdown into man files and html files.
MD2ROFF_BIN=github.com/github/hub/md2roff-bin
# This produces a 0 in some envirnoments (like Docker), but it's only used for packages.
ITERATION:=$(shell git rev-list --count HEAD || echo 0)
# Travis CI passes the version in. Local builds get it from the current git tag.
ifeq ($(VERSION),)
	VERSION:=$(shell git tag -l --merged | tail -n1 | tr -d v || echo development)
endif
# rpm is wierd and changes - to _ in versions.
RPMVERSION:=$(shell echo $(VERSION) | tr -- - _)
DATE:=$(shell date)

# This parameter is passed in as -X to go build. Used to override the Version variable in a package.
# This makes a path like github.com/davidnewhall/unifi-poller/unifipoller.Version=1.3.3
# Name the Version-containing library the same as the github repo, without dashes.
VERSION_PATH:=github.com/$(GHUSER)/$(BINARY)/$(shell echo $(BINARY) | tr -d -- -).Version=$(VERSION)

# Makefile targets follow.

all: man build

# Prepare a release. Called in Travis CI.
release: clean vendor test macos windows $(BINARY)-$(RPMVERSION)-$(ITERATION).x86_64.rpm $(BINARY)_$(VERSION)-$(ITERATION)_amd64.deb
	# Prepareing a release!
	mkdir -p release
	mv $(BINARY).linux $(BINARY).macos release/
	gzip -9r release/
	zip -9qm release/$(BINARY).exe.zip $(BINARY).exe
	mv $(BINARY)-$(RPMVERSION)-$(ITERATION).x86_64.rpm $(BINARY)_$(VERSION)-$(ITERATION)_amd64.deb release/
	# Generating File Hashes
	for i in release/*; do /bin/echo -n "$$i " ; (openssl dgst -r -sha256 "$$i" | head -c64 ; echo) | tee "$$i.sha256.txt"; done

# Delete all build assets.
clean:
	# Cleaning up.
	rm -f $(BINARY){.macos,.linux,.1,}{,.gz} $(BINARY).rb
	rm -f $(BINARY){_,-}*.{deb,rpm} v*.tar.gz.sha256
	rm -f cmd/$(BINARY)/README{,.html} README{,.html} ./$(BINARY)_manual.html
	rm -rf package_build_* release

# Build a man page from a markdown file using md2roff.
# This also turns the repo readme into an html file.
# md2roff is needed to build the man file and html pages from the READMEs.
man: $(BINARY).1.gz
$(BINARY).1.gz: md2roff
	# Building man page. Build dependency first: md2roff
	go run $(MD2ROFF_BIN) --manual $(BINARY) --version $(VERSION) --date "$(DATE)" examples/MANUAL.md
	gzip -9nc examples/MANUAL > $(BINARY).1.gz
	mv examples/MANUAL.html $(BINARY)_manual.html

md2roff:
	go get $(MD2ROFF_BIN)

# TODO: provide a template that adds the date to the built html file.
readme: README.html
README.html: md2roff
	# This turns README.md into README.html
	go run $(MD2ROFF_BIN) --manual $(BINARY) --version $(VERSION) --date "$(DATE)" README.md

# Binaries

build: $(BINARY)
$(BINARY):
	go build -o $(BINARY) -ldflags "-w -s -X $(VERSION_PATH)"

linux: $(BINARY).linux
$(BINARY).linux:
	# Building linux binary.
	GOOS=linux go build -o $(BINARY).linux -ldflags "-w -s -X $(VERSION_PATH)"

macos: $(BINARY).macos
$(BINARY).macos:
	# Building darwin binary.
	GOOS=darwin go build -o $(BINARY).macos -ldflags "-w -s -X $(VERSION_PATH)"

exe: $(BINARY).exe
windows: $(BINARY).exe
$(BINARY).exe:
	# Building windows binary.
	GOOS=windows go build -o $(BINARY).exe -ldflags "-w -s -X $(VERSION_PATH)"

# Packages

rpm: $(BINARY)-$(RPMVERSION)-$(ITERATION).x86_64.rpm
$(BINARY)-$(RPMVERSION)-$(ITERATION).x86_64.rpm: check_fpm package_build_linux
	@echo "Building 'rpm' package for $(BINARY) version '$(RPMVERSION)-$(ITERATION)'."
	fpm -s dir -t rpm \
		--name $(BINARY) \
		--rpm-os linux \
		--version $(RPMVERSION) \
		--iteration $(ITERATION) \
		--after-install scripts/after-install.sh \
		--before-remove scripts/before-remove.sh \
		--license MIT \
		--url $(URL) \
		--maintainer "$(MAINT)" \
		--description "$(DESC)" \
		--config-files "/etc/$(BINARY)/$(CONFIG_FILE)" \
		--chdir package_build_linux

deb: $(BINARY)_$(VERSION)-$(ITERATION)_amd64.deb
$(BINARY)_$(VERSION)-$(ITERATION)_amd64.deb: check_fpm package_build_linux
	@echo "Building 'deb' package for $(BINARY) version '$(VERSION)-$(ITERATION)'."
	fpm -s dir -t deb \
		--name $(BINARY) \
		--deb-no-default-config-files \
		--version $(VERSION) \
		--iteration $(ITERATION) \
		--after-install scripts/after-install.sh \
		--before-remove scripts/before-remove.sh \
		--license MIT \
		--url $(URL) \
		--maintainer "$(MAINT)" \
		--description "$(DESC)" \
		--config-files "/etc/$(BINARY)/$(CONFIG_FILE)" \
		--chdir package_build_linux

docker:
	docker build -f init/docker/Dockerfile -t $(DHUSER)/$(BINARY) .

# Build an environment that can be packaged for linux.
package_build_linux: readme man linux
	# Building package environment for linux.
	mkdir -p $@/usr/bin $@/etc/$(BINARY) $@/lib/systemd/system
	mkdir -p $@/usr/share/man/man1 $@/usr/share/doc/$(BINARY)
	# Copying the binary, config file, unit file, and man page into the env.
	cp $(BINARY).linux $@/usr/bin/$(BINARY)
	cp *.1.gz $@/usr/share/man/man1
	cp examples/$(CONFIG_FILE).example $@/etc/$(BINARY)/
	cp examples/$(CONFIG_FILE).example $@/etc/$(BINARY)/$(CONFIG_FILE)
	cp LICENSE *.html examples/*.* $@/usr/share/doc/$(BINARY)/
	# These go to their own folder so the img src in the html pages continue to work.
	cp init/systemd/$(BINARY).service $@/lib/systemd/system/

check_fpm:
	@fpm --version > /dev/null || (echo "FPM missing. Install FPM: https://fpm.readthedocs.io/en/latest/installing.html" && false)

# This builds a Homebrew formula file that can be used to install this app from source.
# The source used comes from the released version on GitHub. This will not work with local source.
# This target is used by Travis CI to update the released Forumla when a new tag is created.
formula: $(BINARY).rb
v$(VERSION).tar.gz.sha256:
	# Calculate the SHA from the Github source file.
	curl -sL $(URL)/archive/v$(VERSION).tar.gz | openssl dgst -r -sha256 | tee v$(VERSION).tar.gz.sha256
$(BINARY).rb: v$(VERSION).tar.gz.sha256
	# Creating formula from template using sed.
	sed "s/{{Version}}/$(VERSION)/g;s/{{SHA256}}/`head -c64 v$(VERSION).tar.gz.sha256`/g;s/{{Desc}}/$(DESC)/g;s%{{URL}}%$(URL)%g" init/homebrew/$(BINARY).rb.tmpl | tee $(BINARY).rb

# Extras

# Run code tests and lint.
test: lint
	# Testing.
	go test -race -covermode=atomic ./...
lint:
	# Checking lint.
	golangci-lint run $(GOLANGCI_LINT_ARGS)

# This is safe; recommended even.
dep: vendor
vendor:
	dep ensure

# Don't run this unless you're ready to debug untested vendored dependencies.
deps:
	dep ensure -update

# Homebrew stuff. macOS only.

# Used for Homebrew only. Other distros can create packages.
install: man readme $(BINARY)
	@echo -  Done Building!  -
	@echo -  Local installation with the Makefile is only supported on macOS.
	@echo If you wish to install the application manually on Linux, check out the wiki: $(URL)/wiki/Installation
	@echo -  Otherwise, build and install a package: make rpm -or- make deb
	@echo See the Package Install wiki for more info: $(URL)/wiki/Package-Install
	@[ "$$(uname)" = "Darwin" ] || (echo "Unable to continue, not a Mac." && false)
	@[ "$(PREFIX)" != "" ] || (echo "Unable to continue, PREFIX not set. Use: make install PREFIX=/usr/local ETC=/usr/local/etc" && false)
	@[ "$(ETC)" != "" ] || (echo "Unable to continue, ETC not set. Use: make install PREFIX=/usr/local ETC=/usr/local/etc" && false)
	# Copying the binary, config file, unit file, and man page into the env.
	/usr/bin/install -m 0755 -d $(PREFIX)/bin $(PREFIX)/share/man/man1 $(ETC)/$(BINARY) $(PREFIX)/share/doc/$(BINARY)
	/usr/bin/install -m 0755 -cp $(BINARY) $(PREFIX)/bin/$(BINARY)
	/usr/bin/install -m 0644 -cp $(BINARY).1.gz $(PREFIX)/share/man/man1
	/usr/bin/install -m 0644 -cp examples/$(CONFIG_FILE).example $(ETC)/$(BINARY)/
	[ -f $(ETC)/$(BINARY)/$(CONFIG_FILE) ] || /usr/bin/install -m 0644 -cp  examples/$(CONFIG_FILE).example $(ETC)/$(BINARY)/$(CONFIG_FILE)
	/usr/bin/install -m 0644 -cp LICENSE *.html examples/* $(PREFIX)/share/doc/$(BINARY)/
	# These go to their own folder so the img src in the html pages continue to work.

# If you installed with `make install` run `make uninstall` before installing a binary package. (even on Linux!!!)
# This will remove the package install from macOS, it will not remove a package install from Linux.
uninstall:
	@echo "  ==> You must run make uninstall as root on Linux. Recommend not running as root on macOS."
	[ -x /bin/systemctl ] && /bin/systemctl disable $(BINARY) || true
	[ -x /bin/systemctl ] && /bin/systemctl stop $(BINARY) || true
	[ -x /bin/launchctl ] && [ -f ~/Library/LaunchAgents/com.github.$(GHUSER).$(BINARY).plist ] \
		&& /bin/launchctl unload ~/Library/LaunchAgents/com.github.$(GHUSER).$(BINARY).plist || true
	[ -x /bin/launchctl ] && [ -f /Library/LaunchAgents/com.github.$(GHUSER).$(BINARY).plist ] \
		&& /bin/launchctl unload /Library/LaunchAgents/com.github.$(GHUSER).$(BINARY).plist || true
	rm -rf /usr/local/{etc,bin,share/doc}/$(BINARY)
	rm -f ~/Library/LaunchAgents/com.github.$(GHUSER).$(BINARY).plist
	rm -f /Library/LaunchAgents/com.github.$(GHUSER).$(BINARY).plist || true
	rm -f /etc/systemd/system/$(BINARY).service /usr/local/share/man/man1/$(BINARY).1.gz
	[ -x /bin/systemctl ] && /bin/systemctl --system daemon-reload || true
	@[ -f /Library/LaunchAgents/com.github.$(GHUSER).$(BINARY).plist ] && echo "  ==> Unload and delete this file manually:" && echo "  sudo launchctl unload /Library/LaunchAgents/com.github.$(GHUSER).$(BINARY).plist" && echo "  sudo rm -f /Library/LaunchAgents/com.github.$(GHUSER).$(BINARY).plist" || true
