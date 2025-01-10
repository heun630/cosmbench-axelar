#!/bin/bash
source ./env.sh
source ./run_env.sh

INDEX=$1
CURRENT_DATA_DIR=$TESTDIR/node$INDEX
echo "$BINARY start --home $CURRENT_DATA_DIR"
#$BINARY start --home $CURRENT_DATA_DIR


rm -f output$INDEX.log

# committed state라는 문장이 들어있는 문장에 Unix millisecond timestamp 출력
$BINARY start --home $CURRENT_DATA_DIR 2>&1 | while IFS= read -r line; do
  if [[ "$line" == *"committed state"* ]]; then
    echo "$(date '+%s%3N') $line"
  fi
done >> output$INDEX.log