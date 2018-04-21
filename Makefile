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

.PHONY: clean
clean:
	rm -fr vendor

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
	go test ./... -coverprofile=profile.out

.PHONY: test-cov
test-cov: test
	go tool cover -html=profile.out

.PHONY: all
all: clean install build test

.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_PREFIX)/$(NAME) .

.PHONY: docker-push
docker-push:
	docker push $(DOCKER_PREFIX)/$(NAME)

.PHONY: docker-all
docker-all: docker-build docker-push
