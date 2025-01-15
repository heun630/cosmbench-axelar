#!/bin/bash

# Load environment variables
source ./env.sh
source ./run_env.sh

# Get node index from the first argument
INDEX=$1
CURRENT_DATA_DIR=$TESTDIR/node$INDEX

# Set log directory and log file
LOG_DIR="logs"
LOG_FILE="$LOG_DIR/output$INDEX.log"

# Create logs directory if it does not exist
mkdir -p "$LOG_DIR"

# Print the command that will be run (for debugging purposes)
echo "$BINARY start --home $CURRENT_DATA_DIR"

# Remove existing log file
rm -f "$LOG_FILE"

# Start the binary and filter lines containing "committed state"
$BINARY start --home "$CURRENT_DATA_DIR" 2>&1 | while IFS= read -r line; do
  if [[ "$line" == *"committed state"* ]]; then
    # Append Unix millisecond timestamp to the line
    timestamped_line="$(date '+%s%3N') $line"
    # Write to log file and print to screen
    echo "$timestamped_line" | tee -a "$LOG_FILE"
  fi
done
