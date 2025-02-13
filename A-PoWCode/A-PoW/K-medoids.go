package apow

import (
	"math"
	"math/rand"
	"sort"
	//"fmt"
)

// KMedoids 实现K-medoids算法
func KMedoids(data []float64, k int, maxIterations int) (int, []int) {
	// 初始化medoids
	NodeInt := make(map[float64]int, len(data))
	for index, value := range data {
		NodeInt[value] = index
	}
	medoids := make([]float64, k)
	//	copy(medoids, data[:k])
	randMedoids := make([]int, 3)
	for i := 0; i < 3; i++ {
		if i == 0 {
			randMedoids[i] = rand.Intn(len(data))
		}
		if i == 1 {
			randMedoids[i] = rand.Intn(len(data))
			for {
				if randMedoids[1] == randMedoids[0] {
					randMedoids[1] = rand.Intn(len(data))
				} else {
					break
				}
			}
		}
		if i == 2 {
			randMedoids[i] = rand.Intn(len(data))
			for {
				if randMedoids[2] == randMedoids[0] || randMedoids[2] == randMedoids[1] {
					randMedoids[2] = rand.Intn(len(data))
				} else {
					break
				}
			}
		}
	}
	for i := 0; i < 3; i++ {
		medoids[i] = data[randMedoids[i]]
	}
	clusters := make(map[float64][]float64)

	// 迭代优化
	for iteration := 0; iteration < maxIterations; iteration++ {
		// 分配数据点到最近的medoid
		clusters = make(map[float64][]float64)
		for _, point := range data {
			nearestMedoid := findNearestMedoid(point, medoids)
			clusters[nearestMedoid] = append(clusters[nearestMedoid], point)
		}

		// 更新medoids
		for i := 0; i < k; i++ {
			bestMedoid := medoids[i]
			bestCost := calculateCost(clusters[bestMedoid], bestMedoid)

			for _, candidate := range clusters[medoids[i]] {
				candidateCost := calculateCost(clusters[medoids[i]], candidate)
				if candidateCost < bestCost {
					bestCost = candidateCost
					bestMedoid = candidate
				}
			}

			medoids[i] = bestMedoid
		}
	}
	nodeSort := make([]int, len(data))
	sort.Float64s(medoids)
	for index, value1 := range medoids {
		for _, value2 := range clusters[value1] {
			nodeSort[NodeInt[value2]] = index
		}
	}
	lenP := make([]int, 3)
	for index := 0; index < len(data); index++ {
		lenP[nodeSort[index]]++
	}
	MaxPNum := 0
	for i := 0; i < 3; i++ {
		if lenP[i] > lenP[MaxPNum] {
			MaxPNum = i
		}
	}
	return MaxPNum, nodeSort
}

// findNearestMedoid 找到最近的medoid
func findNearestMedoid(point float64, medoids []float64) float64 {
	nearestMedoid := medoids[0]
	minDistance := distance(point, nearestMedoid)

	for _, medoid := range medoids[1:] {
		d := distance(point, medoid)
		if d < minDistance {
			minDistance = d
			nearestMedoid = medoid
		}
	}

	return nearestMedoid
}

// distance 计算两个点之间的欧氏距离
func distance(p1, p2 float64) float64 {
	return math.Sqrt(math.Pow(p1-p2, 2))
}

// calculateCost 计算一个簇内所有点到medoid的距离之和
func calculateCost(cluster []float64, medoid float64) float64 {
	cost := 0.0
	for _, point := range cluster {
		cost += distance(point, medoid)
	}
	return cost
}
