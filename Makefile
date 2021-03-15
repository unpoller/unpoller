# This Makefile is written as generic as possible.
# Setting the variables in settings.sh and creating the paths in the repo makes this work.
# See more: https://github.com/golift/application-builder

# Suck in our application information.
IGNORED:=$(shell bash -c "source settings.sh ; env | grep -v BASH_FUNC | sed 's/=/:=/;s/^/export /' > .metadata.make")

# md2roff turns markdown into man files and html files.
MD2ROFF_BIN=github.com/davidnewhall/md2roff

# rsrc adds an ico file to a Windows exe file.
RSRC_BIN=github.com/akavel/rsrc

# If upx is available, use it to compress the binaries.
UPXPATH=$(shell which upx)

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
# used for freebsd packages.
BINARYU:=$(shell echo $(BINARY) | tr -- - _)

PACKAGE_SCRIPTS=
ifeq ($(FORMULA),service)
	PACKAGE_SCRIPTS=--after-install after-install-rendered.sh --before-remove before-remove-rendered.sh
endif

define PACKAGE_ARGS
$(PACKAGE_SCRIPTS) \
--name $(BINARY) \
--deb-no-default-config-files \
--rpm-os linux \
--iteration $(ITERATION) \
--license $(LICENSE) \
--url $(SOURCE_URL) \
--maintainer "$(MAINT)" \
--vendor "$(VENDOR)" \
--description "$(DESC)" \
--config-files "/etc/$(BINARY)/$(CONFIG_FILE)" \
--freebsd-origin "$(SOURCE_URL)"
endef

