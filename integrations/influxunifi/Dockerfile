#
# building static go binary with golang container
#
FROM golang:stretch as builder

RUN mkdir -p $GOPATH/pkg/mod $GOPATH/bin

RUN apt-get update \
  && apt-get install -y curl  \
  && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh \
  && mkdir -p $GOPATH/src/github.com/davidnewhall

COPY . $GOPATH/src/github.com/davidnewhall/unifi-poller
WORKDIR $GOPATH/src/github.com/davidnewhall/unifi-poller

RUN dep ensure \
  && CGO_ENABLED=0 make linux

#
# creating container for run 
# to use this container use the following command: 
#
# docker run -d -v /your/config/up.conf:/etc/unifi-poller/up.conf golift/unifi-poller
#
# by using "-e UNIFI_PASSWORD=your-secret-pasword" you can avoid this configuration in the config file
#
FROM scratch 

COPY --from=builder /go/src/github.com/davidnewhall/unifi-poller/unifi-poller.linux /unifi-poller
COPY --from=builder /go/src/github.com/davidnewhall/unifi-poller/examples/up.conf.example /etc/unifi-poller/up.conf

VOLUME [ "/etc/unifi-poller"]

ENTRYPOINT [ "/unifi-poller" ]
