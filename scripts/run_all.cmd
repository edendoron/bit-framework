
echo Starting BIT Framework...
start /b scripts\run_storage_access.cmd
timeout 1
start /b scripts\run_indexer.cmd
timeout 1
start /b scripts\run_exporter.cmd
timeout 1
start /b scripts\run_config.cmd
timeout 1
start /b scripts\run_handler.cmd
timeout 1
start /b scripts\run_query.cmd