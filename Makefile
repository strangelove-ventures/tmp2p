VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT  := $(shell git log -1 --format='%H')
DIRTY := $(shell git status --porcelain | wc -l | xargs)

GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin

ldflags = -X github.com/strangelove-ventures/tmp2p/cmd.Version=$(VERSION) \
					-X github.com/strangelove-ventures/tmp2p/cmd.Commit=$(COMMIT) \
					-X github.com/strangelove-ventures/tmp2p/cmd.Dirty=$(DIRTY)

ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

all: install build-static-amd64 build-static-arm64

build: go.sum
	@echo "building tmp2p binary..."
	@go build -mod=readonly -o build/tmp2p -ldflags '$(ldflags)' .

install: go.sum
	@echo "installing tmp2p binary..."
	@go build -mod=readonly -o $(GOBIN)/tmp2p -ldflags '$(ldflags)' .

build-static: build-static-amd64 build-static-arm64

build-static-amd64:
	@echo "building tmp2p amd64 static binary..."
	@GOOS=linux GOARCH=amd64 go build -mod=readonly -o build/tmp2p-amd64 -a -tags netgo -ldflags '$(ldflags) -extldflags "-static"' .

build-static-arm64:
	@echo "building tmp2p arm64 static binary..."
	@GOOS=linux GOARCH=arm64 go build -mod=readonly -o build/tmp2p-arm64 -a -tags netgo -ldflags '$(ldflags) -extldflags "-static"' .

.PHONY: all build build-static-amd64 build-static-arm64