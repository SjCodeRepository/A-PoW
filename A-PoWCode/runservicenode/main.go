package main

import (
	"A-PoW/service"
)

var NodeTimes [][]float64

func main() {
	service.StartServer(NodeTimes)
}
