# Makefile for building and installing the wday CLI tool

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=wday

# Default target executed when no arguments are given to make.
all: build

# Build the project
build:
	$(GOBUILD) -o $(BINARY_NAME) main.go

# Install the binary
install: build
	cp $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

# Clean the build
clean:
	rm -f $(BINARY_NAME)

.PHONY: all build clean
