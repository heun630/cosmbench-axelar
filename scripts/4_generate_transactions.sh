#!/bin/bash

# Load environment variables relative to the script's location
SCRIPT_DIR=$(dirname "$0")
source "$SCRIPT_DIR/env.sh"

# Clean and recreate necessary directories
rm -rf $UNSIGNED_TX_ROOT_DIR
rm -rf $SIGNED_TX_ROOT_DIR
rm -rf $ENCODED_TX_ROOT_DIR

mkdir -p $UNSIGNED_TX_ROOT_DIR
mkdir -p $SIGNED_TX_ROOT_DIR
mkdir -p $ENCODED_TX_ROOT_DIR

# Generate transactions
for ((i=0; i<$NODE_COUNT; i++)); do
    CURRENT_DATA_DIR=$NODE_ROOT_DIR/node$i
    for ((j=0; j<$ACCOUNT_COUNT_PER_LOOP; j++)); do
        NUMBER=$((i * $ACCOUNT_COUNT_PER_LOOP + j))
        ACCOUNT_NUMBER=$((NUMBER + 4))  # Adjust for validators
        ACCOUNT_NAME=${ACCOUNT_NAME_PREFIX}${NUMBER}

        # Retrieve account address
        ACCOUNT_ADDRESS=$($BINARY keys show $ACCOUNT_NAME -a --home $CURRENT_DATA_DIR --keyring-backend test)

        # Generate unsigned transaction
        $BINARY tx bank send $ACCOUNT_ADDRESS $ACCOUNT_ADDRESS $SEND_AMOUNT$UNIT \
            --chain-id $CHAIN_ID \
            --home $CURRENT_DATA_DIR \
            --keyring-backend test \
            --generate-only > $UNSIGNED_TX_ROOT_DIR/$UNSIGNED_TX_PREFIX$NUMBER

        # Sign the transaction
        $BINARY tx sign $UNSIGNED_TX_ROOT_DIR/$UNSIGNED_TX_PREFIX$NUMBER \
            --chain-id $CHAIN_ID \
            --from $ACCOUNT_NAME \
            --home $CURRENT_DATA_DIR \
            --offline \
            --sequence 0 \
            --account-number $ACCOUNT_NUMBER \
            --keyring-backend test > $SIGNED_TX_ROOT_DIR/$SIGNED_TX_PREFIX$NUMBER 2>&1

        # Encode the signed transaction
        ENCODED=$($BINARY tx encode $SIGNED_TX_ROOT_DIR/$SIGNED_TX_PREFIX$NUMBER)

        # Save the encoded transaction
        echo $ENCODED > $ENCODED_TX_ROOT_DIR/$ENCODED_TX_PREFIX$NUMBER
    done
done

echo "### 4_generate_transactions.sh done"
