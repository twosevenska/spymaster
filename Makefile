BUILD_CONTEXT := ./build

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
	