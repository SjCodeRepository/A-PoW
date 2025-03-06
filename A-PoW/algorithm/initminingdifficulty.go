package algorithm

import (
	"apow/config"
	"fmt"
	"math/rand"
)

type MiningDifficulty struct {
	MiningDifficultyBase int
	SuccessfulNum        int
}

// Mining difficulty initialization algorithm Using Q-Learning
func QLearning(NodeTime [][]float64, NodeGroups, ActionsNumber, Max_Epoch int, Max_Variance float64) []MiningDifficulty {
	rm := make([]MiningDifficulty, NodeGroups)
	//Initialize the Q-table, action array, and state array
	QTable := make([][]float64, NodeGroups)
	for i := 0; i < NodeGroups; i++ {
		QTable[i] = make([]float64, ActionsNumber)
	}
	//Set the Q-table to 0
	for i := 0; i < NodeGroups; i++ {
		for j := 0; j < ActionsNumber; j++ {
			QTable[i][j] = 0
		}
	}
	Actions := make([]int, NodeGroups)
	State := make([]float64, NodeGroups)

	//Initialize the Actions and State arrays with random values
	for i := 0; i < NodeGroups; i++ {
		Actions[i] = rand.Intn(ActionsNumber)
	}
	for i := 0; i < NodeGroups; i++ {
		State[i] = NodeTime[i][Actions[i]]
	}

	//Initialize the minimal variance Vmin
	Vmin := Variance(State)

	//Define the reward CrowdReward for the node group with the most clusters, and the reward NormalReward for the other node groups
	CrowdReward := 0
	NormalReward := 0

	//Start the loop
	for count := 0; count < Max_Epoch; count++ {
		//Select Actions and update States
		for i := 0; i < NodeGroups; i++ {
			Actions[i] = ChooseAction(QTable, i, ActionsNumber)
		}
		for i := 0; i < NodeGroups; i++ {
			State[i] = NodeTime[i][Actions[i]]
		}

		//Calculate the variance and perform clustering, preparing for the reward calculation
		varianceNow := Variance(State)
		MaxPNum, nodeSort := KMedoids(State, 3, 100) //Using the K-medoids method to group nodes into clusters, the resulting clusters should be a two-dimensional array

		//Exit condition judgment
		if varianceNow < Max_Variance {
			for x := 0; x < NodeGroups; x++ {
				rm[x].MiningDifficultyBase = Actions[x]/config.MaxSNum + config.Nmin
				rm[x].SuccessfulNum = Actions[x]%config.MaxSNum + 1
			}
			fmt.Printf("Find!%v\n", count)
			for i := 0; i < NodeGroups; i++ {
				fmt.Printf("%v\n", NodeTime[i][Actions[i]])
			}
			return rm
		}

		//Obtain reward
		if varianceNow > Vmin {
			CrowdReward = -1
			NormalReward = -1
		} else {
			if varianceNow == Vmin {
				CrowdReward = 0
				NormalReward = 0
			} else {
				CrowdReward = 10
				NormalReward = 0
			}
		}

		//Update Q-table
		maxValue := 0.0
		maxQValue := make([]float64, NodeGroups)
		for agent := 0; agent < NodeGroups; agent++ {
			maxValue = 0.0
			for i := 0; i < ActionsNumber; i++ {
				if QTable[agent][i] > maxValue {
					maxValue = QTable[agent][i]
				}
			}
			maxQValue[agent] = maxValue
		}
		reward := 0
		for i := 0; i < NodeGroups; i++ {
			if nodeSort[i] == MaxPNum {
				reward = CrowdReward
			} else {
				reward = NormalReward
			}
			QTable[i][Actions[i]] += config.Alpha * (float64(reward) + config.Gamma*maxQValue[i] - QTable[i][Actions[i]])
		}

		Vmin = Variance(State)
	}

	//Returns the action with the largest Q value
	D := make([]int, NodeGroups)
	for i := 0; i < NodeGroups; i++ {
		maxNum := 0
		maxValue := 0.0
		for j := 0; j < ActionsNumber; j++ {
			if QTable[i][j] > maxValue {
				maxNum = j
				maxValue = QTable[i][j]
			}
		}
		D[i] = maxNum
	}

	for i := 0; i < NodeGroups; i++ {
		fmt.Printf("%v\n", NodeTime[i][D[i]])
	}
	for x := 0; x < NodeGroups; x++ {
		rm[x].MiningDifficultyBase = D[x]/config.MaxSNum + config.Nmin
		rm[x].SuccessfulNum = D[x]%config.MaxSNum + 1
	}
	return rm
}
func Variance(data []float64) float64 {
	if len(data) < 2 {
		return 0
	}

	// Calculated mean
	mean := 0.0
	for _, value := range data {
		mean += value
	}
	mean /= float64(len(data))

	// Calculated variance
	variance := 0.0
	for _, value := range data {
		diff := value - mean
		variance += diff * diff
	}
	variance /= float64(len(data) - 1)

	return variance
}
func ChooseAction(qtable [][]float64, agent, ActionsNumber int) int {
	randNumber := rand.Float64()
	maxNum := 0
	maxValue := 0.0
	judge := true
	for count := 0; count < ActionsNumber; count++ {
		if qtable[agent][count] > 0 {
			judge = false
		}
	}
	if judge {
		return rand.Intn(ActionsNumber)
	}
	for i := 0; i < ActionsNumber; i++ {
		if qtable[agent][i] > maxValue {
			maxValue = qtable[agent][i]
			maxNum = i
		}
	}
	if randNumber <= config.Epsilon {
		return rand.Intn(ActionsNumber)
	} else {
		return maxNum
	}
}
