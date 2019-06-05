PACKAGE=./cmd/unifi-poller
BINARY=unifi-poller
VERSION=`git tag -l --merged | tail -n1`

all: man unifi-poller

# Prepare a release. Called in Travis CI.
release: clean test man linux macos rpm deb
	mkdir -p build_assets
	gzip -9k unifi-poller.linux
	gzip -9k unifi-poller.macos
	mv unifi-poller.macos.gz unifi-poller.linux.gz build_assets/
	cp *.rpm *.deb build_assets/

clean:
	rm -f `echo $(PACKAGE)|cut -d/ -f3`{.macos,.linux,.1,}{,.gz}
	rm -f `echo $(PACKAGE)|cut -d/ -f3`{_,-}*.{deb,rpm,pkg}
	rm -rf package_build build_assets

build: unifi-poller
unifi-poller:
	go build -ldflags "-w -s -X main.Version=$(VERSION)" $(PACKAGE)

linux: unifi-poller.linux
unifi-poller.linux:
	GOOS=linux go build -o unifi-poller.linux -ldflags "-w -s -X main.Version=$(VERSION)" $(PACKAGE)

macos: unifi-poller.macos
unifi-poller.macos:
	GOOS=darwin go build -o unifi-poller.macos -ldflags "-w -s -X main.Version=$(VERSION)" $(PACKAGE)

test: lint
	go test -race -covermode=atomic $(PACKAGE)

lint:
	golangci-lint run --enable-all -D gochecknoglobals

man: unifi-poller.1.gz
unifi-poller.1.gz:
	scripts/build_manpages.sh ./

rpm: man linux
	scripts/build_linux_packages.sh rpm

deb: man linux
	scripts/build_linux_packages.sh deb

osxpkg: man macos
	scripts/build_osx_package.sh

install: man
	scripts/local_install.sh

uninstall:
	scripts/local_uninstall.sh

deps:
	dep ensure -update
