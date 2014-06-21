
GOPATH=$(shell pwd)

all: e4f-tool


.PHONY: e4f-tool
e4f-tool: e4f-tool.go
	GOPATH=$(GOPATH) go build $<

check:
	GOPATH=$(GOPATH) go test

