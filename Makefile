.PHONY: test tmp

compile-happy: tmp
	@go run cmd/quill/quill.go c -i internal/compiler/test/src/happy/test.adv -o tmp/happy-test.db

compile-ao: tmp
	@go run cmd/quill/quill.go c -i internal/compiler/test/src/ao/ao.adv -o tmp/test.db

tmp:
	@rm -Rf tmp
	@mkdir -p tmp

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
