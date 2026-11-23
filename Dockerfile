FROM busybox:latest@sha256:a46257e06c82ad85242d0d38077d0a52d9896389ca60a1a59d3463db586e5fa4 AS builder
# we have to do this hop because distroless is bare without common shell commands

RUN mkdir -p /etc/unpoller
# copy over example config for cnfg environment-based default config
COPY examples/up.conf.example /etc/unpoller/up.conf
COPY unpoller_manual.html /etc/unpoller/manual.html
COPY README.html /etc/unpoller/readme.html

FROM gcr.io/distroless/static-debian13@sha256:4c2cd1dee442c9a99924a378c7c126fac8fef14771579ef567c3f2bdd7a5122a
ARG TARGETPLATFORM

COPY ${TARGETPLATFORM}/unpoller /usr/bin/unpoller
COPY --from=builder /etc/unpoller /etc/unpoller

HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
  CMD ["/usr/bin/unpoller", "--health"]

ENTRYPOINT [ "/usr/bin/unpoller" ]
