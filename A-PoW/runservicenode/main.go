package main

import (
	"apow/config"
	"apow/server"
)

func main() {

	//start server node
	server.StartServer(config.NodeTimes)
}
