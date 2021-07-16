#!/usr/bin/env bash

echo Running BIT query service...
go run ./internal/bitQuery/cmd/bit_query_main.go -config-file ./configs/prog_configs/configs.yml