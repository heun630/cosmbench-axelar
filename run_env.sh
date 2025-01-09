#!/bin/bash

export TESTDIR="nodes"

export HOSTS=()
export HOSTS_0="127.0.0.1"
export HOSTS_1="127.0.0.1"
export HOSTS_2="127.0.0.1"
export HOSTS_3="127.0.0.1"
export NODE_COUNT=4 # 노드 수
export IP_COUNT=1
export PROC_PER_IP=$((NODE_COUNT / IP_COUNT)) # IP별 동작할 프로세스 수

export PRIVATE_HOSTS=()
export PRIVATE_HOSTS_0="127.0.0.1"
export PRIVATE_HOSTS_1="127.0.0.1"
export PRIVATE_HOSTS_2="127.0.0.1"
export PRIVATE_HOSTS_3="127.0.0.1"

export RPC_PORTS=()
export RPC_PORTS_0="22000"
export RPC_PORTS_1="22001"
export RPC_PORTS_2="22002"
export RPC_PORTS_3="22003"

export P2P_PORTS=()
export P2P_PORTS_0="22100"
export P2P_PORTS_1="22101"
export P2P_PORTS_2="22102"
export P2P_PORTS_3="22103"

export API_PORTS=()
export API_PORTS_0="22200"
export API_PORTS_1="22201"
export API_PORTS_2="22202"
export API_PORTS_3="22203"

export GRPC_PORTS=()
export GRPC_PORTS_0="22300"
export GRPC_PORTS_1="22301"
export GRPC_PORTS_2="22302"
export GRPC_PORTS_3="22303"

export PPROF_PORTS=()
export PPROF_PORTS_0="22400"
export PPROF_PORTS_1="22401"
export PPROF_PORTS_2="22402"
export PPROF_PORTS_3="22403"

export PROXYAPP_PORTS=()
export PROXYAPP_PORTS_0="22500"
export PROXYAPP_PORTS_1="22501"
export PROXYAPP_PORTS_2="22502"
export PROXYAPP_PORTS_3="22503"

export GRPC_WEB_PORTS=()
export GRPC_WEB_PORTS_0="22600"
export GRPC_WEB_PORTS_1="22601"
export GRPC_WEB_PORTS_2="22602"
export GRPC_WEB_PORTS_3="22603"

export USE_SEED_NODE="disable"