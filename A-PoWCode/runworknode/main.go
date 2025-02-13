package main

import (
	"A-PoW/worknode"
)


func main() {
	NodesID := make([]string, 6)
	NodesAddress := make([]string, 6)

	worknode.StartWorker(NodesID, NodesAddress)
}
