SHELL := /bin/bash
SRC ?= $(shell pwd)
TARGET_DIR ?= $(shell pwd)/target
NAME := tman

GO_BUILD_CONTEXT = $(SRC)/cmd/tman

build-binary:
	@mkdir -p $(TARGET_DIR)
	@echo ">> building binary $(TARGET_DIR)/$(NAME)"
	go build -C $(GO_BUILD_CONTEXT) -o $(TARGET_DIR)/$(NAME)

