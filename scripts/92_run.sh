#!/bin/bash

# Load environment variables
SCRIPT_DIR=$(dirname "$0")
source "$SCRIPT_DIR/env.sh"
source "$SCRIPT_DIR/run_env.sh"

# Get node index from the first argument
INDEX=$1
CURRENT_DATA_DIR=$TESTDIR/node$INDEX

# Set log directory and log file
LOG_DIR="logs"
LOG_FILE="$LOG_DIR/output$INDEX.log"

## Create logs directory if it does not exist
mkdir -p "$LOG_DIR"
#
## Print the command that will be run (for debugging purposes)
echo "$BINARY start --home $CURRENT_DATA_DIR"
#
## Remove existing log file
rm -f "$LOG_FILE"

$BINARY start --home "$CURRENT_DATA_DIR" 2>&1 | while IFS= read -r line; do
  if [[ "$line" == *"committed state"* ]]; then
    # ANSI 이스케이프 코드 제거
    clean_line=$(echo "$line" | sed -r "s/\x1B\[[0-9;]*[mK]//g")

    # Unix millisecond 타임스탬프 추가
    timestamped_line="$(date '+%s%3N') $clean_line"

    # 로그 파일에 기록 및 화면에 출력
    echo "$timestamped_line" | tee -a "logs/output$INDEX.log"
  fi
done