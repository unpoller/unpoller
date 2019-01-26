PACKAGES=`find ./cmd -mindepth 1 -maxdepth 1 -type d`
LIBRARYS=./unidev

all: clean man build

clean:
	for p in $(PACKAGES); do rm -f `echo $${p}|cut -d/ -f3`{,.1,.1.gz}; done

build:
	for p in $(PACKAGES); do go build -ldflags "-w -s" $${p}; done

linux:
	for p in $(PACKAGES); do GOOS=linux go build -ldflags "-w -s" $${p}; done

install: man test build
	@echo "If you get errors, you may need sudo."
	# Install binary.
	GOBIN=/usr/local/bin go install -ldflags "-w -s" ./...
	# Make folders and install man page.
	mkdir -p /usr/local/etc/unifi-poller /usr/local/share/man/man1
	mv *.1.gz /usr/local/share/man/man1
	# Install config file, man page and launch agent or systemd unit file.
	test -f /usr/local/etc/unifi-poller/up.conf || cp up.conf.example /usr/local/etc/unifi-poller/up.conf
	test -d ~/Library/LaunchAgents && cp startup/launchd/com.github.davidnewhall.unifi-poller.plist ~/Library/LaunchAgents || true
	test -d /etc/systemd/system && cp startup/systemd/unifi-poller.service /etc/systemd/system || true
	# Make systemd happy by telling it to reload.
	test -x /bin/systemctl && /bin/systemctl --system daemon-reload
	@echo "Installation Complete. Edit the config file @ /usr/local/etc/unifi-poller/up.conf "
	@echo "Then start the daemon with:"
	@test -d ~/Library/LaunchAgents && echo "   launchctl load ~/Library/LaunchAgents/com.github.davidnewhall.unifi-poller.plist"
	@test -d /etc/systemd/system && echo "   sudo /bin/systemctl start unifi-poller"
	@echo "Examine the log file at: /usr/local/var/log/unifi-poller.log (logs may go elsewhere on linux, check syslog)"

uninstall:
	@echo "If you get errors, you may need sudo."
	# Stop the daemon
	test -x /bin/systemctl && /bin/systemctl stop unifi-poller
	test -x /bin/launchctl && /bin/launchctl unload ~/Library/LaunchAgents/com.github.davidnewhall.unifi-poller.plist
	# Delete config file, binary, man page, launch agent or unit file.
	rm -rf /usr/local/{etc,bin}/unifi-poller /usr/local/share/man/man1/unifi-poller.1.gz
	rm -f ~/Library/LaunchAgents/com.github.davidnewhall.unifi-poller.plist
	rm -f /etc/systemd/system/unifi-poller.service
	# Make systemd happy by telling it to reload.
	test -x /bin/systemctl && /bin/systemctl --system daemon-reload

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
