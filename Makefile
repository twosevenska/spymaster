SHELL := /bin/bash

BUILD_CONTEXT := ./build
BUILD_TARGET := cmd/spymaster/main.go
DOCKER_CONTEXT := ./docker

.PHONY: go-test
go-test:
	@echo "Run all project tests..."
	go test -p 1 ./...

.PHONY: go-convey
go-convey:
	goconvey -workDir=./tests/

.PHONY: bin
bin: clean
	@echo "Build project binaries..."
	GOOS=linux GOARCH=386 go build -v -o $(BUILD_CONTEXT)/spymaster_linux_386 $(BUILD_TARGET)
	GOOS=darwin GOARCH=amd64 go build -v -o $(BUILD_CONTEXT)/spymaster_darwin_amd64 $(BUILD_TARGET)
	GOOS=windows GOARCH=386 go build -v -o $(BUILD_CONTEXT)/spymaster_windows_386 $(BUILD_TARGET)

.PHONY: clean
clean:
	rm -rf $(BUILD_CONTEXT)/*

.PHONY: run
run:
	go run $(BUILD_TARGET)

.PHONY: get
get:
	@echo "Fetch project dependencies..."
	go mod tidy

.PHONY: docker-pull
docker-pull:
	@echo "Updating containers"
	docker-compose -f $(DOCKER_CONTEXT)/docker-compose.yml pull

.PHONY: docker-kill
docker-kill:
	@echo "Killing old containers"
	docker-compose -f $(DOCKER_CONTEXT)/docker-compose.yml kill

.PHONY: docker-clean
docker-clean:
	@echo "Ensuring everything is scrubbed clean"
	@docker volume rm $(shell docker volume ls -qf dangling=true) 2>/dev/null ||:
	@docker rmi $(shell docker images -q -f dangling=true) 2>/dev/null ||:

.PHONY: docker-up
docker-up: docker-pull docker-kill docker-clean
	@echo "Starting containers"
	docker-compose -f $(DOCKER_CONTEXT)/docker-compose.yml up -d
