GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=http-echo

DOCKER_CMD=docker

build:
	$(GOBUILD) -o $(BINARY_NAME) -v
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
docker-build:
	$(DOCKER_CMD) build ./ -t thiagoquintoandar/http-echo:latest
docker-push:
	$(DOCKER_CMD) push thiagoquintoandar/http-echo:latest
