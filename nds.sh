#!/bin/bash
set -e

NODE_COUNT=${1:-50} # default 50

cat <<EOF > docker-compose.yml
services:
EOF

for i in $(seq 1 $NODE_COUNT); do
  echo "  node$i:" >> docker-compose.yml
  if [ "$i" -eq 1 ]; then
    echo "    build: ." >> docker-compose.yml
    echo "    stdin_open: true" >> docker-compose.yml
    echo "    tty: true" >> docker-compose.yml
    echo "    environment:" >> docker-compose.yml
    echo "      - PORT=9001" >> docker-compose.yml
    echo "      - ISBOOTSTRAP=TRUE" >> docker-compose.yml
    echo "      - ENABLECLI=TRUE" >> docker-compose.yml
    
  else
    echo "    build: ." >> docker-compose.yml
    echo "    stdin_open: true" >> docker-compose.yml
    echo "    tty: true" >> docker-compose.yml
    echo "    environment:" >> docker-compose.yml
    echo "      - PORT=$((9000 + i))" >> docker-compose.yml
    echo "      - BOOTSTRAPNODE=node1" >> docker-compose.yml
    echo "      - ISBOOTSTRAP=FALSE" >> docker-compose.yml
  fi
  echo "    networks: [testnet]" >> docker-compose.yml
done

cat <<EOF >> docker-compose.yml
networks:
  testnet:
    driver: bridge
EOF
