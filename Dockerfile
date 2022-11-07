FROM gcr.io/distroless/static-debian11

COPY unpoller /usr/bin/unpoller

ENTRYPOINT [ "/usr/bin/unpoller" ]
