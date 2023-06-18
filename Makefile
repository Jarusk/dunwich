SOURCES := ./...
BINARY := dunwich
BINDIR := .

FLAGS := CGO_ENABLED=0

.DEFAULT_GOAL := all-noclean

.PHONY: tooling
tooling:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: coverage
coverage: test
	@echo "Generating coverage report"
	go tool cover -html=coverage.out

.PHONY: check-coverage
check-coverage: test
	@echo "Get total coverage"
	$(eval COVERED=$(shell go tool cover -func coverage.out | grep total | awk '{print substr($$3, 1, length($$3)-1)}'))
	@if [ "80.0" != "$(word 1, $(sort 80.0 $(COVERED)))" ]; then \
		echo "Test coverage (${COVERED}%) is less than 80%";\
		exit 1;\
	fi

.PHONY: all
all: clean build fmt lint test

.PHONY: all-noclean
all-noclean: build fmt lint test

.PHONY: clean
clean:
	@echo "Removing $(BINARY)"
	@rm -f $(BINARY)
	@echo "Cleaning build and test cache"
	@go clean -cache -testcache

.PHONY: build
build: 
	@echo "Building ${BINARY}"
	${FLAGS} go build -buildvcs=true -o ${BINDIR} ${SOURCES}

.PHONY: fmt
fmt:
	@echo "Running gofmt"
	${FLAGS} gofmt -w .

.PHONY: lint
lint:
	@echo "Running lint"
	golangci-lint run -v ${SOURCES}

.PHONY: run
run: 
	@echo "Running ${BINARY}"
	${FLAGS} go run -buildvcs=true ${SOURCES}

.PHONY: test
test:
	@echo "Running tests"
	${FLAGS} CGO_ENABLED=1 go test -v -cover -coverprofile=coverage.out -race ${SOURCES}

.PHONY: --internal-update
--internal-update:
	@echo "Updating all dependencies"
	${FLAGS} go get -d -u -t ${SOURCES}

.PHONY: update
update: --internal-update vendor

.PHONY: vendor
vendor:
	@echo "Running mod tidy and mod vendor"
	${FLAGS} go mod tidy
	${FLAGS} go mod vendor

