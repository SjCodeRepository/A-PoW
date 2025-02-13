package main

import (
	"A-PoW/core"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

var runiteration = 20

const (
	NodeID         = "nodetest"
	DifficultyBase = 8
)

var Target []byte = []byte{'0', '1'}

func main() {
	nodetime := 0.0
	Length := 0
	times := make([]float64, runiteration)
	for i := 0; i < runiteration; i++ {
		Block := core.Block{
			Header: &core.BlockHeader{
				ParentBlockHash: make([]byte, 10),
				BlockType:       core.BlockType1,
				Number:          uint64(Length),
				Time:            time.Now(),
				Miner:           NodeID,
				Nounce:          0,
			},
		}
		nounce := 0
		time1 := time.Now()
	OutLoop:
		for {
			nounce++
			Block.Header.Nounce = uint64(nounce)
			h1 := sha256.Sum256(Block.Encode())
			hashcode1 := hex.EncodeToString(h1[:])
			count := 0
			for i2 := 0; i2 < DifficultyBase; i2++ {
				for _, value := range Target {
					if hashcode1[i2] == value {
						count++
					}
				}
			}
			if count == DifficultyBase {
				time2 := time.Now()
				times[i] = time2.Sub(time1).Seconds()
				break OutLoop
			}
		}
	}
	totaltimes := 0.0
	for i := 0; i < runiteration; i++ {
		totaltimes += times[i]
	}
	nodetime = totaltimes / 20.0
	fmt.Printf("node mining time is %v\n", nodetime)
}
