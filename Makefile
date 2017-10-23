#!/usr/bin/make -f

export CGO_ENABLED=0

PROJECT=github.com/previousnext/pingdom-check-certificates

# Builds the project
build:
	gox -os='linux' -arch='amd64' -output='bin/pingdom-check-certificates_{{.OS}}_{{.Arch}}' -ldflags='-extldflags "-static"' $(PROJECT)

# Run all lint checking with exit codes for CI
lint:
	golint -set_exit_status `go list ./... | grep -v /vendor/`

# Run tests with coverage reporting
test:
	go test -cover main.go

IMAGE=previousnext/pingdom-check-certificates
VERSION=$(shell git describe --tags --always)

# Releases the project Docker Hub
release:
	docker build -t ${IMAGE}:${VERSION} .
	docker push ${IMAGE}:${VERSION}

.PHONY: build lint test release
