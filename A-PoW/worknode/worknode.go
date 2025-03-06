package worknode

import (
	"apow/blockchain"
	"apow/config"
	"apow/core"
	"apow/message"
	"apow/network"
	"apow/server"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var SuccessfulNum int = 0

type Worker struct {
	NodeDifficulty *server.MiningDifficulty
	NodesID        []string
	NodesAddress   []string
	BlockChain     *blockchain.BlockChain
	NextLeader     string
	LeaderAddress  string
	CurrentEpisode int

	StartNewEpisode chan int
	StopMining      chan int
	StopEpisode     chan int
	Message         message.Message
	serverAddress   string
	Address         string
	nodeID          string
	isJoinConsensus bool
}

func GetNodesID() []string {
	ni := make([]string, config.TotalNodeNum)
	for i := 1; i <= config.TotalNodeNum; i++ {
		workerEnvVar := fmt.Sprintf("NODE%d_ID", i)
		workerID := os.Getenv(workerEnvVar)
		if workerID == "" {
			fmt.Printf("Node%d_IP is not set!\n", i)
			continue
		}
		ni[i-1] = workerID
	}
	return ni
}
func GetNodesAdress() []string {
	as := make([]string, config.TotalNodeNum)
	for i := 1; i <= config.TotalNodeNum; i++ {
		workerEnvVar := fmt.Sprintf("NODE%d_IP", i)
		workerIP := os.Getenv(workerEnvVar)
		if workerIP == "" {
			fmt.Printf("Node%d_IP is not set!\n", i)
			continue
		}
		as[i-1] = workerIP
	}
	return as
}
func GetID() string {
	nodeID := os.Getenv("CURRENT_NODE_ID")
	if nodeID == "" {
		log.Fatal("NODE_ID not found in environment variables")
	}
	return nodeID
}
func GetAddress() string {
	nodeIP := os.Getenv("CURRENT_NODE_IP")
	if nodeIP == "" {
		log.Fatal("NODE_ID not found in environment variables")
	}
	return nodeIP
}
func StartWorker() {
	//Initialize the Worker class
	W1 := &Worker{
		NodeDifficulty:  new(server.MiningDifficulty),
		CurrentEpisode:  0,
		serverAddress:   os.Getenv("SERVER_IP"),
		NodesID:         GetNodesID(),
		NodesAddress:    GetNodesAdress(),
		StartNewEpisode: make(chan int, 10),
		StopMining:      make(chan int),
		StopEpisode:     make(chan int),
		Address:         GetAddress(),
		nodeID:          GetID(),
		isJoinConsensus: true,
	}
	W1.BlockChain = blockchain.NewBlcokchain(W1.nodeID)

	//Start Listening
	go W1.ListenForMessage()

	network.TcpDial(message.Message{
		MessageType: message.MessageType10,
		MessageBody: message.StartNetworkMessage{
			Sender: W1.nodeID,
		}.EncodeSMMessage(),
	}.EncodeMessage(), W1.serverAddress)

OutLoop:
	for {
		select {
		case <-W1.StopEpisode:
			if W1.CurrentEpisode == config.EndCodeNum { //End Episode
				break OutLoop
			}
			W1.CurrentEpisode++
		case <-W1.StartNewEpisode: //Start Episode
			if !W1.isJoinConsensus { //Consensus node admission
				continue
			}
			if W1.nodeID == W1.NextLeader { //This node is the leader
				timer := time.NewTimer(5 * time.Minute)
			OutLoop1:
				for {
					select {
					case <-timer.C:
						Em := &message.EndEpisodeMessage{
							Sender: W1.nodeID,
						}
						network.Broadcast(W1.Address, W1.NodesAddress, message.Message{
							MessageBody: Em.EncodeEEMessage(),
							MessageType: message.MessageType6,
						}.EncodeMessage())
						network.TcpDial(message.Message{
							MessageBody: Em.EncodeEEMessage(),
							MessageType: message.MessageType6,
						}.EncodeMessage(), W1.serverAddress)
						if W1.CurrentEpisode == config.EndCodeNum {
							break OutLoop
						}
						W1.CurrentEpisode++
						break OutLoop1
					default:
						fmt.Printf("Sending TxBlock\n")
					}
				}
			} else {
				Block := core.Block{
					Header: &core.BlockHeader{
						BlockType:       core.BlockType1,
						ParentBlockHash: W1.BlockChain.CurrentBlock.Hash,
						Number:          uint64(W1.BlockChain.Length + 1),
						Time:            time.Now(),
						Miner:           W1.nodeID,
						Nounce:          0,
					},
				}
				nounce := 0.0
			Outloop2:
				for {
					select {
					case <-W1.StopMining:
						m := message.Message{
							MessageType: message.MessageType8,
							MessageBody: message.IterationNumMessage{
								Sender:         W1.nodeID,
								IterationNum:   nounce,
								CurrentEpisode: W1.CurrentEpisode,
							}.EncodeINMessage(),
						}
						network.TcpDial(m.EncodeMessage(), W1.serverAddress)
						break Outloop2
					default:
						//mining
						nounce++
						Block.Header.Nounce = uint64(nounce)
						h1 := sha256.Sum256(Block.Encode())
						hashcode1 := hex.EncodeToString(h1[:])
						count := 0
						for i2 := 0; i2 < W1.NodeDifficulty.MiningDifficultyBase; i2++ {
							for _, value := range config.Target {
								if hashcode1[i2] == value {
									count++
								}
							}
						}
						if count == W1.NodeDifficulty.MiningDifficultyBase {
							SuccessfulNum++
						}
						if SuccessfulNum == W1.NodeDifficulty.SuccessfulNum {
							m := message.Message{
								MessageType: message.MessageType3,
								MessageBody: message.IdentityBlockMessage{
									Sender:        W1.nodeID,
									IdentityBlock: Block.Encode(),
									Iteration:     nounce,
									LeaderEpisode: W1.CurrentEpisode,
									LeaderAddress: W1.Address,
								}.EncodeIBMessage(),
							}
							network.TcpDial(m.EncodeMessage(), W1.LeaderAddress)
							SuccessfulNum = 0
							break Outloop2
						}
						fmt.Printf("Mining...\n")

					}
				}
			}
		default:
			fmt.Printf("witing for new Epiaode\n")
		}
	}
}
func (W1 *Worker) ListenForMessage() {

	Listener, err1 := net.Listen("tcp", W1.Address)
	if err1 != nil {
		fmt.Printf("Node %v failed to listen at address %v\n", W1.NodesID, W1.Address)
		fmt.Println()
		log.Fatal(err1)
	}
	for {
		//Receive and process incoming data
		conn, err2 := Listener.Accept()
		if err2 != nil {
			fmt.Printf("Node %v failed to obtain a connection\n", W1.NodesID)
			fmt.Println()
			log.Fatal(err2)
		}

		go W1.handleConnection(conn)
	}
}
func (W1 Worker) handleConnection(conn net.Conn) {
	defer conn.Close()

	message1 := make([]byte, 102400)
	n, err3 := conn.Read(message1)
	if err3 != nil {
		fmt.Printf("Node %v failed to get data\n", W1.NodesID)
		fmt.Println()
		log.Fatal(err3)
	}
	m := message.DecodeMessage(message1[:n])
	switch m.MessageType {
	case message.MessageType3:
		W1.HandleBP(m.MessageBody)
	case message.MessageType4:
		W1.HandleIB(m.MessageBody)
	case message.MessageType6:
		W1.HandleEE(m.MessageBody)
	case message.MessageType2:
		W1.HandleMD(m.MessageBody)
	case message.MessageType1:
		W1.HandleID(m.MessageBody)
	}
}
func (W1 *Worker) HandleIB(body []byte) {
	m := message.DecodeBAMessage(body)
	W1.BlockChain.AddBlock(core.DecodeB(m.IdentityBlock))
	W1.BlockChain.CurrentBlock = core.DecodeB(m.IdentityBlock)
	W1.BlockChain.Length++
	W1.NextLeader = m.Sender
	W1.LeaderAddress = m.LeaderAddress

	W1.StopMining <- 1
}
func (W1 *Worker) HandleEE(body []byte) {
	W1.StopEpisode <- 1
}
func (W1 *Worker) HandleMD(body []byte) {
	m := message.DecodeMMessage(body)
	W1.NodeDifficulty = &server.MiningDifficulty{
		MiningDifficultyBase: m.DifficultyBase,
		SuccessfulNum:        m.SuccessfulNum,
	}
	W1.isJoinConsensus = m.IsJoinConsensus
	W1.StartNewEpisode <- 1
}
func (W1 *Worker) HandleID(body []byte) {
	m := message.DecodeIDMessage(body)
	W1.NodeDifficulty = &server.MiningDifficulty{
		MiningDifficultyBase: m.InitDifficulty,
		SuccessfulNum:        m.SuccessfulNum,
	}
	W1.NextLeader = m.Leader
	W1.isJoinConsensus = m.IsJoinConsensus
	W1.StartNewEpisode <- 1
}
func (W1 *Worker) HandleBP(body []byte) {
	//Test whether the identity block is legitimate
	m := message.DecodeIBMessage(body)
	W1.NextLeader = m.Sender
	W1.LeaderAddress = m.LeaderAddress
	network.Broadcast(W1.nodeID, W1.NodesAddress, message.Message{
		MessageType: message.MessageType4,
		MessageBody: message.BlockValidationAckMessage{
			Sender:           W1.nodeID,
			IdentityBlock:    m.IdentityBlock,
			ValidationResult: true,
			Iteration:        m.Iteration,
		}.EncodeBAMessage(),
	}.EncodeMessage())
}
