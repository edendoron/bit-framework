#!/usr/bin/env bash

echo Running BIT exporter service...
go run ./internal/bitTestResultsExporter/cmd/bit_exporter_main.go -config-file ./configs/prog_configs/configs.yml