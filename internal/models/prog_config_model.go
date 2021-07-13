package models

type ProgConfigs struct {
	BitConfigPort                  string  `conf:"bit_config_port" help:"The port of bit_config service."`
	BitHandlerPort                 string  `conf:"bit_handler_port" help:"The port of bit_handler service."`
	BitHistoryCuratorPort          string  `conf:"bit_history_curator_port" help:"The port of bit_history_curator service."`
	BitIndexerPort                 string  `conf:"bit_indexer_port" help:"The port of bit_indexer service."`
	BitQueryPort                   string  `conf:"bit_query_port" help:"The port of bit_query service."`
	BitStoragePort                 string  `conf:"bit_storage_port" help:"The port of bit_storage_access service."`
	BitExporterPort                string  `conf:"bit_exporter_port" help:"The port of bit_test_results_exporter service."`
	BitConfigPath                  string  `conf:"bit_config_service_path" help:"The path of bit_config service."`
	BitExporterPath                string  `conf:"bit_exporter_service_path" help:"The path of bit_test_results_exporter service."`
	BitHandlerPath                 string  `conf:"bit_handler_service_path" help:"The path of bit_handler service."`
	BitIndexerPath                 string  `conf:"bit_indexer_service_path" help:"The path of bit_indexer service."`
	BitQueryPath                   string  `conf:"bit_query_service_path" help:"The path of bit_query service."`
	BitHistoryCuratorPath          string  `conf:"bit_history_curator_service_path" help:"The path of bit_history_curator service."`
	BitStoragePath                 string  `conf:"bit_storage_service_path" help:"The path of bit_storage_access service."`
	StorageWriteURL                string  `conf:"storage_write_url" help:"URL used in services in order to create a write request to storage (POST, PUT)"`
	StorageReadURL                 string  `conf:"storage_read_url" help:"URL used in services in order to create a read request to storage (GET)"`
	StorageDeleteURL               string  `conf:"storage_delete_url" help:"URL used in services in order to create a delete request to storage (DELETE)"`
	BitConfigFailuresPath          string  `conf:"config_failures_path" help:"The path in which bit_config service expects to find configuration failures in order to post to storage"`
	BitConfigUserGroupPath         string  `conf:"config_user_group_path" help:"The path in which bit_config service expects to find configuration user_groups_filtering in order to post to storage"`
	BitExporterDefaultBWSize       float32 `conf:"exporter_default_bw_size" help:"bit_exporter default bandwidth size"`
	BitExporterDefaultBWUnits      string  `conf:"exporter_default_bw_units" help:"bit_exporter default bandwidth units (KiB, Mib, MB...)"`
	BitExporterPostToIndexerUrl    string  `conf:"exporter_post_to_indexer_url" help:"URL used in bit_exporter in order to create a write request to storage through indexer service (POST, PUT)."`
	BitHandlerTriggerPeriod        float64 `conf:"handler_trigger_period" help:"Interval size of the bit_handler trigger in seconds. BitStatus is written to storage every triggerPeriod seconds."`
	BitHandlerTriggerType          string  `conf:"handler_trigger_type" help:"Type of bit_handler trigger (CBIT, PBIT, IBIT...)."`
	BitHistoryCuratorAgedDataLimit string  `conf:"bit_history_curator_aged_data_limit" help:"Limit of time for bit_history_curator. every report or bitStatus reported prior to that duration will be deleted from storage. Valid time units are 'ns', 'us' (or 'µs'), 'ms', 's', 'm', 'h'."`
	SSHKeyPath					           string	`conf:"ssh_key" help:"The SSH key used for HTTPS."`
	SSHCertPath					           string	`conf:"ssh_cert" help:"The SSH certificate used for HTTPS."`
	SSHCsrPath					           string	`conf:"ssh_csr" help:"The SSH csr file used for HTTPS."`
}
