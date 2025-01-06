#!/bin/bash
source ./env.sh

for ((i = 0; i < $NODE_COUNT; i++))
do
    CURRENT_DATA_DIR=$NODE_ROOT_DIR/node$i
    # rm -rf $CURRENT_DATA_DIR/keyring-test

    for ((j = 0; j < $ACCOUNT_COUNT_PER_LOOP; j++))
    do            
        NUMBER=$(($((i * $ACCOUNT_COUNT_PER_LOOP)) + j))
        # NUMBER=$j    
        ACCOUNT_NAME=$ACCOUNT_NAME_PREFIX$NUMBER

        $BINARY keys add $ACCOUNT_NAME --keyring-backend $KEYRING_BACKEND --home $CURRENT_DATA_DIR
        ACCOUNT_ADDRESS=$($BINARY keys show $ACCOUNT_NAME -a --home $CURRENT_DATA_DIR --keyring-backend $KEYRING_BACKEND)
        echo "$ACCOUNT_ADDRESS"

        echo "$BINARY add-genesis-account $ACCOUNT_ADDRESS 10000000000000$UNIT --home $GENESIS_DIR"
        $BINARY add-genesis-account $ACCOUNT_ADDRESS 10000000000000$UNIT --home $GENESIS_DIR
    done
done

rm -rf $ACCOUNT_DIR
mkdir $ACCOUNT_DIR
for ((i = 0; i < $NODE_COUNT; i++))
do
    CURRENT_DATA_DIR=$NODE_ROOT_DIR/node$i

    mkdir -p $ACCOUNT_DIR/node$i
    cp -f $CURRENT_DATA_DIR/keyring-test/*.info $ACCOUNT_DIR/node$i
    # rm -rf $CURRENT_DATA_DIR/keyring-test
done

echo "### 3_create_account.sh done"

for ((i = 0; i < $NODE_COUNT; i++))
do
    CURRENT_DATA_DIR=$NODE_ROOT_DIR/node$i

    for ((j=0;j<$ACCOUNT_COUNT_PER_LOOP;j++))
    do            
        NUMBER=$(($((i*$ACCOUNT_COUNT_PER_LOOP))+j))
        # NUMBER=$j    
        ACCOUNT_NAME=$ACCOUNT_NAME_PREFIX$NUMBER
        ACCOUNT_ADDRESS=$($BINARY keys show $ACCOUNT_NAME -a --home $CURRENT_DATA_DIR --keyring-backend $KEYRING_BACKEND)
        echo "$ACCOUNT_ADDRESS"
    done
done
