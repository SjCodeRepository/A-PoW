package worknode

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	blockchain "A-PoW/blockchain"
	"A-PoW/core"
	"A-PoW/message"
	"A-PoW/network"
	"A-PoW/service"
)

const (
	ServerAddress string = "172.28.0.100:5000"
	EndCodeNum    int    = 1000
	TotalJudgeNum int    = 100
)

var Target []byte = []byte{'0', '1'}

type Worker struct {
	NodeDifficulty *service.MiningDifficulty
	NodesID        []string
	NodesAddress   []string
	BlockChain     *blockchain.BlockChain
	NextLeader     string
	CurrentEpisode int

	StartNewEpisode chan int
	StopMining      chan int
	StopEpisode     chan int
	Message         message.Message
	Address         string
	nodeID          string
}

func StartWorker(nodesID, nodesAddress []string) {

	W1 := &Worker{
		NodeDifficulty:  new(service.MiningDifficulty),
		NodesID:         nodesID,
		NodesAddress:    nodesAddress,
		CurrentEpisode:  0,
		StartNewEpisode: make(chan int),
		StopMining:      make(chan int),
		StopEpisode:     make(chan int),
	}
	W1.nodeID = os.Getenv("NODE_ID")
	W1.BlockChain = blockchain.NewBlcokchain(W1.nodeID)
	ip, err := getIPAddress()
	W1.Address = ip
	if err != nil {
		fmt.Printf("Error getting IP: %v\n", err)
		ip = "unknown"
	}

	Rm := message.RegisterMessage{
		Address: ip,
		Id:      os.Getenv("NODE_ID"),
	}
	network.TcpDial(Rm.EncodeRMessage(), ServerAddress)
	go W1.ListenForMessage()
OutLoop:
	for {
		select {
		case <-W1.StartNewEpisode:
			if W1.nodeID == W1.NextLeader {
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
							MessageType: message.EndEpisodeMessageType,
						}.EncodeMessage())
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
				nounce := 0
			Outloop2:
				for {
					select {
					case <-W1.StopMining:
						break Outloop2
					default:
						//mining
						nounce++
						Block.Header.Nounce = uint64(nounce)
						h1 := sha256.Sum256(Block.Encode())
						hashcode1 := hex.EncodeToString(h1[:])
						count := 0
						for i2 := 0; i2 < W1.NodeDifficulty.DifficultyBase; i2++ {
							for _, value := range Target {
								if hashcode1[i2] == value {
									count++
								}
							}
						}
						if count == W1.NodeDifficulty.DifficultyBase {
							m := message.Message{
								MessageType: message.IdentityBlockMessageType,
								MessageBody: message.IdentityBlockMessage{
									Sender:        W1.nodeID,
									IdentityBlock: Block.Encode(),
									Iteration:     nounce,
								}.EncodeIBMessage(),
							}
							network.Broadcast(W1.Address, W1.NodesAddress, m.EncodeMessage())
						}
						fmt.Printf("执行挖矿")

					}
				}
			}
		default:
			fmt.Printf("witing for new Epiaode\n")
		}
		W1.CurrentEpisode++
		if W1.CurrentEpisode == TotalJudgeNum {
			break OutLoop
		}
	}
}
func (W1 *Worker) ListenForMessage() {

	Listener, err1 := net.Listen("tcp", W1.Address)
	if err1 != nil {
		fmt.Printf("节点%v在地址%v监听失败", W1.NodesID, W1.Address)
		fmt.Println()
		log.Fatal(err1)
	}
	message1 := make([]byte, 10240)
	for {
		conn, err2 := Listener.Accept()
		if err2 != nil {
			fmt.Printf("节点%v获取连接失败", W1.NodesID)
			fmt.Println()
			log.Fatal(err2)
		}
		_, err3 := conn.Read(message1)
		if err3 != nil {
			fmt.Printf("节点%v读取数据失败", W1.NodesID)
			fmt.Println()
			log.Fatal(err3)
		}
		m := message.DecodeMessage(message1)
		switch m.MessageType {
		case message.IdentityBlockMessageType:
			W1.HandleIB(m.MessageBody)
		case message.EndEpisodeMessageType:
			W1.HandleEE(m.MessageBody)
		case message.MiningDifficultyMessageType:
			W1.HandleMD(m.MessageBody)
		case message.InitDifficultyMessageType:
			W1.HandleID(m.MessageBody)
		}
	}
}
func (W1 *Worker) HandleIB(body []byte) {
	m := message.DecodeIBMessage(message.DecodeMessage(body).MessageBody)
	W1.BlockChain.Pmlock.Lock()
	W1.BlockChain.AddBlock(core.DecodeB(m.IdentityBlock))
	W1.BlockChain.CurrentBlock = core.DecodeB(m.IdentityBlock)
	W1.BlockChain.Length++
	W1.BlockChain.Pmlock.Unlock()
	W1.NextLeader = m.Sender

	W1.StopMining <- 1
}
func (W1 *Worker) HandleEE(body []byte) {
	W1.StopEpisode <- 1
}
func (W1 *Worker) HandleMD(body []byte) {
	m := message.DecodeMMessage(message.DecodeMessage(body).MessageBody)
	W1.NodeDifficulty = &service.MiningDifficulty{
		DifficultyBase: m.DifficultyBase,
		SuccessfulNum:  m.SuccessfulNum,
	}
	W1.StartNewEpisode <- 1
}
func (W1 *Worker) HandleID(body []byte) {
	m := message.DecodeIDMessage(message.DecodeMessage(body).MessageBody)
	W1.NodeDifficulty = &service.MiningDifficulty{
		DifficultyBase: m.InitDifficulty,
		SuccessfulNum:  m.SuccessfulNum,
	}
	W1.NextLeader = m.Leader
	W1.StartNewEpisode <- 1
}
func getIPAddress() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		if iface.Name == "eth0" {
			addrs, err := iface.Addrs()
			if err != nil {
				return "", err
			}

			for _, addr := range addrs {
				ipNet, ok := addr.(*net.IPNet)
				if ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						return ipNet.IP.String(), nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("no IP address found")
}
