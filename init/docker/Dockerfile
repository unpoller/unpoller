#
# building static go binary with Debian golang container
#

ARG ARCH=amd64
ARG OS=linux
ARG BUILD_DATE=0
ARG COMMIT=0
ARG VERSION=development

FROM golang:stretch as builder
ARG ARCH
ARG OS

RUN mkdir -p $GOPATH/pkg/mod $GOPATH/bin $GOPATH/src/github.com/davidnewhall/unifi-poller
RUN apt-get update \
  && apt-get install -y curl  \
  && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY . $GOPATH/src/github.com/davidnewhall/unifi-poller
WORKDIR $GOPATH/src/github.com/davidnewhall/unifi-poller

RUN dep ensure --vendor-only \
  && CGO_ENABLED=0 make unifi-poller.${ARCH}.${OS}

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
ARG OS
ARG BUILD_DATE
ARG COMMIT
ARG VERSION
ARG LICENSE=MIT
ARG TITLE=unifi-poller
ARG URL=http://github.com/davidnewhall/unifi-poller
ARG DESC=unifi-poller
ARG VENDOR=golift
ARG AUTHOR=golift

# Build-time metadata as defined at https://github.com/opencontainers/image-spec/blob/master/annotations.md
LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="${TITLE}" \
      org.opencontainers.image.documentation="${URL}/wiki/Docker" \
      org.opencontainers.image.description="${DESC}" \
      org.opencontainers.image.url="${URL}" \
      org.opencontainers.image.revision="${COMMIT}" \
      org.opencontainers.image.source="${URL}" \
      org.opencontainers.image.vendor="${VENDOR}" \
      org.opencontainers.image.authors="${AUTHOR}" \
      org.opencontainers.image.architecture="${OS} ${ARCH}" \
      org.opencontainers.image.licenses="${LICENSE}" \
      org.opencontainers.image.version="${VERSION}"

COPY --from=builder /go/src/github.com/davidnewhall/unifi-poller/unifi-poller.${ARCH}.${OS} /unifi-poller
COPY --from=builder /go/src/github.com/davidnewhall/unifi-poller/examples/up.conf.example /etc/unifi-poller/up.conf

VOLUME [ "/etc/unifi-poller" ]
ENTRYPOINT [ "/unifi-poller" ]
