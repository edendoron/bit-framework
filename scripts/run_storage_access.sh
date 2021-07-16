#!/usr/bin/env bash

echo Running BIT storage access service...
go run ./cmd/bitStorageAccess -config-file ./configs/prog_configs/configs.yml
