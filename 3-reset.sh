#!/bin/bash

# Script to stop and clean up all Axelar Core nodes

# Kill all running axelard processes
echo "Stopping all running Axelar Core nodes..."
ps -ef | grep axelard | grep -v grep | awk '{print $2}' | xargs -r kill -9

# Clean up the data directory
BASE_DIR="/data/axelar"
echo "Cleaning up data directory at $BASE_DIR..."
rm -rf "$BASE_DIR"/*

echo "All Axelar Core nodes have been stopped and the data directory has been cleaned."