#!/bin/bash
set -e

NODE_COUNT=10

cat <<EOF > docker-compose.yml
version: "3.9"
services:
EOF


for i in $(seq 1 $NODE_COUNT); do
  echo "  node$i:" >> docker-compose.yml
  if [ "$i" -eq 1 ]; then
    echo "    build: ." >> docker-compose.yml
    echo "    image: kademlia_node" >> docker-compose.yml
  else
    echo "    image: kademlia_node" >> docker-compose.yml
  fi
  echo "    environment:" >> docker-compose.yml
  echo "      - PORT=$((9000 + i))" >> docker-compose.yml
  if [ "$i" -ne 1 ]; then
    echo "      - PEER=node1:9001" >> docker-compose.yml
  fi
  echo "    networks: [testnet]" >> docker-compose.yml
done

cat <<EOF >> docker-compose.yml
networks:
  testnet:
    driver: bridge
EOF