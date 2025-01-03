#!/bin/bash

# Start the script
echo "Initializing multiple Axelar Core nodes..."

# Get input from the user with default values
read -p "Enter the base data directory path (default: /data/axelar): " BASE_DIR
BASE_DIR=${BASE_DIR:-/data/axelar}

read -p "Enter the base chain ID (default: my-private-chain): " BASE_CHAIN_ID
BASE_CHAIN_ID=${BASE_CHAIN_ID:-my-private-chain}

read -p "Enter the number of nodes to set up (default: 2): " NODE_COUNT
NODE_COUNT=${NODE_COUNT:-2}

# Step 1: Clone the Git repository and build
echo "Cloning Axelar Core source and building..."
git clone https://github.com/axelarnetwork/axelar-core.git
cd axelar-core || exit
make build || exit
cd ..

# Shared CHAIN_ID for all nodes
CHAIN_ID="$BASE_CHAIN_ID"

# Step 2: Loop through and create nodes
for i in $(seq 1 "$NODE_COUNT"); do
  NODE_DIR="${BASE_DIR}/node${i}"
  NODE_NAME="node${i}"

  # Calculate unique ports for each node
  P2P_PORT=$((26656 + 1000 * i))
  RPC_PORT=$((26657 + 1000 * i))
  API_PORT=$((1317 + 1000 * i))
  PROXY_APP_PORT=$((26658 + 1000 * i))
  PROMETHEUS_PORT=$((26660 + 1000 * i))
  PPROF_PORT=$((6060 + 1000 * i))
  GRPC_PORT=$((9090 + 1000 * i))
  GRPC_WEB_PORT=$((9091 + 1000 * i))
  ROSETTA_PORT=$((8080 + 1000 * i))

  echo "Setting up $NODE_NAME with CHAIN_ID $CHAIN_ID..."

  # Create the data directory and initialize
  mkdir -p "$NODE_DIR"
  axelard init "$NODE_NAME" --chain-id "$CHAIN_ID" --home "$NODE_DIR" || exit

  # Add the validator key
  echo "Creating a validator key for $NODE_NAME..."
  KEY_OUTPUT=$(axelard keys add validator --home "$NODE_DIR" --output json)
  VALIDATOR_KEY=$(echo "$KEY_OUTPUT" | jq -r '.address')

  if [[ -z "$VALIDATOR_KEY" ]]; then
    echo "Validator key creation failed for $NODE_NAME. Attempting to retrieve the key..."
    VALIDATOR_KEY=$(axelard keys show validator --home "$NODE_DIR" -a)
    if [[ -z "$VALIDATOR_KEY" ]]; then
      echo "Failed to create or retrieve the validator key for $NODE_NAME. Skipping..." >&2
      continue
    fi
  fi

  echo "Validator key creation output for $NODE_NAME: $VALIDATOR_KEY"

  # Add a genesis account
  axelard add-genesis-account "$VALIDATOR_KEY" 2000000000uaxl --home "$NODE_DIR" || exit

  # Handle genesis transactions for the first node
  if [[ "$i" -eq 1 ]]; then
    axelard gentx validator 1000000000uaxl \
      --chain-id "$CHAIN_ID" \
      --home "$NODE_DIR" || exit

    echo "Updating genesis.json for $NODE_NAME to replace \"stake\" with \"uaxl\"..."
    sed -i 's/\"stake\"/\"uaxl\"/g' "$NODE_DIR/config/genesis.json"

    axelard collect-gentxs --home "$NODE_DIR" || exit
  else
    cp "${BASE_DIR}/node1/config/genesis.json" "$NODE_DIR/config/"
  fi

  # Update configuration for unique ports in config.toml
  sed -i "s/26656/${P2P_PORT}/g" "$NODE_DIR/config/config.toml"
  sed -i "s/26657/${RPC_PORT}/g" "$NODE_DIR/config/config.toml"
  sed -i "s/26658/${PROXY_APP_PORT}/g" "$NODE_DIR/config/config.toml"
  sed -i "s/:26660/:${PROMETHEUS_PORT}/g" "$NODE_DIR/config/config.toml"
  sed -i "s/localhost:6060/localhost:${PPROF_PORT}/g" "$NODE_DIR/config/config.toml"

  # Update configuration for unique ports in app.toml
  sed -i "s/0.0.0.0:1317/0.0.0.0:${API_PORT}/g" "$NODE_DIR/config/app.toml"
  sed -i "s/0.0.0.0:9090/0.0.0.0:${GRPC_PORT}/g" "$NODE_DIR/config/app.toml"
  sed -i "s/0.0.0.0:9091/0.0.0.0:${GRPC_WEB_PORT}/g" "$NODE_DIR/config/app.toml"
  sed -i "s/:8080/:${ROSETTA_PORT}/g" "$NODE_DIR/config/app.toml"
done

# Step 3: Configure persistent peers
echo "Setting persistent peers..."
PERSISTENT_PEERS=""
for j in $(seq 1 "$NODE_COUNT"); do
  PEER_ID=$(axelard tendermint show-node-id --home "${BASE_DIR}/node${j}")
  PEER_PORT=$((26656 + 1000 * j))
  PERSISTENT_PEERS+="${PEER_ID}@127.0.0.1:${PEER_PORT},"
done
PERSISTENT_PEERS=${PERSISTENT_PEERS%,} # Remove trailing comma

for i in $(seq 1 "$NODE_COUNT"); do
  NODE_DIR="${BASE_DIR}/node${i}"
  sed -i "s/persistent_peers = .*/persistent_peers = \"$PERSISTENT_PEERS\"/" "$NODE_DIR/config/config.toml"
done

# Step 4: Start all nodes
echo "Starting all nodes..."
for i in $(seq 1 "$NODE_COUNT"); do
  NODE_DIR="${BASE_DIR}/node${i}"
  NODE_NAME="node${i}"
  nohup axelard start --home "$NODE_DIR" > "$NODE_DIR/${NODE_NAME}.log" 2>&1 &
done

echo "Axelar Core multi-node setup is complete!"