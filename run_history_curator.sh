#!/usr/bin/env bash

echo Running history curator service...
go run ./internal/bitHistoryCurator/cmd/history_curator_main.go -config-file ./configs/prog_configs/configs.yml &