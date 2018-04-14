package main 

import (
	"bytes"
	"time"
	"encoding/gob"
	"log"
)

type Block struct {
	Timestamp		int64
	Data			[]byte
	PrevBlockHash	[]byte
	Hash 			[]byte
	Nonce 			int
}

func NewBlock(data string, PrevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), PrevBlockHash, []byte{}, 0}
	pow :=  NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	
	return block
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Fatal("error ", err)
	}
	return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var block *Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Fatal("error ", err)
	}
	return block
}


func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
