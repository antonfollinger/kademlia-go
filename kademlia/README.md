# Kademlia
A peer-to-peer Distributed Hash Table (DHT) using an XOR metric for distance between points (nodes/keys) in a 160-bit key space

## Node
A participant in the DHT with a random ID somewhere in the key space and maintains <i>n</i> k-buckets.

## Contact
Triplet consisting of <IP Address, UDP port, Node ID> describing how one node refers to another.

## k-Bucket
A list of up to <i>k</i> contacts, covering a specific distance range  [2<sup>i</sup>, 2<sup>i+1</sup>) from the local node, using the XOR metric. Buckets nearer to the local node are more granular; farther buckets cover wider ranges.

## Routing table
The collection of all k-buckets for a node. This gives the node neighbors spread across the ID space, enabling efficient O(log N) lookups.
