# build binary
FROM golang:1.10.3 as builder
COPY . /go/src/github.com/michaeloverton/ddos-laser/
WORKDIR /go/src/github.com/michaeloverton/ddos-laser/cmd/sentinel/
RUN CGO_ENABLED=0 go build

# build final image
FROM alpine:latest
RUN apk add --update bash
COPY --from=builder /go/src/github.com/michaeloverton/ddos-laser/cmd/sentinel/sentinel /opt/
WORKDIR /opt
EXPOSE 3000
CMD ["/opt/sentinel"]
