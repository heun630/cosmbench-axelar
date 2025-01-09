#!/bin/bash

export CHAIN_NAME="axelar"
export CHAIN_ID=$CHAIN_NAME"-cosmbench" # 블록체인 ID
export BINARY="/usr/local/bin/"$CHAIN_NAME"d"
export MONIKER="cosmbench"
# export ADDRESS_PREFIX="mssong"

export KEYRING_BACKEND="test" # Select keyring's backend (os|file|test) (default "os")

export NODE_COUNT=4 # 노드 수
export ACCOUNT_COUNT_PER_LOOP=4 # 노드 당 생성할 어카운트 수
# 총 어카운트 수 = NODE_COUNT * ACCOUNT_COUNT_PER_LOOP

export UNIT="uaxl" ## 전송 코인 이름
export SEND_AMOUNT=100 ## 전송 코인 갯수

export NODE_ROOT_DIR=$CHAIN_ID"_nodes" # 노드들을 가지고 있는 디렉토리
export ACCOUNT_DIR=$CHAIN_ID"_accounts"

export ACCOUNT_NAME_PREFIX="account_" # account 생성 시 .info파일 이름

export GENESIS_DIR=$NODE_ROOT_DIR"/node0" # 기준이 될 genesis.json을 가지고 있는 노드

export UNSIGNED_TX_PREFIX="unsigned_tx_"
export SIGNED_TX_PREFIX="signed_tx_"
export ENCODED_TX_PREFIX="encoded_tx_"

export UNSIGNED_TX_ROOT_DIR=$CHAIN_ID"_unsigned_txs"
export SIGNED_TX_ROOT_DIR=$CHAIN_ID"_signed_txs"
export ENCODED_TX_ROOT_DIR=$CHAIN_ID"_encoded_txs"

export DEPLOY_DIR="deploy_run_nodes_scripts"