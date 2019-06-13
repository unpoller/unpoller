FROM golang:stretch as builder

RUN mkdir -p $GOPATH/pkg/mod $GOPATH/bin

RUN apt-get update \
  && apt-get install -y ruby ruby-dev curl  \
  && gem install --no-document fpm \
  && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh \
  && mkdir -p $GOPATH/src/github.com/davidnewhall

COPY . $GOPATH/src/github.com/davidnewhall/unifi-poller
WORKDIR $GOPATH/src/github.com/davidnewhall/unifi-poller

RUN dep ensure \
  && make dockerbin

FROM scratch 

COPY --from=builder /go/src/github.com/davidnewhall/unifi-poller/unifi-poller.dockerbin /unifi-poller
COPY --from=builder /go/src/github.com/davidnewhall/unifi-poller/examples/up.conf.example /etc/unifi-poller/up.conf

VOLUME [ "/etc/unifi-poller"]

ENTRYPOINT [ "/unifi-poller" ]
CMD [ "--config=/etc/unifi-poller/up.conf" ]