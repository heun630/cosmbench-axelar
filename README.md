# Cosmbench-Axelar

## Multi-Node Setup Workflow

If you want to run multiple Axelar nodes, follow these steps in order:

1. **Setup**: Execute the setup script:
   ```bash
   ./2-multi_app_chain_setup.sh
   ```

2. **Connect Peers**: Establish peer connections:
   ```bash
   ./5-peer_connect.sh
   ```

3. **Monitor**: Monitor the nodes:
   ```bash
   ./4-monit.sh
   ```

---

## Single Node Setup

To set up a single Axelar node, execute:

```bash
./1-setup.sh
```

### Input Parameters for Single Node

- **Data Directory**: Specify the path for the data directory (e.g., `/data/axelar`).
- **Node Name**: Specify the name of the node (e.g., `node1`).
- **Chain ID**: Specify the chain ID (e.g., `my-private-chain`).

The script will prompt for these inputs during execution.

---

## Peer Connection for Multi-Node Setup

After setting up multiple nodes, connect the peers by executing:

```bash
./5-peer_connect.sh
```

This script configures and establishes persistent peer connections among all the nodes in the setup.

---

## Monitoring Nodes

To monitor the current running nodes, execute:

```bash
./4-monit.sh
```

This script allows you to select the number of nodes to monitor and displays their live logs and statuses.

---

## Reset and Cleanup

To stop and reset the current Axelar setup, execute:

```bash
./3-reset.sh
```

### What This Script Does:
- Stops all running Axelar nodes.
- Cleans up the data directory (`e.g., /data/axelar`) to reset the environment.

This script ensures a clean state, allowing you to start fresh.