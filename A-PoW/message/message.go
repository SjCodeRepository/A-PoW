package message

//The definition of message
import (
	"encoding/json"
	"log"
)

// The messages include the following:
// MiningDifficultyUpdate, IdentityBlockProposal, BlockValidationAck, EpisodeTermination, and TrustParameterSubmission
const (
	MessageType1  = "InitMiningDifficultyUpdate"
	MessageType2  = "MiningDifficultyUpdate"
	MessageType3  = "IdentityBlockProposal"
	MessageType4  = "BlockValidationAck"
	MessageType5  = "TxBlockUpdate"
	MessageType6  = "EpisodeTermination"
	MessageType7  = "TrustParameterSubmission"
	MessageType8  = "IterationSubmission"
	MessageType9  = "CodeTermination"
	MessageType10 = "StartNetwork"
)

type Message struct {
	MessageBody []byte
	MessageType string
}

// InitMiningDifficultyUpdate
type InitDifficultyMessage struct {
	InitDifficulty  int
	SuccessfulNum   int
	Leader          string
	IsJoinConsensus bool
}

// MiningDifficultyUpdate
type MiningDifficultyMessage struct {
	DifficultyBase  int
	SuccessfulNum   int
	IsJoinConsensus bool
}

// CodeTermination
type EndCodeMessage struct {
	Sender string
}

// EpisodeTermination
type EndEpisodeMessage struct {
	Sender string
}

// TxBlockUpdate
type TxBlockMessage struct {
	Sender  string
	TxBlock []byte
}

// IdentityBlockProposal
type IdentityBlockMessage struct {
	Sender        string
	LeaderEpisode int
	IdentityBlock []byte
	Iteration     float64
	LeaderAddress string
}

// IterationSubmission
type IterationNumMessage struct {
	Sender         string
	IterationNum   float64
	CurrentEpisode int
}

// BlockValidationAck
type BlockValidationAckMessage struct {
	Sender           string
	ValidationResult bool
	IdentityBlock    []byte
	Iteration        float64
	LeaderAddress    string
}

// TrustParameterSubmission
type TrustParameterSubmissionMessage struct {
	Character      string //include "Leader" and "Regular node"
	TrustParameter []byte
	Sender         string
}
type TrustParameterType1 struct {
	//Data passed by the Leader
	Legal_trans_num     []float64
	Malicious_block_num []float64
}
type TrustParameterType2 struct {
	//Data passed by the Regular node
	Valid_trans_num   float64
	Invalid_block_num float64
}

// Start network message
type StartNetworkMessage struct {
	Sender string
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

// InitMiningDifficultyUpdate
func (i *InitDifficultyMessage) IDMessageType() string {
	return MessageType1
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

// MiningDifficultyUpdate
func (i *MiningDifficultyMessage) MMessageType() string {
	return MessageType2
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

// CodeTermination
func (i *EndCodeMessage) ECMessageType() string {
	return MessageType9
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

// EpisodeTermination
func (i *EndEpisodeMessage) EEMessageType() string {
	return MessageType6
}
func (m EndEpisodeMessage) EncodeEEMessage() []byte {
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

// TxBlockUpdate
func (i *TxBlockMessage) TMessageType() string {
	return MessageType5
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

// IdentityBlockProposal
func (i *IdentityBlockMessage) IBMessageType() string {
	return MessageType3
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

// IterationSubmission
func (i *IterationNumMessage) IterationNumMessageType() string {
	return MessageType8
}
func (m IterationNumMessage) EncodeINMessage() []byte {
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

// BlockValidationAck
func (m BlockValidationAckMessage) BlockValidationAckMType() string {
	return MessageType4
}
func (m BlockValidationAckMessage) EncodeBAMessage() []byte {
	msg, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}
func DecodeBAMessage(data []byte) *BlockValidationAckMessage {
	msg := &BlockValidationAckMessage{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}

// TrustParameterSubmissionMessage
func (m TrustParameterSubmissionMessage) TrustParameterMType() string {
	return MessageType7
}
func (m TrustParameterSubmissionMessage) EncodeTPMessage() []byte {
	msg, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}
func DecodeTPMessage(data []byte) *TrustParameterSubmissionMessage {
	msg := &TrustParameterSubmissionMessage{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}
func (m TrustParameterType1) EncodeTPMessage() []byte {
	msg, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}
func DecodeT1Message(data []byte) *TrustParameterType1 {
	msg := &TrustParameterType1{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}
func (m TrustParameterType2) EncodeT2Message() []byte {
	msg, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}
func DecodeT2Message(data []byte) *TrustParameterType2 {
	msg := &TrustParameterType2{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}

// Start Network
func (m StartNetworkMessage) GetSMType() string {
	return MessageType10
}
func (m StartNetworkMessage) EncodeSMMessage() []byte {
	msg, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}
func DecodeSMessage(data []byte) *StartNetworkMessage {
	msg := &StartNetworkMessage{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}
