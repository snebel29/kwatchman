SHELL := $(shell which bash)

CONTAINER_USER ?= kwatchman
VERSION        ?= development

STABLE_VERSION_REGEX := ^v([0-9]{1,}\.){2}[0-9]{1,}$$

REPOSITORY=github.com/snebel29/kwatchman
COVERAGE_FILE=/tmp/coverage.out
LD_FLAGS="-X ${REPOSITORY}/internal/pkg/cli.Version=$(VERSION) -w -extldflags -static"


report-race-conditions:
	go build -race cmd/*.go

build: deps
	CGO_ENABLED=0 go build -ldflags $(LD_FLAGS) cmd/*.go

test:
	go test -race ./... -cover

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
		-t snebel29/kwatchman:$(VERSION) .

ifeq ($(shell echo $(VERSION) | egrep "$(STABLE_VERSION_REGEX)"),)
	@echo "Version $(VERSION) is not stable"

else
	@echo "Version $(VERSION) is stable therefore tag latest"
	docker tag snebel29/kwatchman:$(VERSION) snebel29/kwatchman:latest
endif

	docker image prune -f

push-docker-image:

	echo $(shell echo $(VERSION) | egrep "$(STABLE_VERSION_REGEX)")

ifeq ($(shell echo $(VERSION) | egrep "$(STABLE_VERSION_REGEX)"),)
	@echo "Version $(VERSION) is not stable"
else
	@echo "Version $(VERSION) is stable therefore push latest tag"
	docker push snebel29/kwatchman:latest
endif

	docker push snebel29/kwatchman:$(VERSION)

.PHONY: build test clean docker-image publish-docker-image test-coverage-report
