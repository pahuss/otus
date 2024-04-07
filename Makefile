SHELL := /bin/bash

# build +
# rebuild +
# start +
# stop -
# remove -

.PHONY: default
default: help

.PHONY: build
build: prepare _build

.PHONY: start
start: prepare _start

.PHONY: stop
stop:
	docker compose stop

.PHONY: rebuild
rebuild: prepare down build

.PHONY: prepare
prepare:
	{ \
	randpw(){ LC_CTYPE=C tr -dc "a-zA-Z0-9-_" < /dev/urandom < /dev/random | head -c 20; }; \
	if [ ! -f .env ]; then \
		echo "MYSQL_ROOT_PASSWORD=$$(randpw)" >> .env; \
		echo "DBUSER=social" >> .env; \
		echo "DBNAME=social" >> .env; \
		echo "DBPASS=$$(randpw)" >> .env; \
		echo "DBHOST=mysql_db:3306" >> .env; \
		echo "REDISHOST=redis:6379" >> .env; \
		set -a; \
		source .env; \
	else \
		set -a; \
		source .env; \
	fi; \
	}

.PHONY: _build
_build:
	docker compose build

.PHONY: _start
_start:
	docker compose up -d

.PHONY: down
down:
	docker compose down
.PHONY: help ## Show this help
help:
	@echo "List of supported commands:"
	@grep -h ".PHONY" $(MAKEFILE_LIST) | fgrep -v fgrep | fgrep -v ".PHONY: default" | sed -e 's/^.PHONY: /\t/g' | sed -e "s/\(.*\)\s\+##\s*\(.*\)/\1 - \2/"
