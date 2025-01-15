#!/bin/bash

# Load environment variables relative to the script's location
SCRIPT_DIR=$(dirname "$0")
source "$SCRIPT_DIR/env.sh"

# Create accounts and add them to genesis
for ((i = 0; i < $NODE_COUNT; i++)); do
    CURRENT_DATA_DIR=$NODE_ROOT_DIR/node$i

    for ((j = 0; j < $ACCOUNT_COUNT_PER_LOOP; j++)); do
        NUMBER=$(($((i * $ACCOUNT_COUNT_PER_LOOP)) + j))
        ACCOUNT_NAME=$ACCOUNT_NAME_PREFIX$NUMBER

        $BINARY keys add $ACCOUNT_NAME --keyring-backend $KEYRING_BACKEND --home $CURRENT_DATA_DIR
        ACCOUNT_ADDRESS=$($BINARY keys show $ACCOUNT_NAME -a --home $CURRENT_DATA_DIR --keyring-backend $KEYRING_BACKEND)
        echo "$ACCOUNT_ADDRESS"

        $BINARY add-genesis-account $ACCOUNT_ADDRESS 10000000000000$UNIT --home $GENESIS_DIR
    done
done

# Backup keyring info
rm -rf $ACCOUNT_DIR
mkdir -p $ACCOUNT_DIR

for ((i = 0; i < $NODE_COUNT; i++)); do
    CURRENT_DATA_DIR=$NODE_ROOT_DIR/node$i

    mkdir -p $ACCOUNT_DIR/node$i
    cp -f $CURRENT_DATA_DIR/keyring-test/*.info $ACCOUNT_DIR/node$i
done

echo "### 3_create_account.sh done"

# Display account addresses
for ((i = 0; i < $NODE_COUNT; i++)); do
    CURRENT_DATA_DIR=$NODE_ROOT_DIR/node$i

    for ((j=0;j<$ACCOUNT_COUNT_PER_LOOP;j++)); do
        NUMBER=$(($((i * $ACCOUNT_COUNT_PER_LOOP)) + j))
        ACCOUNT_NAME=$ACCOUNT_NAME_PREFIX$NUMBER
        ACCOUNT_ADDRESS=$($BINARY keys show $ACCOUNT_NAME -a --home $CURRENT_DATA_DIR --keyring-backend $KEYRING_BACKEND)
        echo "$ACCOUNT_ADDRESS"
    done
done
