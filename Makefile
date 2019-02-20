VERSION=x.x.x-development
REPOSITORY=github.com/snebel29/kwatchman
LD_FLAGS="-X main.Version=$(VERSION) -w -extldflags -static"


build: deps
	CGO_ENABLED=0 go build -ldflags $(LD_FLAGS) cmd/*.go
test:
	go test -v ./...
clean:
	go clean
deps:
	dep ensure -v

docker-image:
	docker build -f build/Dockerfile \
		--build-arg VERSION=$(VERSION) \
		--build-arg REPOSITORY=$(REPOSITORY) \
		-t snebel29/kwatchman:latest \
		-t snebel29/kwatchman:$(VERSION) .

publish-docker-image:
	docker push snebel29/kwatchman:$(VERSION)
	docker push snebel29/kwatchman:latest

.PHONY: build test clean docker-image publish-docker-image
