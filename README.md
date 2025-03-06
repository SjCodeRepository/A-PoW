# A-PoW
Code Implementation of the Paper 'A-PoW: An Adaptive PoW Consensus Algorithm for Blockchain in IoT Networks'

## This code represents a basic implementation of A-PoW. 
The test network is built on BlockEmulator, a framework originally developed by the Huang Lab at Sun Yat-sen University and rigorously validated in peer-reviewed research (e.g., IEEE INFOCOM 2022/2024). For more details, refer to the code repository: [Huang Lab Block Emulator](https://github.com/HuangLab-SYSU/block-emulator).

**The codebase is organized into five layers: storage layer, data layer, network layer, consensus layer, and application layer.**

- The storage folder implements the storage layer using BoltDB for persistent block and data storage.
- The network folder handles data transmission across the network.
- The message folder defines the message types and their logical implementations
- The core folder defines core data structures, including transactions and blocks.
- The blockchain folder manages the blockchain data structure.
- The server and worknode folders contain the operational logic for server nodes and worker nodes, respectively.
- The config folder contains the configuration settings and test data used in the experiment

## Steps to start the test network:

1. Launch a local Docker Golang container, copy the 'getnodetimes' folder into the container, and run the 'main.go' file in the folder to obtain the average mining difficulty time (NodeTimes) for each node group.
2. Replace the NodeTimes in 'config.go' under the 'config' folder with the results obtained in step 1, and then run the command 'docker-compose up -d' in the 'A-PoW' folder to start the test network.

The NodesTimes[][] 2D array records the average mining time for each node in the group at various difficulty levels. Each row stores data for a node group, organized in increasing order of difficulty.

The docker-compose.yaml file defines the startup template for both the Server and Worker nodes. Adjust the number of nodes as required, and ensure that the TotalNodeNum field in the config file is updated to match the new number of nodes.
