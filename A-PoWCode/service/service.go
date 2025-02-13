package service

import (
	"A-PoW/apow"
	"A-PoW/message"
	"A-PoW/network"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	saddress      string = "172.28.0.100:5000" //Service node address
	Range         int    = 2                   //The difficulty selection range of mining difficulty adjustment algorithm
	TotalJudgeNum        = 100                 //The total number of iterations of mining difficulty adjustment
	NodeNum              = 30                  //The default number of nodes in the network
	NodeGroups           = 6
	ActionsNumber        = 2
	Max_Epoch            = 2000
	Max_Variance         = 3000
)

type MiningDifficulty struct {
	DifficultyBase int
	SuccessfulNum  int
}

// Service node
type Server struct {
	CurrentLeader         string
	currentEpisode        int
	NodeLeaderNum         map[string]int
	NodeIterateNum        map[string]int
	NodesNum              int
	NodesAddress          []string
	NodesID               []string
	ServerNodeAddress     string
	NodesMiningDifficulty []MiningDifficulty
	NMDInt                []int
	Message               *message.Message
	JudgeNum              int

	EndEpisode chan int
	HasLeader  bool

	Qtable       [][]float64
	TotalActions []MiningDifficulty
	PreActions   []int
	Actions      []int
	ActionsSet   [][]MiningDifficulty

	TrustValue      []float64
	LastLeaderTimes []int
	DInt            map[MiningDifficulty]int
}

