#
# building static go binary with Debian golang container
#
ARG BINARY=unifi-poller
ARG REPO=github.com/davidnewhall/${BINARY}

FROM golang:stretch as builder
ARG ARCH=amd64
ARG OS=linux
ARG BINARY
ARG REPO

RUN mkdir -p $GOPATH/pkg/mod $GOPATH/bin $GOPATH/src/${REPO}
RUN apt-get update \
  && apt-get install -y curl  \
  && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY . $GOPATH/src/${REPO}
WORKDIR $GOPATH/src/${REPO}

RUN dep ensure --vendor-only \
  && CGO_ENABLED=0 make ${BINARY}.${ARCH}.${OS} \
  && mv ${BINARY}.${ARCH}.${OS} ${BINARY}

#
# creating container for run
# to use this container use the following command:
#
# docker run -d -v /your/config/up.conf:/etc/unifi-poller/up.conf golift/unifi-poller
#
# by using "-e UNIFI_PASSWORD=your-secret-pasword" you can avoid this configuration in the config file
#

FROM scratch
ARG REPO
ARG BINARY

COPY --from=builder /go/src/${REPO}/${BINARY} /${BINARY}
COPY --from=builder /go/src/${REPO}/examples/up.conf.example /etc/${BINARY}/up.conf

VOLUME [ "/etc/${BINARY}" ]
ENTRYPOINT [ "/${BINARY}" ]
