#!/usr/bin/env bash

echo Starting BIT Framework...
./run_storage_access.sh &
sleep 1
./run_indexer.sh &
sleep 1
./run_exporter.sh &
sleep 1
./run_config.sh &
sleep 1
./run_handler.sh &
sleep 1
./run_query.sh &