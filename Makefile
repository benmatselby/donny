NAME := donny

.PHONY: explain
explain:
	### Welcome
	#
	# Makefile for donny
	#
	### Installation
	#
	# $$ make clean install
	#

.PHONY: clean
clean:
	rm -fr build;
	rm -fr vendor
	mkdir build;

.PHONY: install
install:
	dep ensure

.PHONY: build
build:
	go build .

.PHONY: static
static:
	go build -ldflags "-linkmode external -extldflags -static" -o $(NAME) .

.PHONY: test
test:
	go test ./... -coverprofile=build/profile.out

.PHONY: test-cov
test-cov: test
	go tool cover -html=build/profile.out

.PHONY: all
all: clean install build test

.PHONY: docker-build
docker-build:
	docker build -t benmatselby/donny .
