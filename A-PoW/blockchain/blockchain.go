package blockchain

import (
	"apow/core"
	"apow/storage"
	"sync"
	"time"
)

type BlockChain struct {
	CurrentBlock *core.Block
	Storage      *storage.Storage
	Length       int
	Pmlock       sync.RWMutex
}

func (chain *BlockChain) AddBlock(b *core.Block) {
	chain.Pmlock.Lock()
	defer chain.Pmlock.Unlock()
	chain.Storage.AddBlock(b)
}
func (chain *BlockChain) AddGenisisBlock() {
	GenisisBlock := &core.Block{
		Header: &core.BlockHeader{
			BlockType: core.BlockType1,
			Number:    0,
			Time:      time.Now(),
		},
	}
	chain.CurrentBlock = GenisisBlock
	chain.Storage.AddBlock(GenisisBlock)
}
func NewBlcokchain(NodeID string) *BlockChain {
	chain := &BlockChain{
		Storage: storage.NewStorage(NodeID),
		Length:  0,
		Pmlock:  sync.RWMutex{},
	}
	chain.AddGenisisBlock()
	return chain
}
