dir = $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
.PHONY: help

help: ## help command for available tasks
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help


build: ## build the container
	docker build -t stella .

build-nc: ## build the container w/o a cache
	docker build --no-cache -t stella .

run: ## run the container with default parameters
	docker run --net=host -v $(dir)/src/config/:/config -v $(dir)/src/assets/:/assets stella

up: build run ## build the container and boot

image:
	docker tag stella docker.pkg.github.com/adityaxdiwakar/stella/stella:latest
	
push-image:
	docker push docker.pkg.github.com/adityaxdiwakar/stella/stella:latest
