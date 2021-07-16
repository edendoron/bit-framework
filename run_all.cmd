
echo Starting BIT Framework...
start /b go run ./internal/bitStorageAccess/cmd/storage_access_main.go -config-file ./configs/prog_configs/configs.yml
timeout 1
start /b go run ./internal/bitIndexer/cmd/bit_indexer_main.go -config-file ./configs/prog_configs/configs.yml
timeout 1
start /b go run ./internal/bitTestResultsExporter/cmd/bit_exporter_main.go -config-file ./configs/prog_configs/configs.yml
timeout 1
start /b go run ./internal/bitConfig/cmd/bit_config_main.go -config-file ./configs/prog_configs/configs.yml
timeout 1
start /b go run ./internal/bitHandler/cmd/bit_handler_main.go -config-file ./configs/prog_configs/configs.yml
timeout 1
start /b go run ./internal/bitQuery/cmd/bit_query_main.go -config-file ./configs/prog_configs/configs.yml