IMAGE_NAME=application
IMAGE_TAG=1.0.0
OUTPUT_DIR=./output
.PHONY: build

build:
	@go build -o $(OUTPUT_DIR)/$(IMAGE_NAME) .
