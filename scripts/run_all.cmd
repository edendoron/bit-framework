
echo Starting BIT Framework...
start /b scripts\run_storage_access.cmd
timeout 2
start /b scripts\run_indexer.cmd
timeout 2
start /b scripts\run_exporter.cmd
timeout 2
start /b scripts\run_config.cmd
timeout 2
start /b scripts\run_handler.cmd
timeout 2
start /b scripts\run_query.cmd
timeout 2
scripts\run_client.cmd