#!/bin/bash

# Kill swarm network
docker stack rm kademlia 2>/dev/null || true
sleep 5
docker network rm testnet 2>/dev/null || true

# Stop all running containers
docker stop $(docker ps -aq)

# Remove all containers
docker rm $(docker ps -aq)


docker image prune -a

