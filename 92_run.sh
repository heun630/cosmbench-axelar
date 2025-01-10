#!/bin/bash

# Load environment variables
source ./env.sh
source ./run_env.sh

# Get node index from the first argument
INDEX=$1
CURRENT_DATA_DIR=$TESTDIR/node$INDEX
LOG_FILE="output$INDEX.log"

# Print the command that will be run (for debugging purposes)
echo "$BINARY start --home $CURRENT_DATA_DIR"

# Remove existing log file
rm -f "$LOG_FILE"

# Start the binary and filter lines containing "committed state"
$BINARY start --home "$CURRENT_DATA_DIR" 2>&1 | while IFS= read -r line; do
  if [[ "$line" == *"committed state"* ]]; then
    # Append Unix millisecond timestamp to the line and write to log file
    echo "$(date '+%s%3N') $line" >> "$LOG_FILE"
  fi
done
