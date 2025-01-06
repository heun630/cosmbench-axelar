# Cosmbench-Axelar
## Example of Directory Structure
```
data
└── axelar
    ├── node0
    │   ├── config
    │   │   └── gentx
    │   ├── data
    │   │   ├── application.db
    │   │   ├── blockstore.db
    │   │   ├── cs.wal
    │   │   ├── evidence.db
    │   │   ├── snapshots
    │   │   │   └── metadata.db
    │   │   ├── state.db
    │   │   └── tx_index.db
    │   ├── keyring-file
    │   └── keyring-test
    └── node1
        ├── config
        ├── data
        │   ├── application.db
        │   ├── blockstore.db
        │   ├── cs.wal
        │   ├── evidence.db
        │   ├── snapshots
        │   │   └── metadata.db
        │   ├── state.db
        │   └── tx_index.db
        └── keyring-file
```
---
## Initial Setup

Before proceeding with any node setup, ensure you clone the Axelar repository and build the project:

```bash
git clone https://github.com/axelarnetwork/axelar-core.git
cd axelar-core
make build
```