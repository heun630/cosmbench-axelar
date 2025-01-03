#!/bin/bash

# Start the script
echo "Initializing Axelar Core..."

# Get input from the user
read -p "Enter the data directory path (e.g., /data/axelar): " DATA_DIR
read -p "Enter the node name (e.g., node1): " NODE_NAME
read -p "Enter the chain ID (e.g., my-private-chain): " CHAIN_ID

# Step 1: Clone the Git repository and build
echo "Cloning Axelar Core source and building..."
git clone https://github.com/axelarnetwork/axelar-core.git
cd axelar-core || exit
make build || exit
cd ..

# Step 2: Create the data directory and initialize
echo "Creating data directory and initializing..."
mkdir -p "$DATA_DIR"
axelard init "$NODE_NAME" --chain-id "$CHAIN_ID" --home "$DATA_DIR" || exit

# Step 3: Add the validator key
echo "Creating a validator key..."
KEY_OUTPUT=$(axelard keys add validator --home "$DATA_DIR" --output json)
echo "Full key output: $KEY_OUTPUT"

VALIDATOR_KEY=$(echo "$KEY_OUTPUT" | jq -r '.address')

if [[ -z "$VALIDATOR_KEY" ]]; then
  echo "Validator key creation failed. Attempting to retrieve the key..."
  VALIDATOR_KEY=$(axelard keys show validator --home "$DATA_DIR" -a)
  if [[ -z "$VALIDATOR_KEY" ]]; then
    echo "Failed to create or retrieve the validator key. Please check your setup." >&2
    exit 1
  fi
fi

echo "Validator key creation output: $VALIDATOR_KEY"

# Step 4: Add a genesis account
echo "Adding a genesis account..."
axelard add-genesis-account "$VALIDATOR_KEY" 2000000000uaxl --home "$DATA_DIR" || exit

# Step 5: Generate a genesis transaction
echo "Generating a genesis transaction..."
axelard gentx validator 1000000000uaxl \
  --chain-id "$CHAIN_ID" \
  --home "$DATA_DIR" || exit

# Update genesis.json to replace "stake" with "uaxl"
echo "Updating genesis.json to replace \"stake\" with \"uaxl\"..."
sed -i 's/\"stake\"/\"uaxl\"/g' "$DATA_DIR/config/genesis.json"

# Step 6: Collect genesis transactions
echo "Collecting genesis transactions..."
axelard collect-gentxs --home "$DATA_DIR" || exit

# Step 7: Start the node
echo "Starting the node..."
nohup axelard start --home "$DATA_DIR" > "$DATA_DIR/node1.log" 2>&1 &

echo "Axelar Core initialization is complete!"