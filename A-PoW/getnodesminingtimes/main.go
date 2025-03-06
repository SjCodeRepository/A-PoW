package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

var runiteration = 1

const Nmin = 8
const Nmax = 9
const MaxSNum = 10
const Nnum = 2
const NodeID = "nodetest"

type MiningDifficulty struct {
	DifficultyBase int
	SuccessfulNum  int
}

var Target []byte = []byte{'0', '1'}

func main() {
	sf := 0
	nodetime := make([]float64, Nnum*MaxSNum)
	Length := 0
	times := make([]float64, runiteration)
	TotalMD := make([]MiningDifficulty, Nnum*MaxSNum)
	for count1 := 0; count1 < Nnum; count1++ {
		for count2 := 1; count2 <= MaxSNum; count2++ {
			TotalMD[MaxSNum*count1+count2-1] = MiningDifficulty{
				DifficultyBase: Nmin + count1,
				SuccessfulNum:  count2,
			}
		}
	}
	for index, value := range TotalMD {
		for i := 0; i < runiteration; i++ {
			Block := Block{
				Header: &BlockHeader{
					ParentBlockHash: make([]byte, 10),
					BlockType:       BlockType1,
					Number:          uint64(Length),
					Time:            time.Now(),
					Miner:           NodeID,
					Nounce:          0,
				},
				Body: []*Transaction{
					{
						Sender: "TestNode",
						Nonce:  1,
						Time:   time.Now(),
					},
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
				for i2 := 0; i2 < value.DifficultyBase; i2++ {
					for _, value := range Target {
						if hashcode1[i2] == value {
							count++
						}
					}
				}
				if count == value.DifficultyBase {
					sf++
				}
				//Successful mining
				if sf == value.SuccessfulNum {
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
		nodetime[index] = totaltimes / float64(runiteration)
	}
	fmt.Printf("node mining time is %v\n", nodetime)
}
