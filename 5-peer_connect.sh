#!/bin/bash

# Get input from the user
read -p "Enter the base data directory path (default: /data/axelar): " BASE_DIR
BASE_DIR=${BASE_DIR:-/data/axelar}

read -p "Enter the number of nodes to set up (default: 2): " NODE_COUNT
NODE_COUNT=${NODE_COUNT:-2}

# Validate input
if [[ -z "$BASE_DIR" || -z "$NODE_COUNT" ]]; then
  echo "Base directory or node count cannot be empty."
  exit 1
fi

# Build persistent peers string
echo "Building persistent peers..."
PERSISTENT_PEERS=""
for j in $(seq 1 "$NODE_COUNT"); do
  PEER_ID=$(axelard tendermint show-node-id --home "${BASE_DIR}/node${j}")

  PEER_PORT=$((26656 + 1000 * j)) # Other nodes use offset ports

  PERSISTENT_PEERS+="${PEER_ID}@127.0.0.1:${PEER_PORT},"
done
PERSISTENT_PEERS=${PERSISTENT_PEERS%,} # Remove trailing comma

# Apply persistent peers to all nodes
echo "Applying persistent peers to all nodes..."
for i in $(seq 1 "$NODE_COUNT"); do
  NODE_DIR="${BASE_DIR}/node${i}"
  sed -i "s/^persistent_peers = .*/persistent_peers = \"$PERSISTENT_PEERS\"/" "$NODE_DIR/config/config.toml"
  echo "Updated persistent_peers for node${i}"
done

echo "Persistent peers configuration is complete!"