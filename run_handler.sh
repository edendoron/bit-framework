#!/usr/bin/env bash

echo Running BIT handler service...
go run ./internal/bitHandler/cmd/bit_handler_main.go -config-file ./configs/prog_configs/configs.yml