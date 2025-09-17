#!/bin/bash
set -e

# Number of nodes (default 50, can be overridden by first argument)
NODE_COUNT=${1:-50}

./reset.sh
sleep 2

echo "Building Docker image..."
docker build -t kademlia_node .

echo "Initializing Docker Swarm (if not already)..."
docker swarm init || true

echo "Generating docker-stack.yml for $NODE_COUNT nodes..."
./swarm.sh $NODE_COUNT

echo "Deploying stack to Docker Swarm..."
docker stack deploy -c docker-stack.yml kademlia

echo "Deployment complete!"