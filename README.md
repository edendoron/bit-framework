# BIT framework

Note to Eden: when you run each service you should add a flag of -config-file followed by a path in order to load configurations from file.

Example: to run bit-config from the main directory, run:

> go run ./internal/bitConfig/cmd/bit_config_main.go -config-file ./configs/prog_configs/configs.yml

You can also add "-config-file ./configs/prog_configs/configs.yml" as arguments in Clion (under edit configurations).
