#!/usr/bin/make -f

PROJECT=github.com/previousnext/pingdom-check-certificate

# Builds the project
build:
	gox -os='linux' -arch='amd64' -output='bin/pingdom-check-certificate_{{.OS}}_{{.Arch}}' -ldflags='-extldflags "-static"' $(PROJECT)

# Run all lint checking with exit codes for CI
lint:
	golint -set_exit_status `go list ./... | grep -v /vendor/`

# Run tests with coverage reporting
test:
	go test -cover main.go

IMAGE=previousnext/pingdom-check-certificate
VERSION=$(shell git describe --tags --always)

# Releases the project Docker Hub
release:
	docker build -t ${IMAGE}:${VERSION} .
	docker push ${IMAGE}:${VERSION}

.PHONY: build lint test release
