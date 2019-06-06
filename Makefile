SHELL := $(shell which bash)

CONTAINER_USER ?= kwatchman
VERSION        ?= x.x.x-development

REPOSITORY=github.com/snebel29/kwatchman
COVERAGE_FILE=/tmp/coverage.out
LD_FLAGS="-X ${REPOSITORY}/internal/pkg/cli.Version=$(VERSION) -w -extldflags -static"


build: deps
	CGO_ENABLED=0 go build -ldflags $(LD_FLAGS) cmd/*.go

test:
	go test ./... -cover

test-coverage-report:
	go test ./... -coverprofile=$(COVERAGE_FILE)
	go tool cover -html=$(COVERAGE_FILE)

clean:
	go clean

deps:
	dep ensure -v

docker-image:
	docker build -f build/Dockerfile \
		--build-arg VERSION=$(VERSION) \
		--build-arg REPOSITORY=$(REPOSITORY) \
		--build-arg CONTAINER_USER=$(CONTAINER_USER) \
		-t snebel29/kwatchman:latest \
		-t snebel29/kwatchman:$(VERSION) .

push-docker-image:
	docker push snebel29/kwatchman:$(VERSION)
	docker push snebel29/kwatchman:latest

.PHONY: build test clean docker-image publish-docker-image test-coverage-report
