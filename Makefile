NAME := donny
DOCKER_PREFIX = benmatselby

.DEFAULT_GOAL := explain
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
	### Targets
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

GITCOMMIT := $(shell git rev-parse --short HEAD)

.PHONY: clean
clean: ## Clean the local dependencies
	rm -fr vendor

.PHONY: install
install: ## Install the local dependencies
	go get ./...

.PHONY: vet
vet: ## Vet the code
	go vet ./...

.PHONY: lint
lint: ## Lint the code
	golint -set_exit_status $(shell go list ./...)

.PHONY: build
build: ## Build the application
	go build .

.PHONY: static
static: ## Build the application
	CGO_ENABLED=0 go build -ldflags "-extldflags -static -X github.com/benmatselby/donny/version.GITCOMMIT=$(GITCOMMIT)" -o $(NAME) .

.PHONY: test
test: ## Run the unit tests
	go test ./... -coverprofile=profile.out

.PHONY: test-cov
test-cov: test ## Run the unit tests with coverage
	go tool cover -html=profile.out

.PHONY: all ## Run everything
all: clean install lint vet build test

.PHONY: static-all ## Run everything
static-all: clean install vet static test
