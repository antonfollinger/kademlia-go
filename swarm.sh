#!/bin/bash
set -e

NODE_COUNT=${1:-50} # default 50
PEER_COUNT=$((NODE_COUNT-1))

cat <<EOF > docker-stack.yml
services:
  bootstrap:
    build: .
    image: kademlia_node
    environment:
      - PORT=9001
      - ISBOOTSTRAP=TRUE
    deploy:
      replicas: 1
    networks:
      - testnet

  peer:
    image: kademlia_node
    environment:
      - ISBOOTSTRAP=FALSE
      - BOOTSTRAPNODE=bootstrap
    depends_on:
      - bootstrap
    deploy:
      replicas: $PEER_COUNT
    networks:
      - testnet

networks:
  testnet:
    driver: overlay
EOF