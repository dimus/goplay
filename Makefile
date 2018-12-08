GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean


all: install

build: peg
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD)

install: peg
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOINSTALL)

peg:
	pigeon -o names.go names.peg
