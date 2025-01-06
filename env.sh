#!/bin/bash

CHAIN_NAME="axelar"
CHAIN_ID=$CHAIN_NAME"-cosmbench" #블록체인 ID
BINARY="/usr/local/bin/"$CHAIN_NAME"d"
MONIKER="cosmbench"
# ADDRESS_PREFIX="mssong"

KEYRING_BACKEND="test" # Select keyring's backend (os|file|test) (default "os")

NODE_COUNT=4 #노드 수
ACCOUNT_COUNT_PER_LOOP=4 #노드 당 생성할 어카운트 수
#즉, 총 어카운트 수 = NODE_COUNT * ACCOUNT_COUNT_PER_LOOP

UNIT="uaxl" ##전송 코인 이름
SEND_AMOUNT=100 ##전송 코인 갯수

NODE_ROOT_DIR=$CHAIN_ID"_nodes" #노드들을 가지고 있는 디렉토리
ACCOUNT_DIR=$CHAIN_ID"_accounts"

ACCOUNT_NAME_PREFIX="account_" #account 생성 시 .info파일 이름

GENESIS_DIR=$NODE_ROOT_DIR"/node0" #기준이 될 genesis.json을 가지고 있는 노드

UNSIGNED_TX_PREFIX="unsigned_tx_"
SIGNED_TX_PREFIX="signed_tx_"
ENCODED_TX_PREFIX="encoded_tx_"

UNSIGNED_TX_ROOT_DIR=$CHAIN_ID"_unsigned_txs"
SIGNED_TX_ROOT_DIR=$CHAIN_ID"_signed_txs"
ENCODED_TX_ROOT_DIR=$CHAIN_ID"_encoded_txs"

DEPLOY_DIR="deploy_run_nodes_scripts"
