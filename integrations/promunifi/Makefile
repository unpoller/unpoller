PACKAGES=`find ./cmd -mindepth 1 -maxdepth 1 -type d`
BINARY=unifi-poller

all: clean man build

clean:
	for p in $(PACKAGES); do rm -f `echo $${p}|cut -d/ -f3`{,.1,.1.gz}; done

build:
	for p in $(PACKAGES); do go build -ldflags "-w -s" $${p}; done

linux:
	for p in $(PACKAGES); do GOOS=linux go build -ldflags "-w -s" $${p}; done

install: man test build
	@echo "If you get errors, you may need sudo."
	# Install binary. binary.
	GOBIN=/usr/local/bin go install -ldflags "-w -s" ./...
	# Making config folders and installing man page.
	mkdir -p /usr/local/etc/$(BINARY) /usr/local/share/man/man1
	mv *.1.gz /usr/local/share/man/man1
	# Installing config file, man page and launch agent or systemd unit file.
	test -f /usr/local/etc/$(BINARY)/up.conf || cp up.conf.example /usr/local/etc/$(BINARY)/up.conf
	test -d ~/Library/LaunchAgents && cp startup/launchd/com.github.davidnewhall.$(BINARY).plist ~/Library/LaunchAgents || true
	test -d /etc/systemd/system && cp startup/systemd/$(BINARY).service /etc/systemd/system || true
	# Making systemd happy by telling it to reload.
	test -x /bin/systemctl && /bin/systemctl --system daemon-reload || true
	@echo
	@echo "Installation Complete. Edit the config file @ /usr/local/etc/$(BINARY)/up.conf "
	@echo "Then start the daemon with:"
	@test -d ~/Library/LaunchAgents && echo "   launchctl load ~/Library/LaunchAgents/com.github.davidnewhall.$(BINARY).plist" || true
	@test -d /etc/systemd/system && echo "   sudo /bin/systemctl start $(BINARY)" || true
	@echo "Examine the log file at: /usr/local/var/log/$(BINARY).log (logs may go elsewhere on linux, check syslog)"

uninstall:
	@echo "If you get errors, you may need sudo."
	# Stopping the daemon
	test -x /bin/systemctl && /bin/systemctl stop $(BINARY) || true
	test -x /bin/launchctl && /bin/launchctl unload ~/Library/LaunchAgents/com.github.davidnewhall.$(BINARY).plist || true
	# Deleting config file, binary, man page, launch agent or unit file.
	rm -rf /usr/local/{etc,bin}/$(BINARY) /usr/local/share/man/man1/$(BINARY).1.gz
	rm -f ~/Library/LaunchAgents/com.github.davidnewhall.$(BINARY).plist
	rm -f /etc/systemd/system/$(BINARY).service
	# Making systemd happy by telling it to reload.
	test -x /bin/systemctl && /bin/systemctl --system daemon-reload || true

test: lint
	for p in $(PACKAGES) $(LIBRARYS); do go test -race -covermode=atomic $${p}; done

lint:
	goimports -l $(PACKAGES)
	gofmt -l $(PACKAGES)
	errcheck $(PACKAGES)
	golint $(PACKAGES)
	go vet $(PACKAGES)

man:
	script/build_manpages.sh ./

deps:
	dep ensure -update
