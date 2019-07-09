# This Makefile is written as generic as possible.
# Setting the variables in .metadata.sh and creating the paths in the repo makes this work.
#

# Suck in our application information.
IGNORE := $(shell bash -c "source .metadata.sh ; env | sed 's/=/:=/;s/^/export /' > .metadata.make")

# md2roff turns markdown into man files and html files.
MD2ROFF_BIN=github.com/github/hub/md2roff-bin

# Travis CI passes the version in. Local builds get it from the current git tag.
ifeq ($(VERSION),)
	include .metadata.make
else
	# Preserve the passed-in version & iteration (homebrew).
	_VERSION:=$(VERSION)
	_ITERATION:=$(ITERATION)
	include .metadata.make
	VERSION:=$(_VERSION)
	ITERATION:=$(_ITERATION)
endif

# rpm is wierd and changes - to _ in versions.
RPMVERSION:=$(shell echo $(VERSION) | tr -- - _)

# Makefile targets follow.

all: build

# Prepare a release. Called in Travis CI.
release: clean vendor test macos arm windows linux_packages
	# Prepareing a release!
	mkdir -p $@
	mv $(BINARY).*.macos $(BINARY).*.linux $@/
	gzip -9r $@/
	for i in $(BINARY)*.exe; do zip -9qm $@/$$i.zip $$i;done
	mv *.rpm *.deb $@/
	# Generating File Hashes
	openssl dgst -r -sha256 $@/* | sed 's#release/##' | tee $@/checksums.sha256.txt

# Delete all build assets.
clean:
	# Cleaning up.
	rm -f $(BINARY) $(BINARY).*.{macos,linux,exe}{,.gz,.zip} $(BINARY).1{,.gz} $(BINARY).rb
	rm -f $(BINARY){_,-}*.{deb,rpm} v*.tar.gz.sha256 .metadata.make
	rm -f cmd/$(BINARY)/README{,.html} README{,.html} ./$(BINARY)_manual.html
	rm -rf package_build_* release

# Build a man page from a markdown file using md2roff.
# This also turns the repo readme into an html file.
# md2roff is needed to build the man file and html pages from the READMEs.
man: $(BINARY).1.gz
$(BINARY).1.gz: md2roff
	# Building man page. Build dependency first: md2roff
	go run $(MD2ROFF_BIN) --manual $(BINARY) --version $(VERSION) --date "$(DATE)" examples/MANUAL.md
	gzip -9nc examples/MANUAL > $@
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
	go build -o $(BINARY) -ldflags "-w -s -X $(VERSION_PATH)=$(VERSION)-$(ITERATION)"

linux: $(BINARY).amd64.linux
$(BINARY).amd64.linux:
	# Building linux 64-bit x86 binary.
	GOOS=linux GOARCH=amd64 go build -o $@ -ldflags "-w -s -X $(VERSION_PATH)=$(VERSION)-$(ITERATION)"

linux386: $(BINARY).i386.linux
$(BINARY).i386.linux:
	# Building linux 32-bit x86 binary.
	GOOS=linux GOARCH=386 go build -o $@ -ldflags "-w -s -X $(VERSION_PATH)=$(VERSION)-$(ITERATION)"

arm: arm64 armhf

arm64: $(BINARY).arm64.linux
$(BINARY).arm64.linux:
	# Building linux 64-bit ARM binary.
	GOOS=linux GOARCH=arm64 go build -o $@ -ldflags "-w -s -X $(VERSION_PATH)=$(VERSION)-$(ITERATION)"

armhf: $(BINARY).armhf.linux
$(BINARY).armhf.linux:
	# Building linux 32-bit ARM binary.
	GOOS=linux GOARCH=arm GOARM=6 go build -o $@ -ldflags "-w -s -X $(VERSION_PATH)=$(VERSION)-$(ITERATION)"

macos: $(BINARY).amd64.macos
$(BINARY).amd64.macos:
	# Building darwin 64-bit x86 binary.
	GOOS=darwin GOARCH=amd64 go build -o $@ -ldflags "-w -s -X $(VERSION_PATH)=$(VERSION)-$(ITERATION)"

exe: $(BINARY).amd64.exe
windows: $(BINARY).amd64.exe
$(BINARY).amd64.exe:
	# Building windows 64-bit x86 binary.
	GOOS=windows GOARCH=amd64 go build -o $@ -ldflags "-w -s -X $(VERSION_PATH)=$(VERSION)-$(ITERATION)"

# Packages

linux_packages: rpm deb rpm386 deb386 debarm rpmarm debarmhf rpmarmhf

rpm: $(BINARY)-$(RPMVERSION)-$(ITERATION).x86_64.rpm
$(BINARY)-$(RPMVERSION)-$(ITERATION).x86_64.rpm: package_build_linux check_fpm
	@echo "Building 'rpm' package for $(BINARY) version '$(RPMVERSION)-$(ITERATION)'."
	fpm -s dir -t rpm \
		--architecture x86_64 \
		--name $(BINARY) \
		--rpm-os linux \
		--version $(RPMVERSION) \
		--iteration $(ITERATION) \
		--after-install scripts/after-install.sh \
		--before-remove scripts/before-remove.sh \
		--license $(LICENSE) \
		--url $(URL) \
		--maintainer "$(MAINT)" \
		--description "$(DESC)" \
		--config-files "/etc/$(BINARY)/$(CONFIG_FILE)" \
		--chdir $<

deb: $(BINARY)_$(VERSION)-$(ITERATION)_amd64.deb
$(BINARY)_$(VERSION)-$(ITERATION)_amd64.deb: package_build_linux check_fpm
	@echo "Building 'deb' package for $(BINARY) version '$(VERSION)-$(ITERATION)'."
	fpm -s dir -t deb \
		--architecture amd64 \
		--name $(BINARY) \
		--deb-no-default-config-files \
		--version $(VERSION) \
		--iteration $(ITERATION) \
		--after-install scripts/after-install.sh \
		--before-remove scripts/before-remove.sh \
		--license $(LICENSE) \
		--url $(URL) \
		--maintainer "$(MAINT)" \
		--description "$(DESC)" \
		--config-files "/etc/$(BINARY)/$(CONFIG_FILE)" \
		--chdir $<

rpm386: $(BINARY)-$(RPMVERSION)-$(ITERATION).i386.rpm
$(BINARY)-$(RPMVERSION)-$(ITERATION).i386.rpm: package_build_linux_386 check_fpm
	@echo "Building 32-bit 'rpm' package for $(BINARY) version '$(RPMVERSION)-$(ITERATION)'."
	fpm -s dir -t rpm \
		--architecture i386 \
		--name $(BINARY) \
		--rpm-os linux \
		--version $(RPMVERSION) \
		--iteration $(ITERATION) \
		--after-install scripts/after-install.sh \
		--before-remove scripts/before-remove.sh \
		--license $(LICENSE) \
		--url $(URL) \
		--maintainer "$(MAINT)" \
		--description "$(DESC)" \
		--config-files "/etc/$(BINARY)/$(CONFIG_FILE)" \
		--chdir $<

deb386: $(BINARY)_$(VERSION)-$(ITERATION)_i386.deb
$(BINARY)_$(VERSION)-$(ITERATION)_i386.deb: package_build_linux_386 check_fpm
	@echo "Building 32-bit 'deb' package for $(BINARY) version '$(VERSION)-$(ITERATION)'."
	fpm -s dir -t deb \
		--architecture i386 \
		--name $(BINARY) \
		--deb-no-default-config-files \
		--version $(VERSION) \
		--iteration $(ITERATION) \
		--after-install scripts/after-install.sh \
		--before-remove scripts/before-remove.sh \
		--license $(LICENSE) \
		--url $(URL) \
		--maintainer "$(MAINT)" \
		--description "$(DESC)" \
		--config-files "/etc/$(BINARY)/$(CONFIG_FILE)" \
		--chdir $<

rpmarm: $(BINARY)-$(RPMVERSION)-$(ITERATION).arm64.rpm
$(BINARY)-$(RPMVERSION)-$(ITERATION).arm64.rpm: package_build_linux_arm64 check_fpm
	@echo "Building 64-bit ARM8 'rpm' package for $(BINARY) version '$(RPMVERSION)-$(ITERATION)'."
	fpm -s dir -t rpm \
		--architecture arm64 \
		--name $(BINARY) \
		--rpm-os linux \
		--version $(RPMVERSION) \
		--iteration $(ITERATION) \
		--after-install scripts/after-install.sh \
		--before-remove scripts/before-remove.sh \
		--license $(LICENSE) \
		--url $(URL) \
		--maintainer "$(MAINT)" \
		--description "$(DESC)" \
		--config-files "/etc/$(BINARY)/$(CONFIG_FILE)" \
		--chdir $<

debarm: $(BINARY)_$(VERSION)-$(ITERATION)_arm64.deb
$(BINARY)_$(VERSION)-$(ITERATION)_arm64.deb: package_build_linux_arm64 check_fpm
	@echo "Building 64-bit ARM8 'deb' package for $(BINARY) version '$(VERSION)-$(ITERATION)'."
	fpm -s dir -t deb \
		--architecture arm64 \
		--name $(BINARY) \
		--deb-no-default-config-files \
		--version $(VERSION) \
		--iteration $(ITERATION) \
		--after-install scripts/after-install.sh \
		--before-remove scripts/before-remove.sh \
		--license $(LICENSE) \
		--url $(URL) \
		--maintainer "$(MAINT)" \
		--description "$(DESC)" \
		--config-files "/etc/$(BINARY)/$(CONFIG_FILE)" \
		--chdir $<

rpmarmhf: $(BINARY)-$(RPMVERSION)-$(ITERATION).armhf.rpm
$(BINARY)-$(RPMVERSION)-$(ITERATION).armhf.rpm: package_build_linux_armhf check_fpm
	@echo "Building 32-bit ARM6/7 HF 'rpm' package for $(BINARY) version '$(RPMVERSION)-$(ITERATION)'."
	fpm -s dir -t rpm \
		--architecture armhf \
		--name $(BINARY) \
		--rpm-os linux \
		--version $(RPMVERSION) \
		--iteration $(ITERATION) \
		--after-install scripts/after-install.sh \
		--before-remove scripts/before-remove.sh \
		--license $(LICENSE) \
		--url $(URL) \
		--maintainer "$(MAINT)" \
		--description "$(DESC)" \
		--config-files "/etc/$(BINARY)/$(CONFIG_FILE)" \
		--chdir $<

debarmhf: $(BINARY)_$(VERSION)-$(ITERATION)_armhf.deb
$(BINARY)_$(VERSION)-$(ITERATION)_armhf.deb: package_build_linux_armhf check_fpm
	@echo "Building 32-bit ARM6/7 HF 'deb' package for $(BINARY) version '$(VERSION)-$(ITERATION)'."
	fpm -s dir -t deb \
		--architecture armhf \
		--name $(BINARY) \
		--deb-no-default-config-files \
		--version $(VERSION) \
		--iteration $(ITERATION) \
		--after-install scripts/after-install.sh \
		--before-remove scripts/before-remove.sh \
		--license $(LICENSE) \
		--url $(URL) \
		--maintainer "$(MAINT)" \
		--description "$(DESC)" \
		--config-files "/etc/$(BINARY)/$(CONFIG_FILE)" \
		--chdir $<

docker:
	docker build -f init/docker/Dockerfile \
		--build-arg "BUILD_DATE=$(DATE)" \
		--build-arg "COMMIT=$(COMMIT)" \
		--build-arg "VERSION=$(VERSION)-$(ITERATION)" \
		--build-arg "LICENSE=$(LICENSE)" \
		--build-arg "DESC=$(DESC)" \
		--build-arg "URL=$(URL)" \
		--build-arg "VENDOR=$(VENDOR)" \
		--build-arg "AUTHOR=$(MAINT)" \
		--build-arg "BINARY=$(BINARY)" \
		--build-arg "GHREPO=$(GHREPO)" \
		--build-arg "CONFIG_FILE=$(CONFIG_FILE)" \
		--tag $(DHUSER)/$(BINARY):local .

# Build an environment that can be packaged for linux.
package_build_linux: readme man linux
	# Building package environment for linux.
	mkdir -p $@/usr/bin $@/etc/$(BINARY) $@/usr/share/man/man1 $@/usr/share/doc/$(BINARY)
	# Copying the binary, config file, unit file, and man page into the env.
	cp $(BINARY).amd64.linux $@/usr/bin/$(BINARY)
	cp *.1.gz $@/usr/share/man/man1
	cp examples/$(CONFIG_FILE).example $@/etc/$(BINARY)/
	cp examples/$(CONFIG_FILE).example $@/etc/$(BINARY)/$(CONFIG_FILE)
	cp LICENSE *.html examples/*?.?* $@/usr/share/doc/$(BINARY)/
	[ ! -f init/systemd/template.unit.service ] || mkdir -p $@/lib/systemd/system
	[ ! -f init/systemd/template.unit.service ] || \
		sed -e "s/{{BINARY}}/$(BINARY)/g" -e "s/{{DESC}}/$(DESC)/g" \
		init/systemd/template.unit.service > $@/lib/systemd/system/$(BINARY).service

package_build_linux_386: package_build_linux linux386
	mkdir -p $@
	cp -r $</* $@/
	cp $(BINARY).i386.linux $@/usr/bin/$(BINARY)

package_build_linux_arm64: package_build_linux arm64
	mkdir -p $@
	cp -r $</* $@/
	cp $(BINARY).arm64.linux $@/usr/bin/$(BINARY)

package_build_linux_armhf: package_build_linux armhf
	mkdir -p $@
	cp -r $</* $@/
	cp $(BINARY).armhf.linux $@/usr/bin/$(BINARY)

check_fpm:
	@fpm --version > /dev/null || (echo "FPM missing. Install FPM: https://fpm.readthedocs.io/en/latest/installing.html" && false)

# This builds a Homebrew formula file that can be used to install this app from source.
# The source used comes from the released version on GitHub. This will not work with local source.
# This target is used by Travis CI to update the released Forumla when a new tag is created.
formula: $(BINARY).rb
v$(VERSION).tar.gz.sha256:
	# Calculate the SHA from the Github source file.
	curl -sL $(URL)/archive/v$(VERSION).tar.gz | openssl dgst -r -sha256 | tee $@
$(BINARY).rb: v$(VERSION).tar.gz.sha256
	# Creating formula from template using sed.
	sed -e "s/{{Version}}/$(VERSION)/g" \
		-e "s/{{Iter}}/$(ITERATION)/g" \
		-e "s/{{SHA256}}/$(shell head -c64 $<)/g" \
		-e "s/{{Desc}}/$(DESC)/g" \
		-e "s%{{URL}}%$(URL)%g" \
		-e "s%{{GHREPO}}%$(GHREPO)%g" \
		-e "s%{{CONFIG_FILE}}%$(CONFIG_FILE)%g" \
		-e "s%{{Class}}%$(shell echo $(BINARY) | perl -pe 's/(?:\b|-)(\p{Ll})/\u$$1/g')%g" \
		init/homebrew/formula.rb.tmpl | tee $(BINARY).rb

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
	dep ensure --vendor-only

# Don't run this unless you're ready to debug untested vendored dependencies.
deps:
	dep ensure --update

# Homebrew stuff. macOS only.

# Used for Homebrew only. Other distros can create packages.
install: man readme $(BINARY)
	@echo -  Done Building!  -
	@echo -  Local installation with the Makefile is only supported on macOS.
	@echo If you wish to install the application manually on Linux, check out the wiki: $(URL)/wiki/Installation
	@echo -  Otherwise, build and install a package: make rpm -or- make deb
	@echo See the Package Install wiki for more info: $(URL)/wiki/Package-Install
	@[ "$(shell uname)" = "Darwin" ] || (echo "Unable to continue, not a Mac." && false)
	@[ "$(PREFIX)" != "" ] || (echo "Unable to continue, PREFIX not set. Use: make install PREFIX=/usr/local ETC=/usr/local/etc" && false)
	@[ "$(ETC)" != "" ] || (echo "Unable to continue, ETC not set. Use: make install PREFIX=/usr/local ETC=/usr/local/etc" && false)
	# Copying the binary, config file, unit file, and man page into the env.
	/usr/bin/install -m 0755 -d $(PREFIX)/bin $(PREFIX)/share/man/man1 $(ETC)/$(BINARY) $(PREFIX)/share/doc/$(BINARY)
	/usr/bin/install -m 0755 -cp $(BINARY) $(PREFIX)/bin/$(BINARY)
	/usr/bin/install -m 0644 -cp $(BINARY).1.gz $(PREFIX)/share/man/man1
	/usr/bin/install -m 0644 -cp examples/$(CONFIG_FILE).example $(ETC)/$(BINARY)/
	[ -f $(ETC)/$(BINARY)/$(CONFIG_FILE) ] || /usr/bin/install -m 0644 -cp  examples/$(CONFIG_FILE).example $(ETC)/$(BINARY)/$(CONFIG_FILE)
	/usr/bin/install -m 0644 -cp LICENSE *.html examples/* $(PREFIX)/share/doc/$(BINARY)/

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
