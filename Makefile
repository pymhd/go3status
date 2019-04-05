GOCMD=$(shell which go)
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=$(shell basename `pwd`)
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build

build: 
	$(info building "$(BINARY_NAME)" binary)
	@$(GOBUILD) -o $(BINARY_NAME) -v

test: 
	$(info Cheking dependencies)
	@$(GOTEST) -v ./... > /dev/null

clean: 
	$(info  go clean and binary file deleting)
	@$(GOCLEAN)
	@rm -f $(BINARY_NAME)
	@rm -f $(BINARY_UNIX)

install: 
	$(info copy "$(BINARY_NAME)" file to /usr/local/bin)
	@cp $(BINARY_NAME) /usr/local/bin


deps:
	$(info go get: github.com/pymhd/go-logging github.com/pymhd/go-simple-cache)
	@$(GOGET) github.com/pymhd/go-logging
	@$(GOGET) github.com/pymhd/go-simple-cache
