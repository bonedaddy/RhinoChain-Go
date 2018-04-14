package main 

import (
 bolt "github.com/coreos/bbolt"
 "log"
)

type Blockchain struct {
	// here we only store block hash of the tip of the chain
	tip 	[]byte
	db 		*bolt.DB
}

// since we dont want to store all block when iterating over a database
// we will just store the latest one we are examining
type BlockchainIterator struct {
	currentHash []byte
	db 			*bolt.DB
}

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}
	return bci
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block
	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})
	if err != nil {
		log.Fatal("error ", err)
	}
	i.currentHash = block.PrevBlockHash
	return block
}

func (bc *Blockchain) AddBlock(data string) {
	// this will store the last hash
	var lastHash []byte

	// create a read transaction to pull the hash of hte last block
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Fatal("error ", err)
	}
	// generate a new block
	newBlock := NewBlock(data, lastHash)

	// update the bolt database
	err = bc.db.Update(func(tx *bolt.Tx) error {
		// fetch the bucket
		b := tx.Bucket([]byte(blocksBucket))
		// serialize block and store
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Fatal("error ", err)
		}
		// store hash of the new block
		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Fatal("error ", err)
		}
		// set the tip
		bc.tip = newBlock.Hash

		return nil
	})
	if err != nil {
		log.Fatal("error ", err)
	}
}

func NewBlockchain() *Blockchain {
	// here we only want to store the hash of the tip of the chain
	var tip []byte
	// establish a connection to a bolt db file
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal("error ", err)
	}
	// creating a read-write transaction
	err = db.Update(func(tx *bolt.Tx) error {
		// fetch the bucket
		b := tx.Bucket([]byte(blocksBucket))

		// if the bucket doesn't exist, create it
		if b == nil {
			// create genesis block
			genesis := NewGenesisBlock()
			// create the bucket
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Fatal("error ", err)
			}
			// seralize genesis block and store
			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Fatal("error ", err)
			}
			// store the hash of the genesis block (the last block)
			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Fatal("error ", err)
			}
			// set the tip hash
			tip = genesis.Hash
		} else {
			// set the tip hash
			tip = b.Get([]byte("l"))
		}

		return nil
	})
	if err != nil {
		log.Fatal("error ", err)
	}
	bc := Blockchain{tip, db}

	return &bc
}
