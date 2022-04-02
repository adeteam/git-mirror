SHELL:=/bin/bash
ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
SRC_BASE=${ROOT_DIR}
SRC_MAIN_DIR=${SRC_BASE}/main
BINARY=git-mirror
BIN_DIR=${SRC_BASE}/build
ARCH=$(shell arch)

VERSION?=?
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-s -w"

# Build the project
all: ready clean linux

ready:
	BIN_DIR=${BIN_DIR}; \
	if [ ! -d "$${BIN_DIR}" ]; then \
	    mkdir $${BIN_DIR}; \
	fi

linux:
	cd ${SRC_MAIN_DIR}; \
	if [ "${ARCH}" == "x86_64" ]; then \
		GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BIN_DIR}/${BINARY}-linux-amd64 . ; \
	fi; \
	GOOS=linux GOARCH=386 go build ${LDFLAGS} -o ${BIN_DIR}/${BINARY}-linux-386 . ; \
	cd - >/dev/null

clean:
	-rm -f ${BIN_DIR}/${BINARY}-*

.PHONY: ready linux clean
