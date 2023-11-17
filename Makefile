SOURCES := ./...
BINARY := dunwich
BINDIR := out
COVERAGE_REPORT := coverage.out

FLAGS := CGO_ENABLED=0

.DEFAULT_GOAL := all-noclean

.PHONY: tooling
tooling:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: coverage
coverage: test-coverage
	@echo "Generating coverage report"
	go tool cover -html=$(COVERAGE_REPORT)

.PHONY: check-coverage
check-coverage: test-coverage
	@echo "Get total coverage"
	$(eval COVERED=$(shell go tool cover -func $(COVERAGE_REPORT) | grep total | awk '{print substr($$3, 1, length($$3)-1)}'))
	@if [ "80.0" != "$(word 1, $(sort 80.0 $(COVERED)))" ]; then \
		echo "Test coverage (${COVERED}%) is less than 80%";\
		exit 1;\
	fi

.PHONY: doc
doc:
	@echo "Start godoc"
	godoc -http=:6060

.PHONY: all
all: clean fmt lint test build

.PHONY: all-noclean
all-noclean: fmt lint test build

.PHONY: clean
clean:
	@echo "Removing $(BINDIR) and coverage report"
	@rm -rf $(BINDIR) $(COVERAGE_REPORT)
	@echo "Cleaning build and test cache"
	@go clean -cache -testcache

.PHONY: build
build: 
	@echo "Building binaries"
	mkdir -p ${BINDIR}
	${FLAGS} go build -buildvcs=true -o ${BINDIR} ${SOURCES}

.PHONY: generate
generate:
	@echo "Generating source files"
	${FLAGS} go generate ${SOURCES}

.PHONY: fmt
fmt:
	@echo "Running gofmt"
	${FLAGS} gofmt -w .

.PHONY: install
install:
	@echo "Installing ${BINARY} to $(shell go env GOPATH)/bin"
	${FLAGS} go install  -buildvcs=true ${SOURCES}

.PHONY: lint
lint:
	@echo "Running lint"
	golangci-lint run -v ${SOURCES}

.PHONY: run
run: 
	@echo "Running ${BINARY}"
	${FLAGS} go run -buildvcs=true ${SOURCES}

.PHONY: bench
bench:
	@echo "Running benchs"
	${FLAGS} CGO_ENABLED=1 go test -v -bench=. -benchmem -race ${SOURCES}

.PHONY: test
test:
	@echo "Running tests"
	${FLAGS} CGO_ENABLED=1 go test -v -race ${SOURCES}

.PHONY: test-coverage
test-coverage:
	@echo "Running tests"
	${FLAGS} CGO_ENABLED=1 go test -v -cover -coverprofile=$(COVERAGE_REPORT) -race ${SOURCES}

.PHONY: tidy
tidy:
	@echo "Running mod tidy"
	${FLAGS} go mod tidy

.PHONY: --internal-update
--internal-update:
	@echo "Updating all dependencies"
	${FLAGS} go get -d -u -t ${SOURCES}

.PHONY: update
update: --internal-update vendor tidy

.PHONY: vendor
vendor:
	@echo "Running mod vendor"
	${FLAGS} go mod vendor
