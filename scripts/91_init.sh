#!/bin/bash

# Load environment variables relative to the script's location
SCRIPT_DIR=$(dirname "$0")
source "$SCRIPT_DIR/env.sh"
source "$SCRIPT_DIR/run_env.sh"

if [ -d "$TESTDIR" ]; then
    echo "Removing $TESTDIR directory..."
    rm -rf "$TESTDIR"
fi

cp -rf "$NODE_ROOT_DIR" "$TESTDIR"

# Copy genesis.json for all nodes
for ((i = 1; i < $NODE_COUNT; i++)); do
    CURRENT_DATA_DIR="$TESTDIR/node$i"
    cp -f "$TESTDIR/node0/config/genesis.json" "$CURRENT_DATA_DIR/config/genesis.json"
done

for ((i = 0; i < $NODE_COUNT; i++)); do
    INDEX=$i
    CURRENT_DATA_DIR="$TESTDIR/node$i"

    # Replace "stake" with "uaxl" in genesis.json
    sed -i 's/"stake"/"uaxl"/g' "$CURRENT_DATA_DIR/config/genesis.json"

    # Update config.toml with custom ports
    sed -i "s#proxy_app = \"tcp://127.0.0.1:26658\"#proxy_app = \"tcp://${PRIVATE_HOSTS[$INDEX]}:${PROXYAPP_PORTS[$INDEX]}\"#g" "$CURRENT_DATA_DIR/config/config.toml"
    sed -i "s#laddr = \"tcp://127.0.0.1:26657\"#laddr = \"tcp://${PRIVATE_HOSTS[$INDEX]}:${RPC_PORTS[$INDEX]}\"#g" "$CURRENT_DATA_DIR/config/config.toml"
    sed -i "s#laddr = \"tcp://0.0.0.0:26656\"#laddr = \"tcp://${PRIVATE_HOSTS[$INDEX]}:${P2P_PORTS[$INDEX]}\"#g" "$CURRENT_DATA_DIR/config/config.toml"
    sed -i "s#pprof_laddr = \"localhost:6060\"#pprof_laddr = \"${PRIVATE_HOSTS[$INDEX]}:${PPROF_PORTS[$INDEX]}\"#g" "$CURRENT_DATA_DIR/config/config.toml"

    # Enable duplicate IPs and increase mempool size
    sed -i 's/allow_duplicate_ip = false/allow_duplicate_ip = true/g' "$CURRENT_DATA_DIR/config/config.toml"
    sed -i 's/size = 200/size = 60000/g' "$CURRENT_DATA_DIR/config/config.toml"

    # Update app.toml with custom ports and settings
    sed -i 's/minimum-gas-prices = "0.007uaxl"/minimum-gas-prices = "0uaxl"/g' "$CURRENT_DATA_DIR/config/app.toml"
    sed -i "s/address = \"0.0.0.0:9090\"/address = \"${PRIVATE_HOSTS[$INDEX]}:${GRPC_PORTS[$INDEX]}\"/g" "$CURRENT_DATA_DIR/config/app.toml"
    sed -i "s/address = \"0.0.0.0:9091\"/address = \"${PRIVATE_HOSTS[$INDEX]}:${GRPC_WEB_PORTS[$INDEX]}\"/g" "$CURRENT_DATA_DIR/config/app.toml"
    sed -i "s#address = \"tcp://0.0.0.0:1317\"#address = \"tcp://${PRIVATE_HOSTS[$INDEX]}:${API_PORTS[$INDEX]}\"#g" "$CURRENT_DATA_DIR/config/app.toml"

    # Enable REST API
    sed -i '/# Enable defines if the API server should be enabled\./ {n; s/enable = false/enable = true/}' "$CURRENT_DATA_DIR/config/app.toml"
done

# Update persistent peers
echo "Updating persistent_peers..."
PERSISTENT_PEERS=""
for ((j = 0; j < $NODE_COUNT; j++)); do
    PEER_ID=$(axelard tendermint show-node-id --home "$TESTDIR/node$j")
    if [ -n "$PEER_ID" ]; then
        PERSISTENT_PEERS+="${PEER_ID}@${PRIVATE_HOSTS[$j]}:${P2P_PORTS[$j]},"
    else
        echo "Failed to retrieve PEER_ID for node${j}. Skipping..."
    fi
done

# Remove the trailing comma
PERSISTENT_PEERS=${PERSISTENT_PEERS%,}

# Apply persistent peers to all nodes
for ((i = 0; i < $NODE_COUNT; i++)); do
    CURRENT_DATA_DIR="$TESTDIR/node$i"
    sed -i "s/persistent_peers = \"\"/persistent_peers = \"$PERSISTENT_PEERS\"/g" "$CURRENT_DATA_DIR/config/config.toml"
done

echo "Persistent peers updated: $PERSISTENT_PEERS"