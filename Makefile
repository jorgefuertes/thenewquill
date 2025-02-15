.PHONY: test

compile:
	mkdir -p tmp
	@go run cmd/quill/quill.go -i internal/compiler/test/adv_files/happy/test.adv -o tmp/test.db

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
