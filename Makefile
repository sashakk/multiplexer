GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

build:
	$(GOBUILD) -o ./bin/server ./cmd/server

test:
	$(GOTEST) ./...

clean:
	rm -f ./bin/*


.PHONY: build-server test clean
