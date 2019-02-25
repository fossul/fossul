FROM golang:1.11

USER nobody

RUN mkdir -p /go/src/github.com/fossil
WORKDIR /go/src/github.com/fossil

COPY . /go/src/github.com/fossil
RUN go build

CMD ["./fossil"]