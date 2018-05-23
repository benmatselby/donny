NAME := donny
DOCKER_PREFIX = benmatselby

.PHONY: explain
explain:
	### Welcome
	#  _______   ______   .__   __. .__   __. ____    ____
	# |       \ /  __  \  |  \ |  | |  \ |  | \   \  /   /
	# |  .--.  |  |  |  | |   \|  | |   \|  |  \   \/   /
	# |  |  |  |  |  |  | |  .    | |  .    |   \_    _/
	# |  '--'  |   --'  | |  |\   | |  |\   |     |  |
	# |_______/ \______/  |__| \__| |__| \__|     |__|
	#
	#
	### Installation
	#
	# $$ make all
	#

GITCOMMIT := $(shell git rev-parse --short HEAD)

.PHONY: clean
clean:
	rm -fr vendor

.PHONY: install
install:
	dep ensure

.PHONY: vet
vet:
	go vet -v ./...

.PHONY: build
build:
	go build .

.PHONY: static
static:
	CGO_ENABLED=0 go build -ldflags "-extldflags -static -X github.com/benmatselby/donny/version.GITCOMMIT=$(GITCOMMIT)" -o $(NAME) .

.PHONY: test
test:
	go test ./... -coverprofile=profile.out

.PHONY: test-cov
test-cov: test
	go tool cover -html=profile.out

.PHONY: all
all: clean install vet build test

.PHONY: static-all
static-all: clean install vet static test

.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_PREFIX)/$(NAME) .

.PHONY: docker-push
docker-push:
	docker push $(DOCKER_PREFIX)/$(NAME)

.PHONY: docker-all
docker-all: docker-build docker-push
