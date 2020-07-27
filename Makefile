.PHONY: help

help: ## help command for available tasks
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help


build: ## build the container
	docker build -t stella .

build-nc: ## build the container w/o a cache
	docker build --no-cache -t stella .

run: ## run the container with default parameters
	docker run --net=host -v $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))/src/config/:/config stella

up: build run ## build the container and boot

image:
	docker tag stella docker.pkg.github.com/adityaxdiwakar/stella/stella:${TRAVIS_TAG}
	docker tag stella docker.pkg.github.com/adityaxdiwakar/stella/stella:latest
	
push-image:
	docker push docker.pkg.github.com/adityaxdiwakar/stella/stella:${TRAVIS_TAG}
	docker push docker.pkg.github.com/adityaxdiwakar/stella/stella:latest
