all: e4f-go


.PHONY: e4f-go
e4f-go: e4f-tool.go
	go build .

check:
	go test

