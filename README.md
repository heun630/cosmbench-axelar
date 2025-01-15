# Cosmbench-Axelar
## Example of Directory Structure
```
.
├── axelar-cosmbench_accounts
│   ├── node0
│   ├── node1
│   ├── node2
│   └── node3
├── axelar-cosmbench_encoded_txs
├── axelar-cosmbench_nodes
│   ├── node0
│   │   ├── config
│   │   │   └── gentx
│   │   ├── data
│   │   └── keyring-test
│   ├── node1
│   │   ├── config
│   │   │   └── gentx
│   │   ├── data
│   │   └── keyring-test
│   ├── node2
│   │   ├── config
│   │   │   └── gentx
│   │   ├── data
│   │   └── keyring-test
│   └── node3
│       ├── config
│       │   └── gentx
│       ├── data
│       └── keyring-test
├── axelar-cosmbench_signed_txs
├── axelar-cosmbench_unsigned_txs
├── bin
├── logs
├── nodes
│   ├── node0
│   │   ├── config
│   │   │   └── gentx
│   │   ├── data
│   │   │   ├── application.db
│   │   │   ├── blockstore.db
│   │   │   ├── cs.wal
│   │   │   ├── evidence.db
│   │   │   ├── snapshots
│   │   │   │   └── metadata.db
│   │   │   ├── state.db
│   │   │   └── tx_index.db
│   │   └── keyring-test
│   ├── node1
│   │   ├── config
│   │   │   └── gentx
│   │   ├── data
│   │   │   ├── application.db
│   │   │   ├── blockstore.db
│   │   │   ├── cs.wal
│   │   │   ├── evidence.db
│   │   │   ├── snapshots
│   │   │   │   └── metadata.db
│   │   │   ├── state.db
│   │   │   └── tx_index.db
│   │   └── keyring-test
│   ├── node2
│   │   ├── config
│   │   │   └── gentx
│   │   ├── data
│   │   │   ├── application.db
│   │   │   ├── blockstore.db
│   │   │   ├── cs.wal
│   │   │   ├── evidence.db
│   │   │   ├── snapshots
│   │   │   │   └── metadata.db
│   │   │   ├── state.db
│   │   │   └── tx_index.db
│   │   └── keyring-test
│   └── node3
│       ├── config
│       │   └── gentx
│       ├── data
│       │   ├── application.db
│       │   ├── blockstore.db
│       │   ├── cs.wal
│       │   ├── evidence.db
│       │   ├── snapshots
│       │   │   └── metadata.db
│       │   ├── state.db
│       │   └── tx_index.db
│       └── keyring-test
├── results
└── scripts
    ├── logs
    └── nodes
        ├── node0
        │   ├── config
        │   └── data
        │       ├── application.db
        │       ├── blockstore.db
        │       ├── snapshots
        │       │   └── metadata.db
        │       └── state.db
        ├── node1
        │   ├── config
        │   └── data
        │       ├── application.db
        │       ├── blockstore.db
        │       ├── snapshots
        │       │   └── metadata.db
        │       └── state.db
        ├── node2
        │   ├── config
        │   └── data
        │       ├── application.db
        │       ├── blockstore.db
        │       ├── snapshots
        │       │   └── metadata.db
        │       └── state.db
        └── node3
            ├── config
            └── data
                ├── application.db
                ├── blockstore.db
                ├── snapshots
                │   └── metadata.db
                └── state.db
```
---
## Initial Setup

Before proceeding with any node setup, ensure you clone the Axelar repository and build the project:

```bash
git clone https://github.com/axelarnetwork/axelar-core.git
cd axelar-core
make build
```

---

## Makefile Usage Guide

This repository contains a `Makefile` to automate the initialization, running, and management of nodes for a blockchain test environment. Below is a detailed guide for using the available Makefile targets.

---

### 1. **Initialization**
#### Command:
```bash
make init
```
#### Description:
Performs the full setup process, including:
1. Initializing nodes with default configurations.
2. Assigning validators to nodes.
3. Creating accounts for transactions.
4. Setting up the environment and configuring persistent peers.
5. Generating transaction files for testing.

#### Use Case:
Run this command to prepare the environment from scratch.

---

###  **Individual Initialization Step 1**

#### a. **Initialize Nodes**
```bash
make init-nodes
```
Initializes the nodes and creates default `genesis.json` files.

#### b. **Assign Validators**
```bash
make assign-validators
```
Assigns validator roles to nodes and updates `genesis.json` accordingly.

#### c. **Create Accounts**
```bash
make create-accounts
```
Creates accounts for transactions and adds them to `genesis.json`.

#### d. **Initialize Environment**
```bash
make initialize-env
```
Sets up configuration files (e.g., updating ports, enabling APIs, and configuring persistent peers).

#### e. **Generate Transactions**
```bash
make generate-transactions
```
Generates unsigned, signed, and encoded transactions for testing.

---

### 2. **Run Nodes**
#### Command:
```bash
make run
```
#### Description:
Starts all nodes concurrently using `scripts/92_run.sh` for each node.

---

### 3. **Send Transactions**
#### Command:
```bash
make send ARGS="TPS RunTime"
```
#### Example:
```bash
make send ARGS="200 10"
```
#### Description:
Sends transactions with the specified transactions per second (TPS) and runtime (in seconds).

---

### 4. **Restart Environment**
#### Command:
```bash
make restart
```
#### Description:
1. Re-initializes the environment by running `make initialize-env`.
2. Restarts all nodes using `make run`.

---

### 5. **Stop Nodes**
#### Command:
```bash
make stop
```
#### Description:
Stops all running nodes by killing their processes.

---

### 6. **Calculate Metrics**
#### Command:
```bash
make calculate
```
#### Description:
Runs the `metrics_calculator.go` script to analyze and calculate metrics from the blockchain logs and transactions.

---

### Example Workflow
1. **Set up the environment:**
   ```bash
   make init
   ```

2. **Start all nodes:**
   ```bash
   make run
   ```

3. **Send transactions:**
   ```bash
   make send ARGS="100 20"
   ```

4. **Calculate metrics:**
   ```bash
   make calculate
   ```

5. **Restart environment if needed:**
   ```bash
   make restart
   ```

6. **Stop all nodes:**
   ```bash
   make stop
   ```

---

### Notes
- Ensure all required scripts are present in the `scripts/` directory.
- Modify the Makefile variables if you need to customize paths or parameters.
- Use `make` commands in the root directory of the project.