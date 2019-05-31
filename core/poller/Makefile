PACKAGES=`find ./cmd -mindepth 1 -maxdepth 1 -type d`
BINARY=unifi-poller
VERSION=`git tag -l --merged | tail -n1`

all: clean man build

clean:
	for p in $(PACKAGES); do rm -f `echo $${p}|cut -d/ -f3`{,.1,.1.gz}; done
	rm -rf package_build unifi-poller_*.deb unifi-poller-*.rpm unifi-poller-*.pkg
	rm -f unifi-poller.*.gz

build:
	for p in $(PACKAGES); do go build -ldflags "-w -s -X main.Version=$(VERSION)" $${p}; done

linux:
	for p in $(PACKAGES); do GOOS=linux go build -ldflags "-w -s -X main.Version=$(VERSION)" $${p}; done

darwin:
	for p in $(PACKAGES); do GOOS=darwin go build -ldflags "-w -s -X main.Version=$(VERSION)" $${p}; done

test: lint
	for p in $(PACKAGES) $(LIBRARYS); do go test -race -covermode=atomic $${p}; done

man:
	scripts/build_manpages.sh ./

rpm: man linux
	scripts/build_linux_packages.sh rpm

deb: man linux
	scripts/build_linux_packages.sh deb

osxpkg: man darwin
	scripts/build_osx_package.sh

install: all
	scripts/local_install.sh

uninstall:
	scripts/local_uninstall.sh

lint:
	goimports -l $(PACKAGES)
	gofmt -l $(PACKAGES)
	errcheck $(PACKAGES)
	golint $(PACKAGES)
	go vet $(PACKAGES)

deps:
	dep ensure -update
