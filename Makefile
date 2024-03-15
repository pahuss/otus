SHELL := /bin/bash

# build +
# rebuild +
# start +
# stop -
# remove -

.PHONY: default
default: help

.PHONY: remove ## Setup project from the beginning
remove: prepare remove

.PHONY: rebuild ## Setup project from the beginning
build: prepare remove build

.PHONY: prepare
prepare:
	{ \
    source .env; \
 	}

.PHONY: build
build:
	{ \
    source .env; \
    docker compose build \
    }

.PHONY: remove
remove:
	{ \
    source .env; \
    docker compose down \
    }


.PHONY: help ## Show this help
help:
	@echo "List of supported commands:"
	@grep -h ".PHONY" $(MAKEFILE_LIST) | fgrep -v fgrep | fgrep -v ".PHONY: default" | sed -e 's/^.PHONY: /\t/g' | sed -e "s/\(.*\)\s\+##\s*\(.*\)/\1 - \2/"
