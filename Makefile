.PHONY: test tmp dist lint format

dist: lint test
	@rm -rf dist
	@mkdir -p dist
	make compiler
	make runtime

compiler:
	@go build -o dist/qc cmd/compiler/qc.go

runtime:
	@go build -o dist/quill cmd/runtime/*.go

compile-ao:
	@mkdir -p tmp
	@go run cmd/compiler/qc.go c -i internal/compiler/test/src/ao/ao.adv -o tmp/ao.db

run-ao: compile-ao runtime
	@./dist/quill tmp/ao.db

test:
	@gotestsum -- -v -test.v ./...
	@echo "Trying to compile AO..."
	@make compile-ao > /dev/null

test-clean:
	@go clean -testcache

test-input: test-clean
	@go test -run "TestConsoleInput" ./internal/output/console/. -tags manual

format:
	@echo "Formatting..."
	@go mod tidy
	@go tool gofumpt -w .
	@go tool goimports -w .
	@go tool golines -m 120 -t 4 --ignore-generated --chain-split-dots -w .

lint:
	@echo "Linting Go..."
	@go tool gofumpt -l -w .
	@go tool staticcheck ./...
	@go tool golangci-lint cache clean
	@go tool golangci-lint run ./...
	@echo "Linting Markdown..."
	@npx --yes markdownlint-cli2 "**/*.md" "#tmp" "#dist" "#work"
	@echo "Checking vulnerabilities..."
	@go tool govulncheck ./...

clean: test-clean
	@rm -Rf dist
	@rm -Rf tmp
	@pushd docs/manual > /dev/null && \
		rm -f *.aux *.log *.out *.toc *.lol *.lot *.lof *.bbl *.blg *.idx *.ilg *.ind \
			*.nlo *.nls *.nlg *.spl *.synctex.gz *.fdb_latexmk *.fls *.listing *.pdf; \
		popd > /dev/null

doc:
	@pushd docs/manual > /dev/null && \
		pdflatex -shell-escape -interaction=nonstopmode \
			-file-line-error manual.tex; \
		popd > /dev/null

go-upgrade-deps:
	@go get -u ./...
