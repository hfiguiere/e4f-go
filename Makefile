
GOPATH=$(shell pwd)

all: e4f-tool


.PHONY: e4f-tool
e4f-tool: e4f-tool.go
	go build $<

check:
	go test

