
echo Starting BIT Framework...
start /b run_storage_access.cmd
timeout 1
start /b run_indexer.cmd
timeout 1
start /b run_exporter.cmd
timeout 1
start /b run_config.cmd
timeout 1
start /b run_handler.cmd
timeout 1
start /b run_query.cmd