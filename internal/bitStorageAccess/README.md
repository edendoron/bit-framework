# BIT Storage Access

'BIT Storage Access' microservice handles all the read, write & delete requests concerning the framework's storage.

Currently, the storage is implemented on top of the local file system of the machine running this service, this means that the service communicates with the local file system through the Go standard 'os' package.

Alternatively, the service can be implemented with databases if needed.

All the data in the storage is encoded using protobuf serialization for memory efficiency.

###File Structure

- **`./bit_status`**: Contains all the system's statuses reported by the 'BIT Handler' service, ordered by their creation timestamp. Further details about the BIT status can be found in [BIT Handler](https://github.com/edendoron/bit-framework/tree/master/internal/bitHandler)
- **`./test_reports`**: Contains all test reports which are reported by the different bit-clients and passed through the 'BIT Test Results Exporter' & 'BIT Indexer' services. The reports are ordered by their timestamps, where multiple reports sharing a timestamp reside in the same `test_results.txt` file.
- **`./config/`**:
    - **`./filtering_rules/`**: Contains the failures definitions and rules as explained in [BIT Config](https://github.com/edendoron/bit-framework/tree/master/internal/bitConfig)
    - **`./perm_filtering_rules/`**: Same files as in `./filtering_rules/`.
    - **`./user_groups/`**: A single file named `user_groups.txt` containing the different user groups in the system, which are pulled from the files in `./filtering_rules/` when 'BIT Config' sends them to the storage.
    - **`./user_groups_masks/`**: Explained in [BIT Config](https://github.com/edendoron/bit-framework/tree/master/internal/bitConfig)