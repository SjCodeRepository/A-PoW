# A-PoW
Code Implementation of the Paper 'A-PoW: An Adaptive PoW Consensus Algorithm for Blockchain in IoT Networks'

This code is for the A-PoW related experimental code. The experimental test network consists of the storage layer, network layer, consensus layer, and service layer. The 'storage' folder contains code related to the storage layer, using the BoltDB database to achieve persistent storage of blocks and other data. The 'network' folder implements network data transmission. The 'core' folder contains core data structures, including transactions and blocks. The 'blockchain' folder defines the blockchain data structure. The 'Service' and 'worknode' folders define the operational logic for service nodes and worker nodes, respectively. The 'A-PoW' folder contains the implementation of methods related to the A-PoW consensus algorithm in the consensus layer.

Steps to start the test network:

1.Launch a local Docker Golang container, copy the 'getnodetimes' folder into the container, and run the 'main.go' file in the folder to obtain the average mining difficulty time (NodeTimes) for each node group.
2.Replace the NodeTimes in 'main.go' under the 'runservice' folder with the results obtained in step 1, and then run the command 'docker-compose up -d' in the 'APoWCode' folder to start the test network.
