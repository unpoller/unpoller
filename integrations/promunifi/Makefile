PACKAGES=`find ./cmd -mindepth 1 -maxdepth 1 -type d`
LIBRARYS=

all: clean man build

clean:
	for p in $(PACKAGES); do rm -f `echo $${p}|cut -d/ -f3`{,.1,.1.gz}; done

build:
	for p in $(PACKAGES); do go build -ldflags "-w -s" $${p}; done

linux:
	for p in $(PACKAGES); do GOOS=linux go build -ldflags "-w -s" $${p}; done

install:
	go install -ldflags "-w -s" ./...

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
	rm -rf Godeps vendor
	godep save ./...
	godep update ./...
