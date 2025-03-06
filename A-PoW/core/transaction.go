package core

//The definition of transaction
import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
	"time"
)

type Transaction struct {
	Sender    string
	Recipient string
	Nonce     uint64
	Signature []byte // not implemented now.
	Value     *big.Int
	TxHash    []byte

	Time time.Time // TimeStamp the tx proposed.
}

func (tx *Transaction) PrintTx() string {
	vals := []interface{}{
		tx.Sender[:],
		tx.Recipient[:],
		tx.Value,
		string(tx.TxHash[:]),
	}
	res := fmt.Sprintf("%v\n", vals)
	return res
}

// Encode transaction for storing
func (tx *Transaction) Encode() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// Decode transaction
func DecodeTx(to_decode []byte) *Transaction {
	var tx Transaction

	decoder := gob.NewDecoder(bytes.NewReader(to_decode))
	err := decoder.Decode(&tx)
	if err != nil {
		log.Panic(err)
	}

	return &tx
}

// new a transaction
func NewTransaction(sender, recipient string, value *big.Int, nonce uint64, proposeTime time.Time) *Transaction {
	tx := &Transaction{
		Sender:    sender,
		Recipient: recipient,
		Value:     value,
		Nonce:     nonce,
		Time:      proposeTime,
	}

	hash := sha256.Sum256(tx.Encode())
	tx.TxHash = hash[:]
	return tx
}
