#
# building static go binary with Debian golang container
#
ARG REPO=github.com/davidnewhall/unifi-poller

FROM golang:stretch as builder
ARG ARCH=amd64
ARG OS=linux
ARG REPO

RUN mkdir -p $GOPATH/pkg/mod $GOPATH/bin $GOPATH/src/${REPO}
RUN apt-get update \
  && apt-get install -y curl  \
  && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY . $GOPATH/src/${REPO}
WORKDIR $GOPATH/src/${REPO}

RUN dep ensure --vendor-only \
  && CGO_ENABLED=0 make unifi-poller.${ARCH}.${OS} \
  && mv unifi-poller.${ARCH}.${OS} unifi-poller

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

COPY --from=builder /go/src/${REPO}/unifi-poller /unifi-poller
COPY --from=builder /go/src/${REPO}/examples/up.conf.example /etc/unifi-poller}/up.conf

VOLUME [ "/etc/unifi-poller" ]
ENTRYPOINT [ "/unifi-poller" ]
