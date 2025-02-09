.DEFAULT_GOAL = verify

BIN=bin

export GOBIN=$(CURDIR)/$(BIN)


$(BIN)/golangci-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.4

$(BIN)/gotestsum:
	go install gotest.tools/gotestsum@v1.11.0

$(BIN)/goimports:
	go install golang.org/x/tools/cmd/goimports@latest

.PHONY: install
install: $(BIN)/golangci-lint  $(BIN)/goimports  $(BIN)/gotestsum

.PHONY: lint
lint:
	$(BIN)/golangci-lint run --config=.golangci.yml ./...

.PHONY: fix
fix:
	gofmt -s -w .

	$(BIN)/goimports -l -w .

	$(BIN)/golangci-lint run --config=.golangci.yml ./... --fix

.PHONY: test
test:
	$(BIN)/gotestsum ./... -race -v -coverprofile=cover.out -covermode=atomic

.PHONY: build
build:
	go build -v -o main ./cmd/main.go

.PHONY: docker-dev
docker-dev:
	docker compose -f docker-compose.dev.yml up --build

.PHONY: testcover
testcover:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...
