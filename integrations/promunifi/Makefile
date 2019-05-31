PACKAGES=`find ./cmd -mindepth 1 -maxdepth 1 -type d`
BINARY=unifi-poller

all: clean test man build

clean:
	for p in $(PACKAGES); do rm -f `echo $${p}|cut -d/ -f3`{,.1,.1.gz}; done
	rm -rf package_build unifi-poller_*.deb unifi-poller-*.rpm unifi-poller-*.pkg

build:
	for p in $(PACKAGES); do go build -ldflags "-w -s" $${p}; done

linux:
	for p in $(PACKAGES); do GOOS=linux go build -ldflags "-w -s" $${p}; done

darwin:
	for p in $(PACKAGES); do GOOS=darwin go build -ldflags "-w -s" $${p}; done

test: lint
	for p in $(PACKAGES) $(LIBRARYS); do go test -race -covermode=atomic $${p}; done

man:
	scripts/build_manpages.sh ./

rpm: clean test man linux
	scripts/build_linux_packages.sh rpm

deb: clean test man linux
	scripts/build_linux_packages.sh deb

osxpkg: clean test man darwin
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
