#!/usr/bin/env bash

echo Running BIT config service...
go run ./internal/bitConfig/cmd/bit_config_main.go -config-file ./configs/prog_configs/configs.yml