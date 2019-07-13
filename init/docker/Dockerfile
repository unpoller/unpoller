#
# This is part of Application Builder.
# https://github.com/golift/application-builder
#

ARG ARCH=amd64
ARG OS=linux
ARG BUILD_DATE=0
ARG COMMIT=0
ARG VERSION=unknown
ARG BINARY=application-builder
ARG GHREPO=golift/application-builder

FROM golang:stretch as builder
ARG ARCH
ARG OS
ARG BINARY
ARG GHREPO

RUN mkdir -p $GOPATH/pkg/mod $GOPATH/bin $GOPATH/src/github.com/${GHREPO}
RUN apt-get update \
  && apt-get install -y curl  \
  && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY . $GOPATH/src/github.com/${GHREPO}
WORKDIR $GOPATH/src/github.com/${GHREPO}

RUN dep ensure --vendor-only \
  && CGO_ENABLED=0 make ${BINARY}.${ARCH}.${OS}

FROM scratch
ARG ARCH
ARG OS
ARG BUILD_DATE
ARG COMMIT
ARG VERSION
ARG LICENSE=MIT
ARG BINARY
ARG GHREPO
ARG URL=http://github.com/golift/application-builder
ARG DESC=application-builder
ARG VENDOR=golift
ARG AUTHOR=golift
ARG CONFIG_FILE=config.conf

# Build-time metadata as defined at https://github.com/opencontainers/image-spec/blob/master/annotations.md
LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="${BINARY}" \
      org.opencontainers.image.documentation="https://github.com/${GHREPO}/wiki/Docker" \
      org.opencontainers.image.description="${DESC}" \
      org.opencontainers.image.url="${URL}" \
      org.opencontainers.image.revision="${COMMIT}" \
      org.opencontainers.image.source="https://github.com/${GHREPO}" \
      org.opencontainers.image.vendor="${VENDOR}" \
      org.opencontainers.image.authors="${AUTHOR}" \
      org.opencontainers.image.architecture="${OS} ${ARCH}" \
      org.opencontainers.image.licenses="${LICENSE}" \
      org.opencontainers.image.version="${VERSION}"

COPY --from=builder /go/src/github.com/${GHREPO}/${BINARY}.${ARCH}.${OS} /image
COPY --from=builder /go/src/github.com/${GHREPO}/examples/${CONFIG_FILE}.example /etc/${BINARY}/${CONFIG_FILE}

VOLUME [ "/etc/${BINARY}" ]
ENTRYPOINT [ "/image" ]
