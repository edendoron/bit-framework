#!/usr/bin/env bash

echo Running BIT indexer service...
go run ./internal/bitIndexer/cmd/bit_indexer_main.go -config-file ./configs/prog_configs/configs.yml