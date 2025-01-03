#!/bin/bash

# Monitor logs for all Axelar Core nodes

# Get input from the user
read -p "Enter the base data directory path (default: /data/axelar): " BASE_DIR
BASE_DIR=${BASE_DIR:-/data/axelar}

read -p "Enter the number of nodes to set up (default: 2): " NODE_COUNT
NODE_COUNT=${NODE_COUNT:-2}

# Validate input
if [[ -z "$BASE_DIR" || -z "$NODE_COUNT" ]]; then
  echo "Base directory or node count cannot be empty."
  exit 1
fi

# Array to store background process IDs
PIDS=()

# Set up a trap to terminate all background processes on exit
cleanup() {
  echo "Terminating all log monitors..."
  for PID in "${PIDS[@]}"; do
    kill "$PID" 2>/dev/null
  done
  exit
}
trap cleanup SIGINT SIGTERM

# Monitor logs for all nodes
for i in $(seq 1 "$NODE_COUNT"); do
  NODE_LOG="${BASE_DIR}/node${i}/node${i}.log"
  if [[ -f "$NODE_LOG" ]]; then
    echo "Tailing log for node${i}..."
    tail -f "$NODE_LOG" | sed "s/^/[node${i}] /" &
    PIDS+=($!) # Store the process ID of the background process
  else
    echo "Log file $NODE_LOG not found for node${i}. Skipping..."
  fi
done

# Wait for all background processes to complete
wait