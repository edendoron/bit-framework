#!/usr/bin/env bash

echo Starting BIT Framework...
./scripts/run_storage_access.sh &
sleep 1
./scripts/run_indexer.sh &
sleep 1
./scripts/run_exporter.sh &
sleep 1
./scripts/run_config.sh &
sleep 1
./scripts/run_handler.sh &
sleep 1
./scripts/run_query.sh &