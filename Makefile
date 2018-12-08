GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean


all: build

build: peg
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD)


peg:
	echo "package main\n" > names.go
	pigeon names.peg >> names.go