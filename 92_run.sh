#!/bin/bash
source ./env.sh
source ./run_env.sh

INDEX=$1
CURRENT_DATA_DIR=$TESTDIR/node$INDEX
echo "$BINARY start --home $CURRENT_DATA_DIR"
$BINARY start --home $CURRENT_DATA_DIR
