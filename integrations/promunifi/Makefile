PACKAGES=`find ./cmd -mindepth 1 -maxdepth 1 -type d`
BINARY=unifi-poller
VERSION=`git tag -l --merged | tail -n1`

all: man unifi-poller

clean:
	for p in $(PACKAGES); do rm -f `echo $${p}|cut -d/ -f3`{.macos,.linux,.1,}{,.gz}; done
	for p in $(PACKAGES); do rm -f `echo $${p}|cut -d/ -f3`{_,-}*.{deb,rpm,pkg}; done
	rm -rf package_build

build: unifi-poller
unifi-poller:
	for p in $(PACKAGES); do go build -ldflags "-w -s -X main.Version=$(VERSION)" $${p}; done

linux: unifi-poller.linux
unifi-poller.linux:
	for p in $(PACKAGES); do GOOS=linux go build -o unifi-poller.linux -ldflags "-w -s -X main.Version=$(VERSION)" $${p}; done

darwin: unifi-poller.macos
unifi-poller.macos:
	for p in $(PACKAGES); do GOOS=darwin go build -o unifi-poller.macos -ldflags "-w -s -X main.Version=$(VERSION)" $${p}; done

test: lint
	for p in $(PACKAGES) $(LIBRARYS); do go test -race -covermode=atomic $${p}; done

lint:
	goimports -l $(PACKAGES)
	gofmt -l $(PACKAGES)
	errcheck $(PACKAGES)
	golint $(PACKAGES)
	go vet $(PACKAGES)

man: unifi-poller.1.gz
unifi-poller.1.gz:
	scripts/build_manpages.sh ./

rpm: man linux
	scripts/build_linux_packages.sh rpm

deb: man linux
	scripts/build_linux_packages.sh deb

osxpkg: man darwin
	scripts/build_osx_package.sh

install: man
	scripts/local_install.sh

uninstall:
	scripts/local_uninstall.sh

deps:
	dep ensure -update
