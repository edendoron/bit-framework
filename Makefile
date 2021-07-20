GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GORUN=$(GOCMD) run
LINTER=golangci-lint
CONFIGS=-config-file ./configs/prog_configs/configs.yml

#all: test build

build: build_bit_query build_bit_indexer build_bit_config build_bit_exporter build_bit_handler build_bit_history_curator build_bit_storage_access

test: test_bit_handler test_bit_query test_bit_exporter

test_bit_query:
	$(GOTEST) ./internal/bitQuery/tests

test_bit_exporter:
	$(GOTEST) ./internal/bitTestResultsExporter/tests

test_bit_handler:
	$(GOTEST) ./internal/bitHandler/tests

build_bit_query:
	$(GOBUILD) -o bin/ ./cmd/bitQuery

build_bit_indexer:
	$(GOBUILD) -o bin/ ./cmd/bitIndexer

build_bit_config:
	$(GOBUILD) -o bin/ ./cmd/bitConfig

build_bit_handler:
	$(GOBUILD) -o bin/ ./cmd/bitHandler

build_bit_history_curator:
	$(GOBUILD) -o bin/ ./cmd/bitHistoryCurator

build_bit_storage_access:
	$(GOBUILD) -o bin/ ./cmd/bitStorageAccess

build_bit_exporter:
	$(GOBUILD) -o bin/ ./cmd/bitTestResultsExporter

run_bit_query:
	$(GORUN) ./cmd/bitQuery $(CONFIGS)

run_bit_indexer:
	$(GORUN) ./cmd/bitIndexer $(CONFIGS)

run_bit_config:
	$(GORUN) ./cmd/bitConfig $(CONFIGS)

run_bit_handler:
	$(GORUN) ./cmd/bitHandler $(CONFIGS)

run_bit_history_curator:
	$(GORUN) ./cmd/bitHistoryCurator $(CONFIGS)

run_bit_storage_access:
	$(GORUN) ./cmd/bitStorageAccess $(CONFIGS)

run_bit_exporter:
	$(GORUN) ./cmd/bitTestResultsExporter $(CONFIGS)

lint:
	$(LINTER) run
