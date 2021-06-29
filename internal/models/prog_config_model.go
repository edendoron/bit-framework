package models

type ProgConfigs struct {
	BitConfigPort               string  `conf:"bit_config_port" help:"A message to print."`
	BitHandlerPort              string  `conf:"bit_handler_port" help:"A message to print."`
	BitHistoryCuratorPort       string  `conf:"bit_history_curator_port" help:"A message to print."`
	BitIndexerPort              string  `conf:"bit_indexer_port" help:"A message to print."`
	BitQueryPort                string  `conf:"bit_query_port" help:"A message to print."`
	BitStoragePort              string  `conf:"bit_storage_port" help:"A message to print."`
	BitExporterPort             string  `conf:"bit_exporter_port" help:"A message to print."`
	BitConfigPath               string  `conf:"bit_config_service_path" help:"A message to print."`
	BitExporterPath             string  `conf:"bit_exporter_service_path" help:"A message to print."`
	BitHandlerPath              string  `conf:"bit_handler_service_path" help:"A message to print."`
	BitIndexerPath              string  `conf:"bit_indexer_service_path" help:"A message to print."`
	BitQueryPath                string  `conf:"bit_query_service_path" help:"A message to print."`
	BitHistoryCuratorPath       string  `conf:"bit_history_curator_service_path" help:"A message to print."`
	BitStoragePath              string  `conf:"bit_storage_service_path" help:"A message to print."`
	StorageWriteURL             string  `conf:"storage_write_url" help:"A message to print."`
	StorageReadURL              string  `conf:"storage_read_url" help:"A message to print."`
	StorageDeleteURL            string  `conf:"storage_delete_url" help:"A message to print."`
	BitConfigFailuresPath       string  `conf:"config_failures_path" help:"A message to print."`
	BitConfigUserGroupPath      string  `conf:"config_user_group_path" help:"A message to print."`
	BitExporterDefaultBWSize    float32 `conf:"exporter_default_bw_size" help:"A message to print."`
	BitExporterDefaultBWUnits   string  `conf:"exporter_default_bw_units" help:"A message to print."`
	BitExporterPostToIndexerUrl string  `conf:"exporter_post_to_indexer_url" help:"A message to print."`
	BitHandlerTriggerPeriod     float64 `conf:"handler_trigger_period" help:"A message to print."`
	BitHandlerTriggerType       string  `conf:"handler_trigger_type" help:"A message to print."`
	BitHistoryCuratorAgedDate   string  `conf:"bit_history_curator_aged_date" help:"A message to print."`
}
