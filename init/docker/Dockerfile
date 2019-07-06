#
# building static go binary with Debian golang container
#
ARG ARCH=amd64

FROM golang:stretch as builder
ARG ARCH

RUN mkdir -p $GOPATH/pkg/mod $GOPATH/bin $GOPATH/src/github.com/davidnewhall
RUN apt-get update \
  && apt-get install -y curl  \
  && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY . $GOPATH/src/github.com/davidnewhall/unifi-poller
WORKDIR $GOPATH/src/github.com/davidnewhall/unifi-poller

RUN dep ensure --vendor-only \
  && CGO_ENABLED=0 make unifi-poller.${ARCH}.linux

#
# creating container for run
# to use this container use the following command:
#
# docker run -d -v /your/config/up.conf:/etc/unifi-poller/up.conf golift/unifi-poller
#
# by using "-e UNIFI_PASSWORD=your-secret-pasword" you can avoid this configuration in the config file
#
FROM scratch
ARG ARCH

COPY --from=builder /go/src/github.com/davidnewhall/unifi-poller/unifi-poller.${ARCH}.linux /unifi-poller
COPY --from=builder /go/src/github.com/davidnewhall/unifi-poller/examples/up.conf.example /etc/unifi-poller/up.conf

VOLUME [ "/etc/unifi-poller" ]
ENTRYPOINT [ "/unifi-poller" ]