func StartServer(Nodetimes [][]float64) {
	S1 := &Server{
		CurrentLeader:         " ",
		currentEpisode:        0,
		NodeLeaderNum:         make(map[string]int, 30),
		NodeIterateNum:        make(map[string]int, 30),
		NodesNum:              0,
		NodesAddress:          make([]string, 30),
		NodesID:               make([]string, 30),
		NodesMiningDifficulty: make([]MiningDifficulty, 30),
		ServerNodeAddress:     saddress,
		Message:               new(message.Message),
		EndEpisode:            make(chan int),
		JudgeNum:              0,
	}

	//Pause for five minutes, waiting for the worker node registration
	timer := time.NewTimer(5 * time.Minute)
	defer timer.Stop()

	Listener, err1 := net.Listen("tcp", S1.ServerNodeAddress)
	if err1 != nil {
		fmt.Printf("Node %v failed to listen on address %v.", S1.NodesID, S1.ServerNodeAddress)
		fmt.Println()
		log.Fatal(err1)
	}
	message1 := make([]byte, 10240)
Outloop1:
	for {
		select {
		case <-timer.C:
			fmt.Println("Wait for the registration period to end")
			break Outloop1
		default:
			conn, err2 := Listener.Accept()
			if err2 != nil {
				fmt.Printf("Node %v failed to obtain the connection\n", S1.NodesID)
				fmt.Println()
				log.Fatal(err2)
			}
			_, err3 := conn.Read(message1)
			if err3 != nil {
				fmt.Printf("Node %v failed to obtain the data\n", S1.NodesID)
				fmt.Println()
				log.Fatal(err3)
			}
			Rm := message.DecodeRMessage(message.DecodeMessage(message1).MessageBody)
			S1.NodesID[S1.NodesNum] = Rm.Id
			S1.NodesAddress[S1.NodesNum] = Rm.Address
			S1.NodesNum++
			conn.Close()
		}
	}
	defer Listener.Close()
	//Start official operation
OutLoop3:
	for {
		if S1.currentEpisode == 0 {
			//Run the mining difficulty initialization algorithm, calculate the initial difficulty of each node, and send to each node
			S1.NodesMiningDifficulty = apow.QLearning(Nodetimes, NodeGroups, ActionsNumber, Max_Epoch, Max_Variance)
			S1.Message.MessageType = message.InitDifficultyMessageType
			S1.InitQLearningSetting()
			for index, value := range S1.NodesAddress {
				S1.Message.MessageBody = message.InitDifficultyMessage{
					InitDifficulty: S1.NodesMiningDifficulty[index].DifficultyBase,
					SuccessfulNum:  S1.NodesMiningDifficulty[index].SuccessfulNum,
				}.EncodeIDMessage()
				network.TcpDial(S1.Message.EncodeMessage(), value)
			}
		} else {
			//Determine whether the round needs to be adjusted. If it does not need to be adjusted, the original difficulty is sent to each node.
			//If it needs to be adjusted, the adjustment algorithm is run
			if (S1.currentEpisode+1)%10 == 0 && S1.JudgeNum <= TotalJudgeNum {
				fmt.Printf("Adjust the mining difficulty\n")
				S1.NodesMiningDifficulty = S1.DifficultyAdjustment()
				S1.JudgeNum++
			}
			S1.Message.MessageType = message.MiningDifficultyMessageType
			for index, value := range S1.NodesAddress {
				S1.Message.MessageBody = message.MiningDifficultyMessage{
					DifficultyBase: S1.NodesMiningDifficulty[index].DifficultyBase,
					SuccessfulNum:  S1.NodesMiningDifficulty[index].SuccessfulNum,
				}.EncodeMMessage()
				network.TcpDial(S1.Message.EncodeMessage(), value)
			}
		}
	OutLoop2:
		for {
			select {
			case <-S1.EndEpisode:
				fmt.Println("End of this Episode")
				S1.currentEpisode++
				if S1.currentEpisode == TotalJudgeNum {
					break OutLoop3
				}
				break OutLoop2
			default:
				conn, err2 := Listener.Accept()
				if err2 != nil {
					fmt.Printf("Node %v failed to obtain the connection\n", S1.NodesID)
					fmt.Println()
					log.Fatal(err2)
				}
				_, err3 := conn.Read(message1)
				if err3 != nil {
					fmt.Printf("Node %v failed to obtain the data\n", S1.NodesID)
					fmt.Println()
					log.Fatal(err3)
				}
				S1.Message = message.DecodeMessage(message1)
				switch S1.Message.MessageType {
				case "SuccessfulMiningMessage":
					S1.HandleSM(S1.Message.MessageBody)
				case "IterationNumMessage":
					S1.HandleIN(S1.Message.MessageBody)
				case "EndEpisode":
					S1.HandleEE(S1.Message.MessageBody)
				}

			}
		}
	}
}
func (S1 *Server) HandleSM(m []byte) {
	Sm := message.DecodeIBMessage(m)
	S1.HasLeader = true
	S1.NodeLeaderNum[Sm.Sender]++
	S1.NodeIterateNum[Sm.Sender]++
}
func (S1 *Server) HandleIN(m []byte) {
	Im := message.DecodeINMessage(m)
	S1.NodeIterateNum[Im.Sender] += Im.IterationNum
}
func (S1 *Server) HandleEE(m []byte) {
	S1.EndEpisode <- 1
}
func (S1 *Server) InitQLearningSetting() {
	for i := 0; i < NodeNum; i++ {
		S1.ActionsSet[i] = make([]MiningDifficulty, Range*2+1)
	}
	for index, value := range S1.NMDInt {
		if value-Range >= 0 && value+Range <= len(S1.TotalActions)-1 {
			count := 0
			for i := value - Range; i <= value+Range; i++ {
				S1.ActionsSet[index][count] = S1.TotalActions[i]
				count++
			}
			S1.Actions[index] = Range
		}
		if value-Range < 0 {
			count := 0
			for i := 0; i <= Range*2+1; i++ {
				S1.ActionsSet[index][count] = S1.TotalActions[i]
				count++
			}
			S1.Actions[index] = value
		}
		if value+Range > len(S1.TotalActions)-1 {
			count := 0
			for i := len(S1.TotalActions) - 1 - Range*2; i <= len(S1.TotalActions)-1; i++ {
				S1.ActionsSet[index][count] = S1.TotalActions[i]
				count++
			}
			S1.Actions[index] = value - (len(S1.TotalActions) - 1 - Range*2)
		}
	}
}
