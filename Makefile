# Makefile for building and running the Go project

# Variables
BINARY_NAME=keryx

# Default target
.PHONY: all
all: build

# Build the Go project
.PHONY: build
build:
	go build -o $(BINARY_NAME)

# Run the project
.PHONY: run
run: build
	@source .env && ./$(BINARY_NAME)

# Clean the build
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
