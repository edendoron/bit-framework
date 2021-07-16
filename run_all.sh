#!/usr/bin/env bash

echo Starting BIT Framework...
go run ./internal/bitStorageAccess/cmd/storage_access_main.go -config-file ./configs/prog_configs/configs.yml &
sleep 1
go run ./internal/bitIndexer/cmd/bit_indexer_main.go -config-file ./configs/prog_configs/configs.yml &
sleep 1
go run ./internal/bitTestResultsExporter/cmd/bit_exporter_main.go -config-file ./configs/prog_configs/configs.yml &
sleep 1
go run ./internal/bitConfig/cmd/bit_config_main.go -config-file ./configs/prog_configs/configs.yml &
sleep 1
go run ./internal/bitHandler/cmd/bit_handler_main.go -config-file ./configs/prog_configs/configs.yml &
sleep 1
go run ./internal/bitQuery/cmd/bit_query_main.go -config-file ./configs/prog_configs/configs.yml &