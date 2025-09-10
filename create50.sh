#!/bin/bash
set -e

NODE_COUNT=${1:-50} # default 50
KAD_ID="0000000000000000000000000000000000000000" # <-- same ID for all nodes

cat <<EOF > docker-compose.yml
version: "3.9"
services:
EOF

for i in $(seq 1 $NODE_COUNT); do
  echo "  node$i:" >> docker-compose.yml
  if [ "$i" -eq 1 ]; then
    echo "    build: ." >> docker-compose.yml
    echo "    image: kademlia_node" >> docker-compose.yml
    echo "    environment:" >> docker-compose.yml
    echo "      - PORT=9001" >> docker-compose.yml
    echo "      - KAD_ID=$KAD_ID" >> docker-compose.yml
  else
    echo "    image: kademlia_node" >> docker-compose.yml
    echo "    environment:" >> docker-compose.yml
    echo "      - PORT=$((9000 + i))" >> docker-compose.yml
    echo "      - PEER=node1:9001" >> docker-compose.yml
    echo "      - KAD_ID=$KAD_ID" >> docker-compose.yml
  fi
  echo "    networks: [testnet]" >> docker-compose.yml
done

cat <<EOF >> docker-compose.yml
networks:
  testnet:
    driver: bridge
EOF
