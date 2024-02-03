#---Build stage---
FROM golang:1.20 AS builder
RUN mkdir -p /go/src/ipam/
COPY . /go/src/ipam/
WORKDIR /go/src/ipam/cmd/ipam-service-server
RUN go mod tidy
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags='-w -s' -o /go/bin/ipam-service

#---Final stage---
FROM alpine:latest
COPY --from=builder /go/bin/ipam-service /go/bin/ipam-service
RUN mkdir -p /migrations
COPY --from=builder /go/src/ipam/migrations /migrations
CMD /go/bin/ipam-service 

