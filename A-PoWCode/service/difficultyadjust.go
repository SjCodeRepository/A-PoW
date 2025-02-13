package service

//Adaptive mining difficulty adjustment algorithm
import (
	"math"
	"math/rand"
)

const Alpha float64 = 0.1
const Beta float64 = 0.9
const r1 int = 5
const r2 int = -1
const Range1 = 1.0
const Range2 = 0.8
const Max_Range = 0.5
const ActionNum int = 3
const CRC int = 200
const MaxEpisode = 1000
const RandP float64 = 0.8

type Node struct {
	NodeHashRate   float64
	NodeTrustValue float64
}
type Difficult struct {
	Difficult  int
	SuccessNum int
}

func (S1 *Server) DifficultyAdjustment() []MiningDifficulty {
	//computer betas
	TotalTrustValue := 0.0
	for _, value := range S1.TrustValue {
		TotalTrustValue += value
	}
	betas := make([]float64, 30)
	for count := 0; count < NodeNum; count++ {
		if S1.LastLeaderTimes[count] == 0 {
			S1.LastLeaderTimes[count]++
		}
		LNe := float64(CRC) * (S1.TrustValue[count] / float64(TotalTrustValue))
		betas[count] = math.Log(float64(S1.LastLeaderTimes[count])/LNe) / math.Log(math.E)
	}
	//obtain reward and update Q-table
	for count2 := 0; count2 < NodeNum; count2++ {
		r := 0.0
		if math.Abs(math.Abs(betas[count2])) < Range2 {
			r = float64(r1)
		} else if math.Abs(math.Abs(betas[count2])) < Range1 {
			r = 0
		} else {
			r = float64(r2)
		}
		_, maxValue := FindMaxQ(S1.Qtable[count2])
		S1.Qtable[count2][S1.DInt[S1.NodesMiningDifficulty[count2]]] = S1.Qtable[count2][S1.DInt[S1.NodesMiningDifficulty[count2]]] + Alpha*(r+Beta*maxValue-S1.Qtable[count2][S1.DInt[S1.NodesMiningDifficulty[count2]]])
		//	fmt.Printf("%v Q-value : %v\n", Actions[count3], QTable[count3][Actions[count3]])
	}
	//Calculate the judgment conditions and determine whether adjustment is needed
	Betam := 0.0
	for count1 := 0; count1 < 30; count1++ {
		Betam += math.Abs(betas[count1])
	}
	Betam /= 30.0
	//adjust the mining difficulty
	if Betam > Max_Range {
		for count3 := 0; count3 < NodeNum; count3++ {
			if math.Abs(betas[count3]) < Range1 {
				continue
			} else if betas[count3] > 0 {
				S1.Actions[count3] = ChooseAction(S1.Qtable[count3], S1.Actions[count3], 1)
			} else {
				S1.Actions[count3] = ChooseAction(S1.Qtable[count3], S1.Actions[count3], 0)
			}
		}
		return S1.NodesMiningDifficulty
	} else {
		return S1.NodesMiningDifficulty
	}
}
func FindMaxQ(Q []float64) (int, float64) {
	maxInt := 0
	maxValue := Q[0]

	for i := 1; i < len(Q); i++ {
		if Q[i] > maxValue {
			maxInt = i
			maxValue = Q[i]
		}
	}
	if maxValue <= 0 {
		maxValue = 0
	}
	return maxInt, maxValue
}
func ChooseAction(QTable []float64, Action int, flag int) int {
	if rand.Float64() > RandP {
		return rand.Intn(len(QTable))
	} else {
		if flag == 0 {
			if Action == 0 || Action == 1 {
				return 0
			}
			judge := true
			for _, value := range QTable[:Action-1] {
				if value != 0 {
					judge = false
				}
			}
			if judge {
				return rand.Intn(Action)
			} else {
				r, _ := FindMaxQ(QTable[:Action-1])
				return r
			}
		} else {
			if Action == len(QTable)-2 || Action == len(QTable)-1 {
				return len(QTable) - 1
			}
			judge := true
			for _, value := range QTable[Action+1:] {
				if value != 0 {
					judge = false
				}
			}
			if judge {
				return rand.Intn(len(QTable)-Action-2) + Action + 1
			} else {
				r, _ := FindMaxQ(QTable[Action+1:])
				return r + Action + 1
			}
		}
	}
}
