#!/bin/bash

# Load environment variables relative to the scripts directory
SCRIPT_DIR=$(dirname "$0")
source "$SCRIPT_DIR/env.sh"

rm -rf $NODE_ROOT_DIR

echo "The number of nodes: $NODE_COUNT"
for ((i=0; i<NODE_COUNT; i++)); do
    CURRENT_DATA_DIR=$NODE_ROOT_DIR/node$i
    $BINARY init $MONIKER$i --chain-id $CHAIN_ID --home $CURRENT_DATA_DIR

    cp -f $CURRENT_DATA_DIR/config/genesis.json $CURRENT_DATA_DIR/config/sample_genesis.json
done
echo "The number of nodes: $NODE_COUNT"

echo "[ SHOW NODE ID ]"
for ((i=0; i<NODE_COUNT; i++)); do
    CURRENT_DATA_DIR=$NODE_ROOT_DIR/node$i
    $BINARY tendermint show-node-id --home $CURRENT_DATA_DIR
done