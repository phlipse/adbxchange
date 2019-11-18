VERSION = $(shell git describe --tags)
BUILD = $(shell git rev-parse --short HEAD)
PROJECTNAME = adbXchange

GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
BINARY_NAME = $(PROJECTNAME)
BINARY_NAME_WIN = $(BINARY_NAME).exe

LDFLAGS=-ldflags "-s -w -X=main.VERSION=$(VERSION) -X=main.BUILD=$(BUILD)"
LDFLAGS_WIN=-ldflags "-s -w -X=main.VERSION=$(VERSION) -X=main.BUILD=$(BUILD) -H=windowsgui"

all: test build
build: 
		$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v
		$(GOBUILD) $(LDFLAGS_WIN) -o $(BINARY_NAME_WIN) -v
test:
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_NAME).exe
run:
		$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)
run_win:
		$(GOBUILD) $(LDFLAGS_WIN) -o $(BINARY_NAME_WIN) -v ./...
		./$(BINARY_NAME_WIN)
