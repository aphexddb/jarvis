ifeq ($(GOBIN),)
GOBIN := $(GOPATH)/bin
endif

DIST=./dist
SERVER_BIN=server
CLIENT_BIN=client
DOCKER_IP=`docker-machine ip`

clean:
	go clean

dependencies:
	go get -u github.com/Masterminds/glide
	glide install

build:
	@go version
	@mkdir -p $(DIST)
	GOOS=linux GOARCH=amd64 go build -o $(DIST)/$(SERVER_BIN) ./cmd/server/server.go
	GOOS=linux GOARCH=arm go build -o $(DIST)/$(CLIENT_BIN) ./cmd/client/client.go

docker:
	docker build -t jarvis:latest .
	@echo "Docker IP: $(DOCKER_IP)"

check:
	go vet `go list ./... | grep -v /vendor/`
	golint `go list ./... | grep -v /vendor/`

test: check
	go test `go list ./... | grep -v /vendor/`

ci: clean dependencies build test

deploy: build docker
	heroku container:push web

default: build

.PHONY: test