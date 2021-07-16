#!/usr/bin/env bash

echo Running BIT exporter service...
go run ./cmd/bitTestResultsExporter -config-file ./configs/prog_configs/configs.yml