SHELL = /bin/bash

.PHONY: all build test

all: build test

build:
	go build

test: build
	@for dir in $(shell find test -type d -name schema); do \
		basedir=`dirname $$dir`; \
		output="$$basedir/generated.graphql"; \
		echo "Merging $$dir into $$output..."; \
		./gqlmerge $$dir $$output; \
	done
