SHELL = /bin/bash

.PHONY: all build test check-diff

all: build test check-diff

build:
	go build

test: build
	@for dir in $(shell find test -type d -name schema); do \
		basedir=`dirname $$dir`; \
		output="$$basedir/generated.graphql"; \
		echo "Merging $$dir into $$output..."; \
		./gqlmerge $$dir $$output; \
	done

check-diff: build test
	@if git diff --exit-code --quiet -- '*.graphql'; then \
		echo "Ok"; \
	else \
		echo "Error: Differences found in generated.graphql files"; \
		exit -1; \
	fi
