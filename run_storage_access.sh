#!/usr/bin/env bash

echo Running BIT storage access service...
go run ./internal/bitStorageAccess/cmd/storage_access_main.go -config-file ./configs/prog_configs/configs.yml
