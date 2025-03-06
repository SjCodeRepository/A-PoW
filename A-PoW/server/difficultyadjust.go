package server

import (
	"apow/config"
	"math"
	"math/rand"
)

func (S1 *Server) DifficultyAdjustment() []MiningDifficulty {
	//computer betas
	TotalTrustValue := 0.0
	for _, value := range S1.ReferenceTrustValue {
		TotalTrustValue += value
	}
	betas := make([]float64, config.TotalNodeNum)
	for count := 0; count < config.TotalNodeNum; count++ {
		if S1.LastLeaderTimes[count] == 0 {
			S1.LastLeaderTimes[count]++
		}
		LNe := float64(config.CRC) * (S1.ReferenceTrustValue[count] / float64(TotalTrustValue))
		betas[count] = math.Log(float64(S1.LastLeaderTimes[count])/LNe) / math.Log(math.E)
	}
	//obtain reward and update Q-table
	for count2 := 0; count2 < config.TotalNodeNum; count2++ {
		r := 0.0
		if math.Abs(math.Abs(betas[count2])) < config.Range2 {
			r = float64(config.R1)
		} else if math.Abs(math.Abs(betas[count2])) < config.Range1 {
			r = 0
		} else {
			r = float64(config.R2)
		}
		_, maxValue := FindMaxQ(S1.Qtable[count2])
		S1.Qtable[count2][S1.Actions[count2]] = S1.Qtable[count2][S1.Actions[count2]] + config.Alpha*(r+config.Beta*maxValue-S1.Qtable[count2][S1.Actions[count2]])
		//	fmt.Printf("%v Q-value : %v\n", Actions[count3], QTable[count3][Actions[count3]])
	}
	//Calculate the judgment conditions and determine whether adjustment is needed
	Betam := 0.0
	for count1 := 0; count1 < 30; count1++ {
		Betam += math.Abs(betas[count1])
	}
	Betam /= 30.0
	//adjust the mining difficulty
	if Betam > config.Max_Range {
		for count3 := 0; count3 < config.TotalNodeNum; count3++ {
			if math.Abs(betas[count3]) < config.Range1 {
				continue
			} else if betas[count3] > 0 {
				S1.Actions[count3] = ChooseAction(S1.Qtable[count3], S1.Actions[count3], 1)
			} else {
				S1.Actions[count3] = ChooseAction(S1.Qtable[count3], S1.Actions[count3], 0)
			}
		}
		for i := 0; i < config.TotalNodeNum; i++ {
			S1.NodesMiningDifficulty[i] = S1.ActionsSet[i][S1.Actions[i]]
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
	if rand.Float64() > config.RandP {
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
