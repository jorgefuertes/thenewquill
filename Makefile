.PHONY: test tmp dist

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
	@go test ./...

test-v:
	go test -v ./...

test-clean:
	@go clean -testcache

test-input: test-clean
	@go test -run "TestConsoleInput" ./internal/output/console/. -tags manual

lint:
	@echo "Linting..."
	@go tool gofumpt -l -d .
	@go tool staticcheck ./...
	@go tool golangci-lint run ./...

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