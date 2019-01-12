PACKAGES=`find ./cmd -mindepth 1 -maxdepth 1 -type d`
LIBRARYS=./unidev

all: clean man build

clean:
	for p in $(PACKAGES); do rm -f `echo $${p}|cut -d/ -f3`{,.1,.1.gz}; done

build:
	for p in $(PACKAGES); do go build -ldflags "-w -s" $${p}; done

linux:
	for p in $(PACKAGES); do GOOS=linux go build -ldflags "-w -s" $${p}; done

install: man
	@echo "If you get errors, you may need sudo."
	GOBIN=/usr/local/bin go install -ldflags "-w -s" ./...
	mkdir -p /usr/local/etc/unifi-poller /usr/local/share/man/man1
	test -f /usr/local/etc/unifi-poller/up.conf || cp up.conf.example /usr/local/etc/unifi-poller/up.conf
	test -d ~/Library/LaunchAgents && cp startup/launchd/com.github.davidnewhall.unifi-poller.plist ~/Library/LaunchAgents || true
	test -d /etc/systemd/system && cp startup/systemd/unifi-poller.service /etc/systemd/system || true
	mv *.1.gz /usr/local/share/man/man1

uninstall:
	@echo "If you get errors, you may need sudo."
	test -f ~/Library/LaunchAgents/com.github.davidnewhall.unifi-poller.plist && launchctl unload ~/Library/LaunchAgents/com.github.davidnewhall.unifi-poller.plist || true
	test -f /etc/systemd/system/unifi-poller.service && systemctl stop unifi-poller || true
	rm -rf /usr/local/{etc,bin}/unifi-poller /usr/local/share/man/man1/unifi-poller.1.gz
	rm -f ~/Library/LaunchAgents/com.github.davidnewhall.unifi-poller.plist
	rm -f /etc/systemd/system/unifi-poller.service

test: lint
	for p in $(PACKAGES) $(LIBRARYS); do go test -race -covermode=atomic $${p}; done

lint:
	goimports -l $(PACKAGES) $(LIBRARYS)
	gofmt -l $(PACKAGES) $(LIBRARYS)
	errcheck $(PACKAGES) $(LIBRARYS)
	golint $(PACKAGES) $(LIBRARYS)
	go vet $(PACKAGES) $(LIBRARYS)

man:
	script/build_manpages.sh ./

deps:
	dep ensure -update
