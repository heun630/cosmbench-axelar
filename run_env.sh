#!/bin/bash

TESTDIR="nodes"

HOSTS=()

HOSTS[0]="127.0.0.1"
HOSTS[1]="127.0.0.1"
HOSTS[2]="127.0.0.1"
HOSTS[3]="127.0.0.1"
NODE_COUNT=4 #나중에 txs_editor.py할 때 파싱하는 과정에서 #HOOSTS[@]를 쓰면 안되기 때문에 직접 값을 입력해야함
IP_COUNT=1
PROC_PER_IP=$((NODE_COUNT / IP_COUNT)) #IP별 동작할 프로세스 수

PRIVATE_HOSTS[0]="127.0.0.1"
PRIVATE_HOSTS[1]="127.0.0.1"
PRIVATE_HOSTS[2]="127.0.0.1"
PRIVATE_HOSTS[3]="127.0.0.1"

RPC_PORTS=()
RPC_PORTS[0]="22000"
RPC_PORTS[1]="22001"
RPC_PORTS[2]="22002"
RPC_PORTS[3]="22003"

P2P_PORTS=()
P2P_PORTS[0]="22100"
P2P_PORTS[1]="22101"
P2P_PORTS[2]="22102"
P2P_PORTS[3]="22103"

###############################

API_PORTS=()
API_PORTS[0]="22200"
API_PORTS[1]="22201"
API_PORTS[2]="22202"
API_PORTS[3]="22203"

GRPC_PORTS=()
GRPC_PORTS[0]="22300"
GRPC_PORTS[1]="22301"
GRPC_PORTS[2]="22302"
GRPC_PORTS[3]="22303"

PPROF_PORTS=()
PPROF_PORTS[0]="22400"
PPROF_PORTS[1]="22401"
PPROF_PORTS[2]="22402"
PPROF_PORTS[3]="22403"

PROXYAPP_PORTS=()
PROXYAPP_PORTS[0]="22500"
PROXYAPP_PORTS[1]="22501"
PROXYAPP_PORTS[2]="22502"
PROXYAPP_PORTS[3]="22503"

#####################################

NODE_IDS=()
NODE_IDS[0]="38d49d4336e8118ab00bd53e59e28ba2679f21f4"
NODE_IDS[1]="9fcb52cb6bb087d9a70d9a4f6d1a4b6b0dd7f1a4"
NODE_IDS[2]="836d16ac9583081bc466ccad566fc1e3e3129622"
NODE_IDS[3]="2d4ccf7b764315cd2e972b825c769e570945d484"

USE_SEED_NODE="disable"