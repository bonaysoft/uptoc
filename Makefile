.PHONY: default install build fmt test vet docker clean

BINARY=uptoc
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR := $(dir $(MKFILE_PATH))

TARGET_DIR=${MKFILE_DIR}build
TARGET_PATH=${TARGET_DIR}/${BINARY}

LDFLAGS="-s -w -X ${BINARY}/version.release=${RELEASE} -X ${BINARY}/version.commit=${COMMIT} -X ${BINARY}/version.repo=${GITREPO}"

# git info
COMMIT := git-$(shell git rev-parse --short HEAD)
GITREPO := $(shell git config --get remote.origin.url)
RELEASE := $(shell git describe --tags | awk -F '-' '{print $$1}')

default: install build

install:
	go mod download

build:
	GOOS=darwin GOARCH=amd64 go build -ldflags ${LDFLAGS} -o ${TARGET_PATH}-macos/${BINARY}
	GOOS=linux GOARCH=amd64 go build -ldflags ${LDFLAGS} -o ${TARGET_PATH}-linux-amd64/${BINARY}
	GOOS=windows GOARCH=amd64 go build -ldflags ${LDFLAGS} -o ${TARGET_PATH}-windows-amd64/${BINARY}

test:
	go test -coverprofile=coverage.txt -covermode=atomic ./...
    go tool cover --func=coverage.txt

covhtml:
	go tool cover -html=coverage.txt

pack:
	tar -C ${TARGET_DIR}/${BINARY}-macos -zvcf ${TARGET_PATH}-macos.tar.gz ${BINARY}
	tar -C ${TARGET_DIR}/${BINARY}-linux-amd64 -zvcf ${TARGET_PATH}-linux-amd64.tar.gz ${BINARY}
	tar -C ${TARGET_DIR}/${BINARY}-windows-amd64 -zvcf ${TARGET_PATH}-windows-x64.tar.gz ${BINARY}

clean:
	rm -rf ${TARGET_DIR}