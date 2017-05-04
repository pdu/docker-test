FROM golang

ADD . /go/src/github.com/pdu/docker-test

RUN go install github.com/pdu/docker-test/server

ENTRYPOINT ["/go/bin/server"]
CMD []

