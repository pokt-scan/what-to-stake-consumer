.PHONY: generate build build_docker start start_as_daemon

# Project root directory
PROJECT_ROOT := $(shell pwd)

# Project Version
VERSION := 1.0.0

# Define the default value for CONFIG_FILE
DEFAULT_CONFIG_FILE := config.json

# Check if CONFIG_FILE environment variable is set and not empty
# If not, use the default value
CONFIG_FILE := $(if $(CONFIG_FILE),$(CONFIG_FILE),$(DEFAULT_CONFIG_FILE))

# Define the project root directory
PROJECT_ROOT := $(shell pwd)

# Concatenate PROJECT_ROOT with CONFIG_FILE to form CONFIG_FILE_PATH
CONFIG_FILE_PATH := $(PROJECT_ROOT)/$(CONFIG_FILE)

# make generate or CONFIG_FILE=<configFileNameHere> make generate
generate:
	@echo "Generating POKTscan GraphQL schema types..."
	@PROJECT_ROOT=$(PROJECT_ROOT) VERSION=$(VERSION) CONFIG_FILE=$(CONFIG_FILE) go generate ./wtsc/...

# build wtsc project
build: generate
	@echo "Building the project in $(PROJECT_ROOT) with version $(VERSION)..."
	@mkdir -p bin
	@CONFIG_FILE=$(CONFIG_FILE) go build -o "$(PROJECT_ROOT)/bin/wtsc" cmd/wtsc/main.go

# build wtsc docker image
build_docker: generate
	@echo "Building docker image..."
	@CONFIG_FILE_PATH=$(CONFIG_FILE_PATH) docker compose build wtsc

# same as build_docker but force to omit cache on docker layers
build_docker_no_cache: generate
	@echo "Building docker image without cache..."
	@CONFIG_FILE_PATH=$(CONFIG_FILE_PATH) docker compose build --no-cache wtsc

# start wtsc on host
start: build
	@PROJECT_ROOT=$(PROJECT_ROOT) VERSION=$(VERSION) CONFIG_FILE=$(CONFIG_FILE) $(PROJECT_ROOT)/bin/wtsc

# start wtsc on docker using docker compose without release terminal
start_docker:
	@CONFIG_FILE_PATH=$(CONFIG_FILE_PATH) VERSION=$(VERSION) docker compose up wtsc

# start wtsc on docker using docker compose and release terminal
start_as_daemon:
	@CONFIG_FILE_PATH=$(CONFIG_FILE_PATH) VERSION=$(VERSION) docker compose up -d wtsc

# stop wtsc docker server without destroy it
stop_docker:
	@CONFIG_FILE_PATH=$(CONFIG_FILE_PATH) docker compose stop

# stop and destroy wtsc service from docker
down:
	@CONFIG_FILE_PATH=$(CONFIG_FILE_PATH) docker compose down -v