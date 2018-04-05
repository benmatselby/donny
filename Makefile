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
	mkdir build;

.PHONY: install
	dep ensure

.PHONY: build
build:
	go build .

.PHONY: test
test:
	@for d in $(shell go list ./...); do \
		go test -race -coverprofile=build/profile.out -covermode=atomic "$$d"; \
		if [ -f build/profile.out ]; then \
			cat build/profile.out >> build/coverage.out; \
			rm build/profile.out; \
		fi; \
	done;

test-cov:
	go tool cover -html=build/coverage.out
