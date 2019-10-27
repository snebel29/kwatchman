SHELL := $(shell which bash)

CONTAINER_USER ?= kwatchman
VERSION        ?= master

BUILD_DATE = $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
COMMIT     = $(shell git rev-parse --short HEAD)

STABLE_VERSION_REGEX := ^v([0-9]{1,}\.){2}[0-9]{1,}$$

REPOSITORY = github.com/snebel29/kwatchman
PKG        = ${REPOSITORY}/internal/pkg

COVERAGE_FILE=/tmp/coverage.out

LD_FLAGS = "-X ${PKG}/version.Version=$(VERSION) \
			-X ${PKG}/version.BuildDate=$(BUILD_DATE) \
			-X ${PKG}/version.Commit=$(COMMIT) \
			-w -extldflags -static"

# To ensure that all built code and potentially non tested is covered by the race detector
report-race-conditions:
	go build -race cmd/*.go

build:
	CGO_ENABLED=0 go build -ldflags $(LD_FLAGS) cmd/*.go

test:
	go test -race ./... -coverprofile=coverage.txt -covermode=atomic

test-coverage-report:
	go test ./... -coverprofile=$(COVERAGE_FILE)
	go tool cover -html=$(COVERAGE_FILE)

          file: {{ coverage_report_filepath }}
clean:
	go clean

deps:
	dep ensure -v

docker-image:

	docker build -f build/Dockerfile \
		--build-arg VERSION=$(VERSION) \
		--build-arg REPOSITORY=$(REPOSITORY) \
		--build-arg PKG=$(PKG) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg CONTAINER_USER=$(CONTAINER_USER) \
		-t snebel29/kwatchman:$(VERSION) .

ifeq ($(shell echo $(VERSION) | egrep "$(STABLE_VERSION_REGEX)"),)
	@echo "Version $(VERSION) is not stable"

else
	@echo "Version $(VERSION) is stable therefore tag latest"
	docker tag snebel29/kwatchman:$(VERSION) snebel29/kwatchman:latest
endif

	docker image prune -f

push-docker-image:

ifeq ($(shell echo $(VERSION) | egrep "$(STABLE_VERSION_REGEX)"),)
	@echo "Version $(VERSION) is not stable"
else
	@echo "Version $(VERSION) is stable therefore push latest tag"
	docker push snebel29/kwatchman:latest
endif

	docker push snebel29/kwatchman:$(VERSION)

.PHONY: build test clean docker-image publish-docker-image test-coverage-report
