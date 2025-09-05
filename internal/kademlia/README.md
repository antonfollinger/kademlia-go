# Kademlia
A peer-to-peer distributed hash table using an XOR metric for distance between points (nodes/keys) in a 160-bit key space

## Node
Initializes with a random ID somewhere in the key space

## Routing Table
For every 0 $\le$ i $\le$ 160, each node keeps a list of <IP Address, UDP port, Node ID> triplets for nodes of distance between 2<sup>i</sup> and 2<sup>i+1</sup>, these lists are called <i>k-buckets</i>.

## Bucket
