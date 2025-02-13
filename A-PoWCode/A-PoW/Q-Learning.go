package apow

import (
	"fmt"
	"math/rand"

	"A-PoW/service"
)

const (
	Nnum    = 2
	Nmin    = 8
	Mmax    = 9
	MaxSNum = 10  // 行为数量
	alpha   = 0.1 // 学习率
	gamma   = 0.9 // 折扣因子
	epsilon = 0.1 // 探索因子
	epochs  = 5   // 迭代次数
)

// 使用Q-Learning的挖矿难度初始化调整方法
func QLearning(NodeTime [][]float64, NodeGroups, ActionsNumber, Max_Epoch int, Max_Variance float64) []service.MiningDifficulty {
	rm := make([]service.MiningDifficulty, NodeGroups)
	//初始化Q表、行为数组以及状态数组
	QTable := make([][]float64, NodeGroups)
	for i := 0; i < NodeGroups; i++ {
		QTable[i] = make([]float64, ActionsNumber)
	}
	//将Q表置0
	for i := 0; i < NodeGroups; i++ {
		for j := 0; j < ActionsNumber; j++ {
			QTable[i][j] = 0
		}
	}
	Actions := make([]int, NodeGroups)
	State := make([]float64, NodeGroups)

	//初始为Actions和State数组赋予随机值
	for i := 0; i < NodeGroups; i++ {
		Actions[i] = rand.Intn(ActionsNumber)
	}
	for i := 0; i < NodeGroups; i++ {
		State[i] = NodeTime[i][Actions[i]]
	}

	//初始化最小方差
	Vmin := Variance(State)

	//定义给最多簇中节点组的奖励CrowdReward，以及给其他节点组的奖励NormalReward
	CrowdReward := 0
	NormalReward := 0

	//开始循环
	for count := 0; count < Max_Epoch; count++ {
		//选择行为并更新State
		for i := 0; i < NodeGroups; i++ {
			Actions[i] = ChooseAction(QTable, i, ActionsNumber)
		}
		for i := 0; i < NodeGroups; i++ {
			State[i] = NodeTime[i][Actions[i]]
		}

		//计算方差并分簇，做好计算奖励前的准备工作
		varianceNow := Variance(State)
		MaxPNum, nodeSort := KMedoids(State, 3, 100) //使用K-medoids方法将节点组分簇，得到的clusters应该是一个二维数组

		//退出条件判断
		if varianceNow < Max_Variance {
			for x := 0; x < NodeGroups; x++ {
				rm[x].DifficultyBase = Actions[x]/MaxSNum + Nmin
				rm[x].SuccessfulNum = Actions[x]%MaxSNum + 1
			}
			fmt.Printf("Find!%v\n", count)
			for i := 0; i < NodeGroups; i++ {
				fmt.Printf("%v\n", NodeTime[i][Actions[i]])
			}
			return rm
		}

		//获取奖励
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

		//更新Q表
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
			QTable[i][Actions[i]] += alpha * (float64(reward) + gamma*maxQValue[i] - QTable[i][Actions[i]])
		}

		Vmin = Variance(State)
	}

	//返回Q值最大的行为
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
		rm[x].DifficultyBase = D[x]/MaxSNum + Nmin
		rm[x].SuccessfulNum = D[x]%MaxSNum + 1
	}
	return rm
}
func Variance(data []float64) float64 {
	if len(data) < 2 {
		return 0
	}

	// 计算均值
	mean := 0.0
	for _, value := range data {
		mean += value
	}
	mean /= float64(len(data))

	// 计算方差
	variance := 0.0
	for _, value := range data {
		diff := value - mean
		variance += diff * diff
	}
	variance /= float64(len(data) - 1) // 使用 n-1 作为分母，即样本方差

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
	if randNumber <= epsilon {
		return rand.Intn(ActionsNumber)
	} else {
		return maxNum
	}
}
