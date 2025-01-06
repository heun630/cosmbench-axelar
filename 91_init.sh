#!/bin/bash
source ./env.sh
source ./run_env.sh

if [ -d "$TESTDIR" ]; then
    echo "Removing $TESTDIR directory..."
    rm -rf $TESTDIR
fi

cp -rf $NODE_ROOT_DIR $TESTDIR

for ((i = 1; i < $NODE_COUNT; i++))
do
    CURRENT_DATA_DIR=$TESTDIR/node$i
    cp -f $TESTDIR/node0/config/genesis.json $CURRENT_DATA_DIR/config/genesis.json 
done

for ((i = 0; i < $NODE_COUNT; i++))
do
    INDEX=$i
    CURRENT_DATA_DIR=$TESTDIR/node$i

    # doenom stake에서 uaxl로 변경
    sed -i 's/\"stake\"/\"uaxl\"/g' "$CURRENT_DATA_DIR/config/genesis.json"

    # Proxy App PORT 변경
    sed -i "s#proxy_app = \"tcp:\/\/127.0.0.1:26658\"#proxy_app = \"tcp:\/\/${PRIVATE_HOSTS[$INDEX]}:${PROXYAPP_PORTS[$INDEX]}\"#g" $CURRENT_DATA_DIR/config/config.toml

    # RPC PORT 변경
    echo "sed -i "s#laddr = \"tcp:\/\/127.0.0.1:26657\"#laddr = \"tcp:\/\/${PRIVATE_HOSTS[$INDEX]}:${RPC_PORTS[$INDEX]}\"#g" $CURRENT_DATA_DIR/config/config.toml"
    sed -i "s#laddr = \"tcp:\/\/127.0.0.1:26657\"#laddr = \"tcp:\/\/${PRIVATE_HOSTS[$INDEX]}:${RPC_PORTS[$INDEX]}\"#g" $CURRENT_DATA_DIR/config/config.toml

    # P2P PORT 변경
    sed -i "s#laddr = \"tcp:\/\/0.0.0.0:26656\"#laddr = \"tcp:\/\/${PRIVATE_HOSTS[$INDEX]}:${P2P_PORTS[$INDEX]}\"#g" $CURRENT_DATA_DIR/config/config.toml 

    # [pprof port]
    sed -i "s#pprof_laddr = \"localhost:6060\"#pprof_laddr = \"${PRIVATE_HOSTS[$INDEX]}:${PPROF_PORTS[$INDEX]}\"#g" $CURRENT_DATA_DIR/config/config.toml

    # 중복 IP 허용
    sed -i 's/allow_duplicate_ip = false/allow_duplicate_ip = true/g' $CURRENT_DATA_DIR/config/config.toml

    # Mempool Size
    sed -i 's/size = 200/size = 60000/g' $CURRENT_DATA_DIR/config/config.toml

    # Minimum Gas Prices
    sed -i 's/minimum-gas-prices = \"160000000inj\"/minimum-gas-prices = \"0stake\"/g' $CURRENT_DATA_DIR/config/app.toml

    # max_bytes
    # sed -i "s/\"max_bytes\": \"22020096\"/\"max_bytes\": \"88080384\"/g" $CURRENT_DATA_DIR/config/genesis.json

    # [app.toml port]
    
    #sed -i 's/address = \"0.0.0.0:9090\"/address = \"${PRIVATE_HOSTS[$INDEX]}:${GRPC_PORTS[$INDEX]}\"/g' $CURRENT_DATA_DIR/config/app.toml
    #echo "sed -i 's/address = \"0.0.0.0:9090\"/address = "'${PRIVATE_HOSTS[$INDEX]}:${GRPC_PORTS[$INDEX]}'"/g' $CURRENT_DATA_DIR/config/app.toml"
    sed -i "s/address = \"0.0.0.0:9900\"/address = \"${PRIVATE_HOSTS[$INDEX]}:${GRPC_PORTS[$INDEX]}\"/g" $CURRENT_DATA_DIR/config/app.toml

    sed -i "s/address = \"tcp:\/\/0.0.0.0:10337\"/address = \"tcp:\/\/${PRIVATE_HOSTS[$INDEX]}:${API_PORTS[$INDEX]}\"/g" $CURRENT_DATA_DIR/config/app.toml
done

# [persistent peers]
echo "Update persistent_peers"
PERSISTENT_PEERS=""

#for ((i = 0; i < $NODE_COUNT; i++))
#do
#    PEER_ID=$(axelard tendermint show-node-id --home "${TESTDIR}/node${j}")
#    echo "axelard tendermint show-node-id --home "${TESTDIR}/node${j}"
#    PERSISTENT_PEERS=$PERSISTENT_PEERS${PEER_ID}'@'${PRIVATE_HOSTS[$i]}':'${P2P_PORTS[$i]}','
#done
#PERSISTENT_PEERS=${PERSISTENT_PEERS%,} #마지막에 ,를 제거하겠다는 의미

for ((j = 0; j < $NODE_COUNT; j++)); do
    # 각 노드의 PEER_ID를 동적으로 가져옴
    PEER_ID=$(axelard tendermint show-node-id --home "${TESTDIR}/node${j}")

    # PEER_ID 확인 후 추가
    if [ -n "$PEER_ID" ]; then
        PERSISTENT_PEERS+="${PEER_ID}@${PRIVATE_HOSTS[$j]}:${P2P_PORTS[$j]},"
    else
        echo "Failed to retrieve PEER_ID for node${j}. Skipping..."
    fi
done

echo "PERSISTENT_PEERS : "$PERSISTENT_PEERS
for ((i = 0; i < $NODE_COUNT; i++))
do
    CURRENT_DATA_DIR=$TESTDIR/node$i
    # echo "sed -i "s/persistent_peers = \"\"/persistent_peers = \"$PERSISTENT_PEERS\"/g" $CURRENT_DATA_DIR/config/config.toml"
    sed -i "s/persistent_peers = \"\"/persistent_peers = \"$PERSISTENT_PEERS\"/g" $CURRENT_DATA_DIR/config/config.toml
done