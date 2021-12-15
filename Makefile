SHELL := /bin/bash

BUILD_CONTEXT := ./build
DOCKER_CONTEXT := ./docker

.PHONY: go-test
go-test:
	@echo "Run all project tests..."
	go test -p 1 ./...

.PHONY: bin
bin: clean
	@echo "Build project binaries..."
	GOOS=linux GOARCH=386 go build -v -o $(BUILD_CONTEXT)/spymaster_linux_386
	GOOS=darwin GOARCH=386 go build -v -o $(BUILD_CONTEXT)/spymaster_darwin_386
	GOOS=linux GOARCH=386 go build -v -o $(BUILD_CONTEXT)/spymaster_windows_386

.PHONY: clean
clean:
	rm -rf $(BUILD_CONTEXT)/

.PHONY: test
test: go-test

.PHONY: run
run:
	go run cmd/spymaster/main.go

.PHONY: get
get:
	@echo "Fetch project dependencies..."
	go mod tidy

.PHONY: docker-up
docker-up: docker-clean
	docker-compose -f $(DOCKER_CONTEXT)/docker-compose.yml pull
	docker-compose -f $(DOCKER_CONTEXT)/docker-compose.yml kill
	docker-compose -f $(DOCKER_CONTEXT)/docker-compose.yml up -d

.PHONY: docker-clean
docker-clean:
	@docker volume rm $(shell docker volume ls -qf dangling=true) 2>/dev/null ||:
	@docker rmi $(shell docker images -q -f dangling=true) 2>/dev/null ||: