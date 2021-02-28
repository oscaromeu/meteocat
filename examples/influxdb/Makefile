VERSION=`git rev-parse HEAD`
BUILD=`date +%FT%T%z`
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"
DOCKER_IMAGE=meteocat-testing
DOCKER_IMAGE_TAG=meteocat-testing
PLATFORM=local

.PHONY: help
help: ## - Show help message
	@printf "\033[32m\xE2\x9c\x93 usage: make [target]\n\n\033[0m"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: docker-pull
docker-pull:	## - docker pull latest images
	@printf "\033[32m\xE2\x9c\x93 docker pull latest images\n\033[0m"
	@docker pull golang:alpine

.PHONY: docker-pull-distroless
docker-pull-distroless:	## - docker pull latest images
	@printf "\033[32m\xE2\x9c\x93 docker pull latest images\n\033[0m"
	@docker pull golang:buster
	@docker pull gcr.io/distroless/base
	@docker pull gcr.io/distroless/static

.PHONY: build-distroless
build-distroless:docker-pull-distroless	## - Build the smallest and secured golang docker image based on distroless
	@printf "\033[32m\xE2\x9c\x93 Build the smallest and secured golang docker image based on distroless\n\033[0m"
	$(eval BUILDER_IMAGE=$(shell docker inspect --format='{{index .RepoDigests 0}}' golang:buster))
	$(eval DISTROLESS_IMAGE=$(shell docker inspect --format='{{index .RepoDigests 0}}' gcr.io/distroless/base))
	@export DOCKER_CONTENT_TRUST=1
	@docker build -f Dockerfile --build-arg BUILDER_IMAGE=$(BUILDER_IMAGE) --build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) -t $(DOCKER_IMAGE_TAG)-distroless .

.PHONY: ls
ls: ## - List '$(DOCKER_IMAGE_TAG)' docker images
	@printf "\033[32m\xE2\x9c\x93 Look at the size dude !\n\033[0m"
	@docker image ls $(DOCKER_IMAGE_TAG)

.PHONY: run
run:	## - Run the smallest and secured golang docker image based on scratch
	@printf "\033[32m\xE2\x9c\x93 Run the smallest and secured golang docker image based on scratch\n\033[0m"
	@docker run $(DOCKER_IMAGE_TAG):latest

