# This Makefile is written as generic as possible.
# Setting the variables in settings.sh and creating the paths in the repo makes this work.
# See more: https://github.com/golift/application-builder

# Suck in our application information.
IGNORED:=$(shell bash -c "source settings.sh ; env | grep -v BASH_FUNC | sed 's/=/:=/;s/^/export /' > .metadata.make")

# md2roff turns markdown into man files and html files.
MD2ROFF_BIN=github.com/davidnewhall/md2roff@v0.0.1

# rsrc adds an ico file to a Windows exe file.
RSRC_BIN=github.com/akavel/rsrc

# If upx is available, use it to compress the binaries.
UPXPATH=$(shell which upx)

# Skip upx in Mac ARM environments: https://github.com/upx/upx/issues/446
ifeq ($(shell uname -ps),Darwin arm)
	UPXPATH=
endif

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

all: clean generate build

####################
##### Releases #####
####################

# Prepare a release. Called in Travis CI.
release: clean windows
	# Prepareing a release!
	mkdir -p $@
	mv $(BINARY).*.linux $(BINARY).*.freebsd $@/
	gzip -9r $@/
	for i in $(BINARY)*.exe ; do zip -9qj $@/$$i.zip $$i examples/*.example *.html; rm -f $$i;done
	mv *.rpm *.deb *.txz $@/
	# Generating File Hashes
	openssl dgst -r -sha256 $@/* | sed 's#release/##' | tee $@/checksums.sha256.txt

# DMG only makes a DMG file is MACAPP is set. Otherwise, it makes a gzipped binary for macOS.
dmg: clean macapp
	mkdir -p release
	[ "$(MACAPP)" = "" ] || hdiutil create release/$(MACAPP).dmg -srcfolder $(MACAPP).app -ov
	[ "$(MACAPP)" != "" ] || mv $(BINARY).*.macos release/
	[ "$(MACAPP)" != "" ] || gzip -9r release/
	openssl dgst -r -sha256 release/* | sed 's#release/##' | tee release/macos_checksum.sha256.txt

# Delete all build assets.
clean:
	rm -rf dist
	rm -f $(BINARY) $(BINARY).*.{macos,freebsd,linux,exe,upx}{,.gz,.zip} $(BINARY).1{,.gz} $(BINARY).rb
	rm -f $(BINARY){_,-}*.{deb,rpm,txz} v*.tar.gz.sha256 examples/MANUAL .metadata.make rsrc_*.syso
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
	cd /tmp ; go get $(RSRC_BIN) ; go install $(RSRC_BIN)@latest

####################
##### Binaries #####
####################

build: $(BINARY)
$(BINARY): main.go
	go build -o $(BINARY) -ldflags "-w -s $(VERSION_LDFLAGS) $(EXTRA_LDFLAGS) "
	[ -z "$(UPXPATH)" ] || $(UPXPATH) -q9 $@

macos: $(BINARY).amd64.macos
$(BINARY).amd64.macos: main.go
	# Building darwin 64-bit x86 binary.
	GOOS=darwin GOARCH=amd64 go build -o $@ -ldflags "-w -s $(VERSION_LDFLAGS) $(EXTRA_LDFLAGS) "
	[ -z "$(UPXPATH)" ] || $(UPXPATH) -q9 $@

exe: $(BINARY).amd64.exe
windows: $(BINARY).amd64.exe
$(BINARY).amd64.exe: rsrc.syso main.go
	# Building windows 64-bit x86 binary.
	GOOS=windows GOARCH=amd64 go build -o $@ -ldflags "-w -s $(VERSION_LDFLAGS) $(EXTRA_LDFLAGS) $(WINDOWS_LDFLAGS)"
	[ -z "$(UPXPATH)" ] || $(UPXPATH) -q9 $@

####################
##### Packages #####
####################

macapp: $(MACAPP).app
$(MACAPP).app: macos
	[ -z "$(MACAPP)" ] || mkdir -p init/macos/$(MACAPP).app/Contents/MacOS
	[ -z "$(MACAPP)" ] || cp $(BINARY).amd64.macos init/macos/$(MACAPP).app/Contents/MacOS/$(MACAPP)
	[ -z "$(MACAPP)" ] || cp -rp init/macos/$(MACAPP).app $(MACAPP).app

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
	/usr/bin/install -m 0755 -d $(PREFIX)/lib/$(BINARY)/web/static/{css,js,images}
	/usr/bin/install -m 0644 -cp init/webserver/index.html $(PREFIX)/lib/$(BINARY)/web/index.html
	/usr/bin/install -m 0644 -cp init/webserver/static/css/* $(PREFIX)/lib/$(BINARY)/web/static/css/
	/usr/bin/install -m 0644 -cp init/webserver/static/js/* $(PREFIX)/lib/$(BINARY)/web/static/js/
	/usr/bin/install -m 0644 -cp init/webserver/static/images/* $(PREFIX)/lib/$(BINARY)/web/static/images/
	[ -f $(ETC)/$(BINARY)/$(CONFIG_FILE) ] || /usr/bin/install -m 0644 -cp  examples/$(CONFIG_FILE).example $(ETC)/$(BINARY)/$(CONFIG_FILE)
	/usr/bin/install -m 0644 -cp LICENSE *.html examples/* $(PREFIX)/share/doc/$(BINARY)/

goreleaser:
	goreleaser release --rm-dist

goreleaser-test:
	goreleaser release --rm-dist --skip-validate --skip-publish --skip-sign --debug
