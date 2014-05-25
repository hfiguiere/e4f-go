

all: e4f


.PHONY: e4f
e4f: e4f.go
	go build $<

