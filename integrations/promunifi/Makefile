BINARY:=unifi-poller
PACKAGE:=./cmd/$(BINARY)
VERSION:=$(shell git tag -l --merged | tail -n1 | tr -d v)
ITERATION:=$(shell git rev-list --all --count)

all: man build

# Prepare a release. Called in Travis CI.
release: clean test man linux macos rpm deb osxpkg
	mkdir -p release
	gzip -9 $(BINARY).linux
	gzip -9 $(BINARY).macos
	mv $(BINARY).macos.gz $(BINARY).linux.gz release/
	mv *.rpm *.deb *.pkg release/

# Delete all build assets.
clean:
	rm -f $(BINARY){.macos,.linux,.1,}{,.gz}
	rm -f $(BINARY){_,-}*.{deb,rpm,pkg}
	rm -rf package_build release

# Build a man page from a markdown file using ronn.
man: $(BINARY).1.gz
$(BINARY).1.gz:
	scripts/build_manpages.sh ./

# Binaries

build: $(BINARY)
$(BINARY):
	go build -o $(BINARY) -ldflags "-w -s -X main.Version=$(VERSION)" $(PACKAGE)

linux: $(BINARY).linux
$(BINARY).linux:
	GOOS=linux go build -o $(BINARY).linux -ldflags "-w -s -X main.Version=$(VERSION)" $(PACKAGE)

macos: $(BINARY).macos
$(BINARY).macos:
	GOOS=darwin go build -o $(BINARY).macos -ldflags "-w -s -X main.Version=$(VERSION)" $(PACKAGE)

# Packages

rpm: man linux $(BINARY)-$(VERSION)-$(ITERATION).x86_64.rpm
$(BINARY)-$(VERSION)-$(ITERATION).x86_64.rpm:
	scripts/build_packages.sh rpm "$(VERSION)" "$(ITERATION)"

deb: man linux $(BINARY)_$(VERSION)-$(ITERATION)_amd64.deb
$(BINARY)_$(VERSION)-$(ITERATION)_amd64.deb:
	scripts/build_packages.sh deb "$(VERSION)" "$(ITERATION)"

osxpkg: man macos $(BINARY)-$(VERSION).pkg
$(BINARY)-$(VERSION).pkg:
	scripts/build_packages.sh osxpkg "$(VERSION)" "$(ITERATION)"

# Extras

test: lint
	go test -race -covermode=atomic $(PACKAGE)
lint:
	golangci-lint run --enable-all -D gochecknoglobals

install: man
	scripts/local_install.sh

uninstall:
	scripts/local_uninstall.sh

deps:
	dep ensure -update
