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

.PHONY: test
test:
	go test ./... -coverprofile=build/profile.out

test-cov: test
	go tool cover -html=build/profile.out

.PHONY: all
all: clean install build test
