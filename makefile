# These will be provided to the target
VERSION := 0.1.0
BUILD := `git rev-parse HEAD`

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-s -w -X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=ece-support-diagnostics
BINARY_UNIX=dist/linux/$(BINARY_NAME)

all: test build
build: 
		$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v
test: 
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)
# deps:
#         $(GOGET) github.com/markbates/goth
#         $(GOGET) github.com/markbates/pop


# Cross compilation
build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_UNIX) -v
# docker-build:
#         docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v