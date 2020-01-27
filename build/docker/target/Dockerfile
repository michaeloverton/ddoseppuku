# build binary
FROM golang:1.10.3 as builder
COPY . /go/src/github.com/michaeloverton/ddos-laser/
WORKDIR /go/src/github.com/michaeloverton/ddos-laser/cmd/target/
RUN CGO_ENABLED=0 go build

# build final image
FROM alpine:latest
RUN apk add --update bash ca-certificates
COPY --from=builder /go/src/github.com/michaeloverton/ddos-laser/cmd/target/target /opt/
WORKDIR /opt
EXPOSE 8080
CMD ["/opt/target"]
