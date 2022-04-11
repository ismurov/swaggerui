PROJECTNAME = Swagger UI

all: help

## lint: Run linters.
lint:
	@echo '>>> Run linters.'
	@golangci-lint version --format=short | awk '{printf "golangci-lint: %s\n", $$0}'
	@golangci-lint run --config .golangci.yaml
.PHONY: lint

## test: Run tests in short mode.
test:
	@echo '>>> Run tests in short mode.'
	@go test \
		-shuffle=on \
		-short \
		-count=1 \
		-timeout=5m \
		./...
.PHONY: test

## test-acc: Run all tests with accurate code coverage.
test-acc:
	@echo '>>> Run all tests with accurate code coverage.'
	@go test \
		-shuffle=on \
		-race \
		-count=1 \
		-timeout=10m \
		./... \
		-coverprofile=./coverage.out
.PHONY: test-acc

## test-coverage: Print out the code coverage information in the console (after test-acc).
test-coverage:
	@go tool cover -func=./coverage.out
.PHONY: test-coverage

## test-coverage-browser: Open a browser window and show the code coverage information (after test-acc).
test-coverage-browser:
	@go tool cover -html=./coverage.out
.PHONY: test-coverage-browser

help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
.PHONY: help
