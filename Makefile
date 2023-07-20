# This Makefile is written as generic as possible.
# Setting the variables in settings.sh and creating the paths in the repo makes this work.
# See more: https://github.com/golift/application-builder

# Suck in our application information.
IGNORED:=$(shell bash -c "source settings.sh ; env | grep -v BASH_FUNC | sed 's/=/:=/;s/^/export /' > .metadata.make")

# md2roff turns markdown into man files and html files.
MD2ROFF_BIN=github.com/davidnewhall/md2roff@v0.0.1

# rsrc adds an ico file to a Windows exe file.
RSRC_BIN=github.com/akavel/rsrc

# CI passes the version in. Local builds get it from the current git tag.
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

# Makefile targets follow.

all: build

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

build-and-release: export DOCKER_REGISTRY = ghcr.io
build-and-release: export DOCKER_IMAGE_NAME = unpoller/unpoller
build-and-release: export GPG_SIGNING_KEY = 

bulid-and-release: clean
	goreleaser release --clean

build: export DOCKER_REGISTRY = ghcr.io
build: export DOCKER_IMAGE_NAME = unpoller/unpoller
build: export GPG_SIGNING_KEY = 

build: clean
	goreleaser release --clean --skip-validate --skip-publish --skip-sign --debug

clean:
	git clean -xdf || true
	(docker images -f "dangling=true" -q | xargs docker rmi) || true

lint:
	golangci-lint run --fix

test:
	go test -timeout=30s ./...

integration-test:
	go test -timeout=30m -args=integration ./...
