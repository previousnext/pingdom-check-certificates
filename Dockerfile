FROM golang:1.8
ADD . /go/src/github.com/previousnext/pingdom-check-certificates
WORKDIR /go/src/github.com/previousnext/pingdom-check-certificates
RUN go get github.com/mitchellh/gox
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/previousnext/pingdom-check-certificates/bin/pingdom-check-certificates_linux_amd64 /usr/local/bin/pingdom-check-certificates
CMD ["pingdom-check-certificates"]