PLUGINS:=$(patsubst plugins/%/main.go,%,$(wildcard plugins/*/main.go))

VERSION_LDFLAGS:= -X \"$(VERSION_PATH).Branch=$(BRANCH) ($(COMMIT))\" \
	-X \"$(VERSION_PATH).BuildDate=$(DATE)\" \
	-X \"$(VERSION_PATH).BuildUser=$(shell whoami)\" \
  -X \"$(VERSION_PATH).Revision=$(ITERATION)\" \
  -X \"$(VERSION_PATH).Version=$(VERSION)\"

# Makefile targets follow.

all: clean build

####################
##### Releases #####
####################

# Prepare a release. Called in Travis CI.
release: clean linux_packages freebsd_packages windows
	# Prepareing a release!
	mkdir -p $@
	mv $(BINARY).*.linux $(BINARY).*.freebsd $@/
	gzip -9r $@/
	for i in $(BINARY)*.exe ; do zip -9qj $@/$$i.zip $$i examples/*.example *.html; rm -f $$i;done
	mv *.rpm *.deb *.txz $@/
	# Generating File Hashes
	openssl dgst -r -sha256 $@/* | sed 's#release/##' | tee $@/checksums.sha256.txt

dmg: clean macapp
	mkdir -p release
	hdiutil create release/$(MACAPP).dmg -srcfolder $(MACAPP).app -ov
	openssl dgst -r -sha256 release/* | sed 's#release/##' | tee release/dmg_checksum.sha256.txt

# Delete all build assets.
clean:
	rm -f $(BINARY) $(BINARY).*.{macos,freebsd,linux,exe,upx}{,.gz,.zip} $(BINARY).1{,.gz} $(BINARY).rb
	rm -f $(BINARY){_,-}*.{deb,rpm,txz} v*.tar.gz.sha256 examples/MANUAL .metadata.make
	rm -f cmd/$(BINARY)/README{,.html} README{,.html} ./$(BINARY)_manual.html rsrc.syso $(MACAPP).app.zip
	rm -rf package_build_* release after-install-rendered.sh before-remove-rendered.sh $(MACAPP).app

####################
##### Sidecars #####
####################

# Build a man page from a markdown file using md2roff.
# This also turns the repo readme into an html file.
# md2roff is needed to build the man file and html pages from the READMEs.
man: $(BINARY).1.gz
$(BINARY).1.gz: md2roff
	# Building man page. Build dependency first: md2roff
	$(shell go env GOPATH)/bin/md2roff --manual $(BINARY) --version $(VERSION) --date "$(DATE)" examples/MANUAL.md
	gzip -9nc examples/MANUAL > $@
	mv examples/MANUAL.html $(BINARY)_manual.html

md2roff: $(shell go env GOPATH)/bin/md2roff
$(shell go env GOPATH)/bin/md2roff:
	cd /tmp ; go get $(MD2ROFF_BIN) ; go install $(MD2ROFF_BIN)

# TODO: provide a template that adds the date to the built html file.
readme: README.html
README.html: md2roff
	# This turns README.md into README.html
	$(shell go env GOPATH)/bin/md2roff --manual $(BINARY) --version $(VERSION) --date "$(DATE)" README.md

rsrc: rsrc.syso
rsrc.syso: init/windows/application.ico init/windows/manifest.xml $(shell go env GOPATH)/bin/rsrc
	$(shell go env GOPATH)/bin/rsrc -ico init/windows/application.ico -manifest init/windows/manifest.xml
$(shell go env GOPATH)/bin/rsrc:
	cd /tmp ; go get $(RSRC_BIN) ; go install $(RSRC_BIN)

####################
##### Binaries #####
####################

build: $(BINARY)
$(BINARY): main.go
	go build -o $(BINARY) -ldflags "-w -s $(VERSION_LDFLAGS) $(EXTRA_LDFLAGS) "
	[ -z "$(UPXPATH)" ] || $(UPXPATH) -q9 $@

linux: $(BINARY).amd64.linux
$(BINARY).amd64.linux: main.go
	# Building linux 64-bit x86 binary.
	GOOS=linux GOARCH=amd64 go build -o $@ -ldflags "-w -s $(VERSION_LDFLAGS) $(EXTRA_LDFLAGS) "
	[ -z "$(UPXPATH)" ] || $(UPXPATH) -q9 $@

linux386: $(BINARY).i386.linux
$(BINARY).i386.linux: main.go
	# Building linux 32-bit x86 binary.
	GOOS=linux GOARCH=386 go build -o $@ -ldflags "-w -s $(VERSION_LDFLAGS) $(EXTRA_LDFLAGS) "
	[ -z "$(UPXPATH)" ] || $(UPXPATH) -q9 $@

arm: arm64 armhf

arm64: $(BINARY).arm64.linux
$(BINARY).arm64.linux: main.go
	# Building linux 64-bit ARM binary.
	GOOS=linux GOARCH=arm64 go build -o $@ -ldflags "-w -s $(VERSION_LDFLAGS) $(EXTRA_LDFLAGS) "
	[ -z "$(UPXPATH)" ] || $(UPXPATH) -q9 $@

armhf: $(BINARY).armhf.linux
$(BINARY).armhf.linux: main.go
	# Building linux 32-bit ARM binary.
	GOOS=linux GOARCH=arm GOARM=6 go build -o $@ -ldflags "-w -s $(VERSION_LDFLAGS) $(EXTRA_LDFLAGS) "
	[ -z "$(UPXPATH)" ] || $(UPXPATH) -q9 $@

macos: $(BINARY).amd64.macos
$(BINARY).amd64.macos: main.go
	# Building darwin 64-bit x86 binary.
	GOOS=darwin GOARCH=amd64 go build -o $@ -ldflags "-w -s $(VERSION_LDFLAGS) $(EXTRA_LDFLAGS) "
	[ -z "$(UPXPATH)" ] || $(UPXPATH) -q9 $@

freebsd: $(BINARY).amd64.freebsd
$(BINARY).amd64.freebsd: main.go
	GOOS=freebsd GOARCH=amd64 go build -o $@ -ldflags "-w -s $(VERSION_LDFLAGS) $(EXTRA_LDFLAGS) "

freebsd386: $(BINARY).i386.freebsd
$(BINARY).i386.freebsd: main.go
	GOOS=freebsd GOARCH=386 go build -o $@ -ldflags "-w -s $(VERSION_LDFLAGS) $(EXTRA_LDFLAGS) "
	[ -z "$(UPXPATH)" ] || $(UPXPATH) -q9 $@ || true

freebsdarm: $(BINARY).armhf.freebsd
$(BINARY).armhf.freebsd: main.go
	GOOS=freebsd GOARCH=arm go build -o $@ -ldflags "-w -s $(VERSION_LDFLAGS) $(EXTRA_LDFLAGS) "

exe: $(BINARY).amd64.exe
windows: $(BINARY).amd64.exe
$(BINARY).amd64.exe: rsrc.syso main.go
	# Building windows 64-bit x86 binary.
	GOOS=windows GOARCH=amd64 go build -o $@ -ldflags "-w -s $(VERSION_LDFLAGS) $(EXTRA_LDFLAGS) $(WINDOWS_LDFLAGS)"
	[ -z "$(UPXPATH)" ] || $(UPXPATH) -q9 $@

####################
##### Packages #####
####################

linux_packages: rpm deb rpm386 deb386 debarm rpmarm debarmhf rpmarmhf

freebsd_packages: freebsd_pkg freebsd386_pkg freebsdarm_pkg

macapp: $(MACAPP).app
$(MACAPP).app: macos
	@[ "$(MACAPP)" != "" ] || (echo "Must set 'MACAPP' in settings.sh!" && exit 1)
	mkdir -p init/macos/$(MACAPP).app/Contents/MacOS
	cp $(BINARY).amd64.macos init/macos/$(MACAPP).app/Contents/MacOS/$(MACAPP)
	cp -rp init/macos/$(MACAPP).app $(MACAPP).app

rpm: $(BINARY)-$(RPMVERSION)-$(ITERATION).x86_64.rpm
$(BINARY)-$(RPMVERSION)-$(ITERATION).x86_64.rpm: package_build_linux check_fpm
	@echo "Building 'rpm' package for $(BINARY) version '$(RPMVERSION)-$(ITERATION)'."
	fpm -s dir -t rpm $(PACKAGE_ARGS) -a x86_64 -v $(RPMVERSION) -C $<
	[ "$(SIGNING_KEY)" == "" ] || expect -c "spawn rpmsign --key-id=$(SIGNING_KEY) --resign $(BINARY)-$(RPMVERSION)-$(ITERATION).x86_64.rpm; expect -exact \"Enter pass phrase: \"; send \"$(PRIVATE_KEY)\r\"; expect eof"

deb: $(BINARY)_$(VERSION)-$(ITERATION)_amd64.deb
$(BINARY)_$(VERSION)-$(ITERATION)_amd64.deb: package_build_linux check_fpm
	@echo "Building 'deb' package for $(BINARY) version '$(VERSION)-$(ITERATION)'."
	fpm -s dir -t deb $(PACKAGE_ARGS) -a amd64 -v $(VERSION) -C $<
	[ "$(SIGNING_KEY)" == "" ] || expect -c "spawn debsigs --default-key="$(SIGNING_KEY)" --sign=origin $(BINARY)_$(VERSION)-$(ITERATION)_amd64.deb; expect -exact \"Enter passphrase: \"; send \"$(PRIVATE_KEY)\r\"; expect eof"

rpm386: $(BINARY)-$(RPMVERSION)-$(ITERATION).i386.rpm
$(BINARY)-$(RPMVERSION)-$(ITERATION).i386.rpm: package_build_linux_386 check_fpm
	@echo "Building 32-bit 'rpm' package for $(BINARY) version '$(RPMVERSION)-$(ITERATION)'."
	fpm -s dir -t rpm $(PACKAGE_ARGS) -a i386 -v $(RPMVERSION) -C $<
	[ "$(SIGNING_KEY)" == "" ] || expect -c "spawn rpmsign --key-id=$(SIGNING_KEY) --resign $(BINARY)-$(RPMVERSION)-$(ITERATION).i386.rpm; expect -exact \"Enter pass phrase: \"; send \"$(PRIVATE_KEY)\r\"; expect eof"

deb386: $(BINARY)_$(VERSION)-$(ITERATION)_i386.deb
$(BINARY)_$(VERSION)-$(ITERATION)_i386.deb: package_build_linux_386 check_fpm
	@echo "Building 32-bit 'deb' package for $(BINARY) version '$(VERSION)-$(ITERATION)'."
	fpm -s dir -t deb $(PACKAGE_ARGS) -a i386 -v $(VERSION) -C $<
	[ "$(SIGNING_KEY)" == "" ] || expect -c "spawn debsigs --default-key="$(SIGNING_KEY)" --sign=origin $(BINARY)_$(VERSION)-$(ITERATION)_i386.deb; expect -exact \"Enter passphrase: \"; send \"$(PRIVATE_KEY)\r\"; expect eof"

rpmarm: $(BINARY)-$(RPMVERSION)-$(ITERATION).arm64.rpm
$(BINARY)-$(RPMVERSION)-$(ITERATION).arm64.rpm: package_build_linux_arm64 check_fpm
	@echo "Building 64-bit ARM8 'rpm' package for $(BINARY) version '$(RPMVERSION)-$(ITERATION)'."
	fpm -s dir -t rpm $(PACKAGE_ARGS) -a arm64 -v $(RPMVERSION) -C $<
	[ "$(SIGNING_KEY)" == "" ] || expect -c "spawn rpmsign --key-id=$(SIGNING_KEY) --resign $(BINARY)-$(RPMVERSION)-$(ITERATION).arm64.rpm; expect -exact \"Enter pass phrase: \"; send \"$(PRIVATE_KEY)\r\"; expect eof"

debarm: $(BINARY)_$(VERSION)-$(ITERATION)_arm64.deb
$(BINARY)_$(VERSION)-$(ITERATION)_arm64.deb: package_build_linux_arm64 check_fpm
	@echo "Building 64-bit ARM8 'deb' package for $(BINARY) version '$(VERSION)-$(ITERATION)'."
	fpm -s dir -t deb $(PACKAGE_ARGS) -a arm64 -v $(VERSION) -C $<
	[ "$(SIGNING_KEY)" == "" ] || expect -c "spawn debsigs --default-key="$(SIGNING_KEY)" --sign=origin $(BINARY)_$(VERSION)-$(ITERATION)_arm64.deb; expect -exact \"Enter passphrase: \"; send \"$(PRIVATE_KEY)\r\"; expect eof"

rpmarmhf: $(BINARY)-$(RPMVERSION)-$(ITERATION).armhf.rpm
$(BINARY)-$(RPMVERSION)-$(ITERATION).armhf.rpm: package_build_linux_armhf check_fpm
	@echo "Building 32-bit ARM6/7 HF 'rpm' package for $(BINARY) version '$(RPMVERSION)-$(ITERATION)'."
	fpm -s dir -t rpm $(PACKAGE_ARGS) -a armhf -v $(RPMVERSION) -C $<
	[ "$(SIGNING_KEY)" == "" ] || expect -c "spawn rpmsign --key-id=$(SIGNING_KEY) --resign $(BINARY)-$(RPMVERSION)-$(ITERATION).armhf.rpm; expect -exact \"Enter pass phrase: \"; send \"$(PRIVATE_KEY)\r\"; expect eof"

debarmhf: $(BINARY)_$(VERSION)-$(ITERATION)_armhf.deb
$(BINARY)_$(VERSION)-$(ITERATION)_armhf.deb: package_build_linux_armhf check_fpm
	@echo "Building 32-bit ARM6/7 HF 'deb' package for $(BINARY) version '$(VERSION)-$(ITERATION)'."
	fpm -s dir -t deb $(PACKAGE_ARGS) -a armhf -v $(VERSION) -C $<
	[ "$(SIGNING_KEY)" == "" ] || expect -c "spawn debsigs --default-key="$(SIGNING_KEY)" --sign=origin $(BINARY)_$(VERSION)-$(ITERATION)_armhf.deb; expect -exact \"Enter passphrase: \"; send \"$(PRIVATE_KEY)\r\"; expect eof"

freebsd_pkg: $(BINARY)-$(VERSION)_$(ITERATION).amd64.txz
$(BINARY)-$(VERSION)_$(ITERATION).amd64.txz: package_build_freebsd check_fpm
	@echo "Building 'freebsd pkg' package for $(BINARY) version '$(VERSION)-$(ITERATION)'."
	fpm -s dir -t freebsd $(PACKAGE_ARGS) -a amd64 -v $(VERSION) -p $(BINARY)-$(VERSION)_$(ITERATION).amd64.txz -C $<

freebsd386_pkg: $(BINARY)-$(VERSION)_$(ITERATION).i386.txz
$(BINARY)-$(VERSION)_$(ITERATION).i386.txz: package_build_freebsd_386 check_fpm
	@echo "Building 32-bit 'freebsd pkg' package for $(BINARY) version '$(VERSION)-$(ITERATION)'."
	fpm -s dir -t freebsd $(PACKAGE_ARGS) -a 386 -v $(VERSION) -p $(BINARY)-$(VERSION)_$(ITERATION).i386.txz -C $<

freebsdarm_pkg: $(BINARY)-$(VERSION)_$(ITERATION).armhf.txz
$(BINARY)-$(VERSION)_$(ITERATION).armhf.txz: package_build_freebsd_arm check_fpm
	@echo "Building 32-bit ARM6/7 HF 'freebsd pkg' package for $(BINARY) version '$(VERSION)-$(ITERATION)'."
	fpm -s dir -t freebsd $(PACKAGE_ARGS) -a arm -v $(VERSION) -p $(BINARY)-$(VERSION)_$(ITERATION).armhf.txz -C $<

# Build an environment that can be packaged for linux.
package_build_linux: readme man plugins_linux_amd64 after-install-rendered.sh before-remove-rendered.sh linux
	# Building package environment for linux.
	mkdir -p $@/usr/bin $@/etc/$(BINARY) $@/usr/share/man/man1 $@/usr/share/doc/$(BINARY) $@/usr/lib/$(BINARY)
	# Copying the binary, config file, unit file, and man page into the env.
	cp $(BINARY).amd64.linux $@/usr/bin/$(BINARY)
	cp *.1.gz $@/usr/share/man/man1
	rm -f $@/usr/lib/$(BINARY)/*.so
	[ ! -f *amd64.so ] || cp *amd64.so $@/usr/lib/$(BINARY)/
	cp examples/$(CONFIG_FILE).example $@/etc/$(BINARY)/
	cp examples/$(CONFIG_FILE).example $@/etc/$(BINARY)/$(CONFIG_FILE)
	cp LICENSE *.html examples/*?.?* $@/usr/share/doc/$(BINARY)/
	[ "$(FORMULA)" != "service" ] || mkdir -p $@/lib/systemd/system
	[ "$(FORMULA)" != "service" ] || \
		sed -e "s/{{BINARY}}/$(BINARY)/g" -e "s/{{DESC}}/$(DESC)/g" \
		init/systemd/template.unit.service > $@/lib/systemd/system/$(BINARY).service

after-install-rendered.sh:
	sed -e "s/{{BINARY}}/$(BINARY)/g" scripts/after-install.sh > after-install-rendered.sh

before-remove-rendered.sh:
	sed -e "s/{{BINARY}}/$(BINARY)/g" scripts/before-remove.sh > before-remove-rendered.sh

package_build_linux_386: package_build_linux linux386
	mkdir -p $@
	cp -r $</* $@/
	[ ! -f *386.so ] || cp *386.so $@/usr/lib/$(BINARY)/
	cp $(BINARY).i386.linux $@/usr/bin/$(BINARY)

package_build_linux_arm64: package_build_linux arm64
	mkdir -p $@
	cp -r $</* $@/
	[ ! -f *arm64.so ] || cp *arm64.so $@/usr/lib/$(BINARY)/
	cp $(BINARY).arm64.linux $@/usr/bin/$(BINARY)

package_build_linux_armhf: package_build_linux armhf
	mkdir -p $@
	cp -r $</* $@/
	[ ! -f *armhf.so ] || cp *armhf.so $@/usr/lib/$(BINARY)/
	cp $(BINARY).armhf.linux $@/usr/bin/$(BINARY)

# Build an environment that can be packaged for freebsd.
package_build_freebsd: readme man after-install-rendered.sh before-remove-rendered.sh freebsd
	mkdir -p $@/usr/local/bin $@/usr/local/etc/$(BINARY) $@/usr/local/share/man/man1 $@/usr/local/share/doc/$(BINARY)
	cp $(BINARY).amd64.freebsd $@/usr/local/bin/$(BINARY)
	cp *.1.gz $@/usr/local/share/man/man1
	cp examples/$(CONFIG_FILE).example $@/usr/local/etc/$(BINARY)/
	cp examples/$(CONFIG_FILE).example $@/usr/local/etc/$(BINARY)/$(CONFIG_FILE)
	cp LICENSE *.html examples/*?.?* $@/usr/local/share/doc/$(BINARY)/
	[ "$(FORMULA)" != "service" ] || mkdir -p $@/usr/local/etc/rc.d
	[ "$(FORMULA)" != "service" ] || \
			sed -e "s/{{BINARY}}/$(BINARY)/g" -e "s/{{BINARYU}}/$(BINARYU)/g" -e "s/{{CONFIG_FILE}}/$(CONFIG_FILE)/g" \
			init/bsd/freebsd.rc.d > $@/usr/local/etc/rc.d/$(BINARY)
	[ "$(FORMULA)" != "service" ] || chmod +x $@/usr/local/etc/rc.d/$(BINARY)

package_build_freebsd_386: package_build_freebsd freebsd386
	mkdir -p $@
	cp -r $</* $@/
	cp $(BINARY).i386.freebsd $@/usr/local/bin/$(BINARY)

package_build_freebsd_arm: package_build_freebsd freebsdarm
	mkdir -p $@
	cp -r $</* $@/
	cp $(BINARY).armhf.freebsd $@/usr/local/bin/$(BINARY)

check_fpm:
	@fpm --version > /dev/null || (echo "FPM missing. Install FPM: https://fpm.readthedocs.io/en/latest/installing.html" && false)

##################
##### Extras #####
##################

plugins: $(patsubst %,%.so,$(PLUGINS))
$(patsubst %,%.so,$(PLUGINS)):
	go build -o $@ -ldflags "$(VERSION_LDFLAGS)" -buildmode=plugin ./plugins/$(patsubst %.so,%,$@)

linux_plugins: plugins_linux_amd64 plugins_linux_i386 plugins_linux_arm64 plugins_linux_armhf
plugins_linux_amd64: $(patsubst %,%.linux_amd64.so,$(PLUGINS))
$(patsubst %,%.linux_amd64.so,$(PLUGINS)):
	GOOS=linux GOARCH=amd64 go build -o $@ -ldflags "$(VERSION_LDFLAGS)" -buildmode=plugin ./plugins/$(patsubst %.linux_amd64.so,%,$@)

plugins_darwin: $(patsubst %,%.darwin.so,$(PLUGINS))
$(patsubst %,%.darwin.so,$(PLUGINS)):
	GOOS=darwin go build -o $@ -ldflags "$(VERSION_LDFLAGS)" -buildmode=plugin ./plugins/$(patsubst %.darwin.so,%,$@)

# Run code tests and lint.
test: lint
	# Testing.
	go test -race -covermode=atomic ./...
lint:
	# Checking lint.
	GOOS=linux $(shell go env GOPATH)/bin/golangci-lint run $(GOLANGCI_LINT_ARGS)
	GOOS=freebsd $(shell go env GOPATH)/bin/golangci-lint run $(GOLANGCI_LINT_ARGS)
	GOOS=windows $(shell go env GOPATH)/bin/golangci-lint run $(GOLANGCI_LINT_ARGS)

# Mockgen and bindata are examples.
# Your `go generate` may require other tools; add them!

mockgen: $(shell go env GOPATH)/bin/mockgen
$(shell go env GOPATH)/bin/mockgen:
	cd /tmp ; go get github.com/golang/mock/mockgen ; go install github.com/golang/mock/mockgen

bindata: $(shell go env GOPATH)/bin/go-bindata
$(shell go env GOPATH)/bin/go-bindata:
	cd /tmp ; go get -u github.com/go-bindata/go-bindata/... ; go install github.com/go-bindata/go-bindata

generate: mockgen bindata
	go generate ./...

##################
##### Docker #####
##################

docker:
	docker build -f init/docker/Dockerfile \
		--build-arg "BUILD_DATE=$(DATE)" \
		--build-arg "COMMIT=$(COMMIT)" \
		--build-arg "VERSION=$(VERSION)-$(ITERATION)" \
		--build-arg "LICENSE=$(LICENSE)" \
		--build-arg "DESC=$(DESC)" \
		--build-arg "VENDOR=$(VENDOR)" \
		--build-arg "AUTHOR=$(MAINT)" \
		--build-arg "BINARY=$(BINARY)" \
		--build-arg "SOURCE_URL=$(SOURCE_URL)" \
		--build-arg "CONFIG_FILE=$(CONFIG_FILE)" \
		--tag $(BINARY) .

####################
##### Homebrew #####
####################

# This builds a Homebrew formula file that can be used to install this app from source.
# The source used comes from the released version on GitHub. This will not work with local source.
# This target is used by Travis CI to update the released Forumla when a new tag is created.
formula: $(BINARY).rb
v$(VERSION).tar.gz.sha256:
	# Calculate the SHA from the Github source file.
	curl -sL $(SOURCE_URL)/archive/v$(VERSION).tar.gz | openssl dgst -r -sha256 | tee $@
$(BINARY).rb: v$(VERSION).tar.gz.sha256 init/homebrew/$(FORMULA).rb.tmpl
	# Creating formula from template using sed.
	sed -e "s/{{Version}}/$(VERSION)/g" \
		-e "s/{{Iter}}/$(ITERATION)/g" \
		-e "s/{{SHA256}}/$(shell head -c64 $<)/g" \
		-e "s/{{Desc}}/$(DESC)/g" \
		-e "s%{{SOURCE_URL}}%$(SOURCE_URL)%g" \
		-e "s%{{SOURCE_PATH}}%$(SOURCE_PATH)%g" \
		-e "s%{{CONFIG_FILE}}%$(CONFIG_FILE)%g" \
		-e "s%{{Class}}%$(shell echo $(BINARY) | perl -pe 's/(?:\b|-)(\p{Ll})/\u$$1/g')%g" \
		init/homebrew/$(FORMULA).rb.tmpl | tee $(BINARY).rb
		# That perl line turns hello-world into HelloWorld, etc.

# Used for Homebrew only. Other distros can create packages.
install: man readme $(BINARY) plugins_darwin
	@echo -  Done Building  -
	@echo -  Local installation with the Makefile is only supported on macOS.
	@echo -  Otherwise, build and install a package: make rpm -or- make deb
	@[ "$(shell uname)" = "Darwin" ] || (echo "Unable to continue, not a Mac." && false)
	@[ "$(PREFIX)" != "" ] || (echo "Unable to continue, PREFIX not set. Use: make install PREFIX=/usr/local ETC=/usr/local/etc" && false)
	@[ "$(ETC)" != "" ] || (echo "Unable to continue, ETC not set. Use: make install PREFIX=/usr/local ETC=/usr/local/etc" && false)
	# Copying the binary, config file, unit file, and man page into the env.
	/usr/bin/install -m 0755 -d $(PREFIX)/bin $(PREFIX)/share/man/man1 $(ETC)/$(BINARY) $(PREFIX)/share/doc/$(BINARY) $(PREFIX)/lib/$(BINARY)
	/usr/bin/install -m 0755 -cp $(BINARY) $(PREFIX)/bin/$(BINARY)
	/usr/bin/install -m 0644 -cp $(BINARY).1.gz $(PREFIX)/share/man/man1
	/usr/bin/install -m 0644 -cp examples/$(CONFIG_FILE).example $(ETC)/$(BINARY)/
	[ -f $(ETC)/$(BINARY)/$(CONFIG_FILE) ] || /usr/bin/install -m 0644 -cp  examples/$(CONFIG_FILE).example $(ETC)/$(BINARY)/$(CONFIG_FILE)
	/usr/bin/install -m 0644 -cp LICENSE *.html examples/* $(PREFIX)/share/doc/$(BINARY)/
