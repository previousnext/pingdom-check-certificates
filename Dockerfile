FROM golang:1.9
ADD workspace /go
RUN go get github.com/mitchellh/gox
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/bin/pingdom-check-certificate_linux_amd64 /usr/local/bin/pingdom-check-certificate
CMD ["pingdom-check-certificate"]
