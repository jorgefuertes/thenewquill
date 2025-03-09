.PHONY: test tmp dist

dist: lint test
	@rm -rf dist
	@mkdir -p dist
	make compiler
	make runtime

compiler:
	@go build -o dist/qc cmd/compiler/qc.go

runtime:
	@go build -o dist/quill cmd/runtime/quill.go

compile-ao:
	@mkdir -p tmp
	@go run cmd/compiler/qc.go c -i internal/compiler/test/src/ao/ao.adv -o tmp/ao.db

test:
	@go test ./...

test-v:
	go test -v ./...

lint:
	@echo "Linting..."
	@staticcheck ./...
	@golangci-lint run ./...

clean:
	@go clean -testcache
	@rm -Rf dist
	@rm -Rf tmp
