package message

//The definition of message
import (
	"encoding/json"
	"log"
)

// The messages include the following:
// Service node side: 1. Initial mining difficulty message (InitDifficulty), 2. Mining difficulty message (MiningDifficulty), 3. End program message (EndCode)
// Worker node side: 1. End of current round Episode message (EndEpisode), 2. Transaction block message (Txblock), 3. Identity block message (IdentityBlock), 4. Iteration count message (Iteration), 5. Registration message (Register)

const (
	InitDifficultyMessageType   = "InitDifficulty"
	RegisterMessageType         = "Register"
	MiningDifficultyMessageType = "MiningDifficulty"
	EndCodeMessageType          = "EndCode"
	EndEpisodeMessageType       = "EndEpisode"
	TxBlockMessageType          = "TxBlock"
	IdentityBlockMessageType    = "IdentityBlock"
	IterationNumMessageType     = "Iteration"
)

type Message struct {
	MessageBody []byte
	MessageType string
}
type InitDifficultyMessage struct {
	InitDifficulty int
	SuccessfulNum  int
	Leader         string
}
type RegisterMessage struct {
	Address string
	Id      string
}
type MiningDifficultyMessage struct {
	DifficultyBase int
	SuccessfulNum  int
}
type EndCodeMessage struct {
	Sender string
}
type EndEpisodeMessage struct {
	Sender string
}
type TxBlockMessage struct {
	Sender  string
	TxBlock []byte
}
type IdentityBlockMessage struct {
	Sender        string
	LeaderEpisode int
	IdentityBlock []byte
	Iteration     int
}
type IterationNumMessage struct {
	Sender         string
	IterationNum   int
	CurrentEpisode int
}

// message
func (m Message) EncodeMessage() []byte {
	msg, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}
func DecodeMessage(data []byte) *Message {
	msg := &Message{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}

// Initial mining difficulty message
func (i *InitDifficultyMessage) IDMessageType() string {
	return InitDifficultyMessageType
}
func (m InitDifficultyMessage) EncodeIDMessage() []byte {
	msg, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}
func DecodeIDMessage(data []byte) *InitDifficultyMessage {
	msg := &InitDifficultyMessage{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}

// Registration message
func (i *RegisterMessage) RMessageType() string {
	return RegisterMessageType
}
func (m *RegisterMessage) EncodeRMessage() []byte {
	msg, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}
func DecodeRMessage(data []byte) *RegisterMessage {
	msg := &RegisterMessage{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}

// mining difficulty message
func (i *MiningDifficultyMessage) MMessageType() string {
	return MiningDifficultyMessageType
}
func (m MiningDifficultyMessage) EncodeMMessage() []byte {
	msg, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}
func DecodeMMessage(data []byte) *MiningDifficultyMessage {
	msg := &MiningDifficultyMessage{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}

// End program message
func (i *EndCodeMessage) ECMessageType() string {
	return EndCodeMessageType
}
func (m *EndCodeMessage) EncodeECMessage() []byte {
	msg, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}
func DecodeECMessage(data []byte) *EndCodeMessage {
	msg := &EndCodeMessage{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}

// End of current round Episode message
func (i *EndEpisodeMessage) EEMessageType() string {
	return EndCodeMessageType
}
func (m *EndEpisodeMessage) EncodeEEMessage() []byte {
	msg, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}
func DecodeEEMessage(data []byte) *EndEpisodeMessage {
	msg := &EndEpisodeMessage{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}

// Transaction block message
func (i *TxBlockMessage) TMessageType() string {
	return TxBlockMessageType
}
func (m *TxBlockMessage) EncodeTMessage() []byte {
	msg, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}
func DecodeTMessage(data []byte) *TxBlockMessage {
	msg := &TxBlockMessage{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}

// Identity block message
func (i *IdentityBlockMessage) IBMessageType() string {
	return IdentityBlockMessageType
}
func (m IdentityBlockMessage) EncodeIBMessage() []byte {
	msg, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}
func DecodeIBMessage(data []byte) *IdentityBlockMessage {
	msg := &IdentityBlockMessage{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}

// Iteration count message
func (i *IterationNumMessage) IterationNumMessageType() string {
	return IterationNumMessageType
}
func (m *IterationNumMessage) EncodeINMessage() []byte {
	msg, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}
func DecodeINMessage(data []byte) *IterationNumMessage {
	msg := &IterationNumMessage{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}
