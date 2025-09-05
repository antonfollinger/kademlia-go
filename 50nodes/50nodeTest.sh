#!/bin/bash
set -e 

NETWORK_NAME=testnet
NODE_COUNT=50
IMAGE=busybox   

# Clean up if network exists
if docker network inspect "$NETWORK_NAME" >/dev/null 2>&1; then
    docker network rm "$NETWORK_NAME"
fi

# Create network
docker network create "$NETWORK_NAME"

# Start containers
for i in $(seq 1 "$NODE_COUNT"); do
    docker run -dit --name "node$i" --network "$NETWORK_NAME" "$IMAGE" sh
done

echo "Started $NODE_COUNT containers in network '$NETWORK_NAME'."


# In node 2
## nc -l -p 1234 

# In node 1
## echo 'Hello from node1' | nc node2 1234



