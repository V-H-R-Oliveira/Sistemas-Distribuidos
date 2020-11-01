package model

import (
	"encoding/json"
	"io"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Collection -> Default data collection.
const Collection = "blockchain"

// Serialization -> Interface for data serialization.
type Serialization interface {
	Deserialize(io.ReadCloser) error
	Serialize() ([]byte, error)
}

// Block -> Data structure for a block.
type Block struct {
	ID        string    `json:"id; omitempty"`
	Name      string    `json:"name"`
	Parent    *string   `json:"parent-id; omitempty"`
	Timestamp time.Time `json:"timestamp; omitempty"`
}

// ConcurrentBlock -> Thread-safe Block.
type ConcurrentBlock struct {
	Block *Block
	Mu    *sync.RWMutex
}

// NewBlock -> Creates a new Block.
func NewBlock() *Block {
	return &Block{
		ID:        uuid.New().String(),
		Timestamp: time.Now().UTC(),
		Parent:    nil,
	}
}

// NewConcurrentBlock -> Creates a new concurrent block
func NewConcurrentBlock(block *Block) *ConcurrentBlock {
	return &ConcurrentBlock{
		Block: block,
		Mu:    &sync.RWMutex{},
	}
}

// Deserialize -> Deserialize implementation.
func (b *Block) Deserialize(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(b)
}

// Serialize -> Serialize the payload to a JSON format.
func (b *Block) Serialize() ([]byte, error) {
	return json.Marshal(b)
}
