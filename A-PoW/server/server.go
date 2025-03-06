package server

import (
	"apow/algorithm"
	"apow/config"
	"apow/message"
	"apow/network"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"sort"
	"sync"
	"time"
)

type MiningDifficulty struct {
	MiningDifficultyBase int
	SuccessfulNum        int
}

// Service node
type Server struct {
	CurrentLeader         string
	currentEpisode        int
	NodeLeaderNum         map[string]float64
	NodeIterateNum        map[string]float64
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

	ReferenceTrustValue     []float64
	CurrentRoundTrustValue  [][]float64
	TotalRoundNum           int
	MaliciousBehaviorRounds []int
	LastLeaderTimes         []int
	IDInt                   map[string]int
	DInt                    []map[MiningDifficulty]int

	NetworkNodeNum  int //totalnumber of nodes entering the network successfully
	startExperiment chan int

	mu sync.Mutex
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
func StartServer(Nodetimes [][]float64) {
	S1 := &Server{
		CurrentLeader:           " ",
		currentEpisode:          0,
		NodeLeaderNum:           make(map[string]float64, config.TotalNodeNum),
		NodeIterateNum:          make(map[string]float64, config.TotalNodeNum),
		NodesAddress:            make([]string, config.TotalNodeNum),
		NodesID:                 make([]string, config.TotalNodeNum),
		NodesMiningDifficulty:   make([]MiningDifficulty, config.TotalNodeNum),
		ServerNodeAddress:       os.Getenv("SERVER_IP"),
		Message:                 new(message.Message),
		EndEpisode:              make(chan int),
		JudgeNum:                0,
		IDInt:                   make(map[string]int),
		DInt:                    make([]map[MiningDifficulty]int, config.TotalNodeNum),
		NMDInt:                  make([]int, config.TotalNodeNum),
		Qtable:                  make([][]float64, config.TotalNodeNum),
		ActionsSet:              make([][]MiningDifficulty, config.TotalNodeNum),
		LastLeaderTimes:         make([]int, config.TotalNodeNum),
		CurrentRoundTrustValue:  make([][]float64, config.TotalNodeNum),
		MaliciousBehaviorRounds: make([]int, config.TotalNodeNum),
		HasLeader:               false,
		TotalActions:            make([]MiningDifficulty, config.ActionsNumber),
		startExperiment:         make(chan int),
		NetworkNodeNum:          0,
	}

	for count := 0; count < config.TotalNodeNum; count++ {
		S1.Qtable[count] = make([]float64, config.ActionsNumber)
	}
	S1.NodesID = GetNodesID()
	S1.NodesAddress = GetNodesAdress()
	S1.ReferenceTrustValue = config.TestReferenceTrustValue
	for count1 := 0; count1 < config.Nnum; count1++ {
		for count2 := 1; count2 <= config.MaxSNum; count2++ {
			S1.TotalActions[count2+config.MaxSNum*count1-1] = MiningDifficulty{
				MiningDifficultyBase: config.Nmin + count1,
				SuccessfulNum:        count2,
			}
		}
	}
	rand.Seed(time.Now().UnixNano())
	//	network.InitNetworkTools()

	go S1.ListenForMessage()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	timer1 := time.NewTimer(1 * time.Minute)
OutLooPFirst:
	for {
		select {
		case <-timer1.C:
			fmt.Printf("The network is not start successfully!\n")
			os.Exit(3)
		case <-S1.startExperiment:
			fmt.Printf("The network is start successfully!\n")
			break OutLooPFirst
		case <-ticker.C:
			fmt.Printf("waiting for network start\n")
		}
	}

	//Start official operation
OutLoop3:
	for {
		if S1.currentEpisode == 0 {
			//Run the mining difficulty initialization algorithm, calculate the initial difficulty of each node, and send to each node
			NM := make([]MiningDifficulty, config.TotalNodeNum)
			NM2 := algorithm.QLearning(Nodetimes, config.NodeGroups, config.ActionsNumber, config.Max_Epoch, config.Max_Variance)
			for index, value := range NM2 {
				for i := 0; i < config.NodeNum; i++ {
					S1.NMDInt[config.NodeNum*index+i] = config.MaxSNum*(value.MiningDifficultyBase-config.Nmin) + value.SuccessfulNum
					NM[config.NodeNum*index+i].MiningDifficultyBase = value.MiningDifficultyBase
					NM[config.NodeNum*index+i].SuccessfulNum = value.SuccessfulNum
				}
			}
			S1.NodesMiningDifficulty = NM
			S1.Message.MessageType = message.MessageType1
			S1.InitQLearningSetting()
			InitLeader := S1.NodesID[rand.Intn(config.TotalNodeNum)]
			for index, value := range S1.NodesAddress {
				S1.Message.MessageBody = message.InitDifficultyMessage{
					InitDifficulty:  S1.NodesMiningDifficulty[index].MiningDifficultyBase,
					SuccessfulNum:   S1.NodesMiningDifficulty[index].SuccessfulNum,
					Leader:          InitLeader,
					IsJoinConsensus: S1.ConsensusAdmission(index), //consensusu node admission
				}.EncodeIDMessage()
				network.TcpDial(S1.Message.EncodeMessage(), value)
			}
		} else {
			//Determine whether the round needs to be adjusted. If it does not need to be adjusted, the original difficulty is sent to each node.
			//If it needs to be adjusted, the adjustment algorithm is run
			if (S1.currentEpisode+1)%100 == 0 && S1.JudgeNum <= config.TotalJudgeNum {
				fmt.Printf("Adjust the mining difficulty\n")
				S1.NodesMiningDifficulty = S1.DifficultyAdjustment()
				for index, value := range S1.NodesMiningDifficulty {
					S1.NMDInt[index] = config.MaxSNum*(value.MiningDifficultyBase-config.Nmin) + value.SuccessfulNum
				}
				for count := 0; count < config.TotalJudgeNum; count++ {
					S1.LastLeaderTimes[count] = 0
				}
				S1.JudgeNum++
			}
			S1.Message.MessageType = message.MessageType2
			for index, value := range S1.NodesAddress {
				S1.Message.MessageBody = message.MiningDifficultyMessage{
					DifficultyBase:  S1.NodesMiningDifficulty[index].MiningDifficultyBase,
					SuccessfulNum:   S1.NodesMiningDifficulty[index].SuccessfulNum,
					IsJoinConsensus: S1.ConsensusAdmission(index),
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
				if S1.currentEpisode == config.TotalJudgeNum {
					//final results
					fmt.Printf("The total number of nodes that become leaders is %v\n", S1.NodeLeaderNum)
					fmt.Printf("The total number of iterations for each node is %v\n", S1.NodeIterateNum)
					fmt.Printf("The Reference numberfor each node is %v\n", S1.ReferenceTrustValue)
					break OutLoop3
				}
				break OutLoop2
			default:
				fmt.Printf("Waiting for new episode\n")
			}
		}
	}
}
func (S1 *Server) ListenForMessage() {
	//Receive and process incoming data
	listener, err := net.Listen("tcp", S1.ServerNodeAddress)
	if err != nil {
		log.Fatalf("Listen failed: %v", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		go S1.handleConnection(conn)
	}
}
func (S1 *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	message1 := make([]byte, 10240)
	n, err3 := conn.Read(message1)
	if err3 != nil {
		fmt.Printf("Node %v failed to obtain the data\n", S1.NodesID)
		fmt.Println()
		log.Fatal(err3)
	}
	S1.Message = message.DecodeMessage(message1[:n])
	switch S1.Message.MessageType {
	case message.MessageType4:
		S1.HandleSM(S1.Message.MessageBody)
	case message.MessageType8:
		S1.HandleIN(S1.Message.MessageBody)
	case message.MessageType6:
		S1.HandleEE(S1.Message.MessageBody)
	case message.MessageType10:
		S1.HandleSN(S1.Message.MessageBody)
	}
}

func (S1 *Server) HandleSM(m []byte) {
	Sm := message.DecodeIBMessage(m)
	S1.HasLeader = true
	S1.NodeLeaderNum[Sm.Sender]++
	S1.LastLeaderTimes[S1.IDInt[Sm.Sender]]++
	S1.NodeIterateNum[Sm.Sender] += float64(Sm.Iteration)
}
func (S1 *Server) HandleIN(m []byte) {
	Im := message.DecodeINMessage(m)
	S1.NodeIterateNum[Im.Sender] += Im.IterationNum
}
func (S1 *Server) HandleEE(m []byte) {
	S1.EndEpisode <- 1
}
func (S1 *Server) InitQLearningSetting() {
	for i := 0; i < config.TotalNodeNum; i++ {
		S1.ActionsSet[i] = make([]MiningDifficulty, config.Range*2+1)
	}
	for index, value := range S1.NMDInt {
		if value-config.Range > 0 && value+config.Range <= len(S1.TotalActions) {
			count := 0
			for i := value - config.Range - 1; i <= value+config.Range-1; i++ {
				S1.ActionsSet[index][count] = S1.TotalActions[i]
				S1.DInt[index][S1.TotalActions[i]] = count
				count++
			}
			S1.Actions[index] = config.Range
		}
		if value-config.Range <= 0 {
			count := 0
			for i := 0; i <= config.Range*2+1; i++ {
				S1.ActionsSet[index][count] = S1.TotalActions[i]
				S1.DInt[index][S1.TotalActions[i]] = count
				count++
			}
			S1.Actions[index] = value - 1
		}
		if value+config.Range > len(S1.TotalActions) {
			count := 0
			for i := len(S1.TotalActions) - 1 - config.Range*2; i <= len(S1.TotalActions)-1; i++ {
				S1.ActionsSet[index][count] = S1.TotalActions[i]
				S1.DInt[index][S1.TotalActions[i]] = count
				count++
			}
			S1.Actions[index] = value - (len(S1.TotalActions) - 1 - config.Range*2)
		}
	}
}

// Accept and calculate the node trust value,conduct the consensus node admission
func (S1 *Server) ConsensusNodeAdmission() []bool {
	IsNextRoundParticipating := make([]bool, config.TotalNodeNum)
	for i := 0; i < config.TotalNodeNum; i++ {

		mean := mean(S1.CurrentRoundTrustValue[i])
		sd := standardDeviation(S1.CurrentRoundTrustValue[i])
		TotalTrustValue := 0.0
		for count := 0; count < config.TotalNodeNum; count++ {
			Z := (S1.CurrentRoundTrustValue[i][count] - mean) / sd
			if Z <= 3*sd && Z >= -3*sd {
				TotalTrustValue += S1.CurrentRoundTrustValue[i][count]
			}
		}
		S1.ReferenceTrustValue[i] = TotalTrustValue / config.CalculateTVRoundNum

	}
	// Create an index array to sort Reference Trust Value
	indexes := make([]int, config.TotalNodeNum)
	for i := range indexes {
		indexes[i] = i
	}
	x := make([]float64, config.TotalNodeNum)
	for i := 0; i < config.TotalNodeNum; i++ {
		x[i] = S1.ReferenceTrustValue[i]
	}
	// Sort indexes based on the value of Reference Trust Value (descending)
	sort.Slice(indexes, func(i, j int) bool {
		return x[indexes[i]] > x[indexes[j]]
	})

	// Mark the first 50% of nodes as true and the second 50% as false
	mid := x[config.TotalNodeNum/2]
	for i := 0; i < config.TotalNodeNum; i++ {
		if S1.ReferenceTrustValue[i] < mid {
			IsNextRoundParticipating[indexes[i]] = true
		} else {
			IsNextRoundParticipating[indexes[i]] = false
		}
	}

	return IsNextRoundParticipating
}
func mean(arr []float64) float64 {
	var sum float64
	for _, v := range arr {
		sum += v
	}
	return sum / float64(len(arr))
}

// Calculated standard deviation
func standardDeviation(arr []float64) float64 {
	m := mean(arr)
	var sumSquares float64
	for _, v := range arr {
		sumSquares += math.Pow(v-m, 2)
	}
	return math.Sqrt(sumSquares / float64(len(arr)-1))
}
func (S1 *Server) ConsensusAdmission(node int) bool {
	tr := make([]float64, config.TotalNodeNum)
	for i := 0; i < config.TotalNodeNum; i++ {
		tr[i] = S1.ReferenceTrustValue[i]
	}
	sort.Slice(tr, func(i, j int) bool { return tr[i] > tr[j] })
	if S1.ReferenceTrustValue[node] > tr[config.TotalNodeNum/2] {
		return true
	} else {
		return false
	}
}
func (S1 *Server) HandleSN([]byte) {
	S1.mu.Lock()
	defer S1.mu.Unlock()
	S1.NetworkNodeNum++
	if S1.NetworkNodeNum == config.TotalNodeNum {
		S1.startExperiment <- 1
	}
}
