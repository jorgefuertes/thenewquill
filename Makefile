.PHONY: test

run:
	@go run .

test:
	@go test ./...

test-v:
	go test -v ./...

lint:
	@echo "Linting..."
	@staticcheck ./...
	@golangci-lint run ./..

clean:
	@go clean -testcache
	@rm -Rf dist
	@rm -Rf tmp
