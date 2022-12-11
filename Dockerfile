FROM busybox:latest as builder
# we have to do this hop because distroless is bare without common shell commands

RUN mkdir -p /etc/unpoller
# copy over example config for cnfg environment-based default config
COPY examples/up.conf.example /etc/unpoller/up.conf
COPY unpoller_manual.html /etc/unpoller/manual.html
COPY README.html /etc/unpoller/readme.html

FROM gcr.io/distroless/static-debian11

COPY unpoller /usr/bin/unpoller
COPY --from=builder /etc/unpoller /etc/unpoller

ENTRYPOINT [ "/usr/bin/unpoller" ]
