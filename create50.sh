#!/bin/bash
set -e

NODE_COUNT=${1:-10} # default 50

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
    echo "      - BOOTSTRAP=TRUE" >> docker-compose.yml
  else
    echo "    image: kademlia_node" >> docker-compose.yml
    echo "    environment:" >> docker-compose.yml
    echo "      - BOOTSTRAP=False" >> docker-compose.yml
    echo "      - PEER=node1" >> docker-compose.yml
  fi
  echo "    networks: [testnet]" >> docker-compose.yml
done

cat <<EOF >> docker-compose.yml
networks:
  testnet:
    driver: bridge
EOF
