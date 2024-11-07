package blockchain
import (
	"time" 
	"strconv"
	"errors"
	"crypto/sha256"
	"encoding/hex"
)

// Abstractions

// Abstraction of a block in the blockchain and it json representation
type Block struct {
	Index	 	int `json:"index"` //Index of the block in the blockchain
	Timestamp 	time.Time `json:"timestamp"` //Date of the block
	Data	 	KeyValue `json:"data"` //Data of the block
	PreviosHash	string `json:"previousHash"` //Hash of the previous block
	Hash		string `json:"hash"` //Hash of the block
}

// Block functions

func GenesisBlock() Block {
	//Create the genesis block
	block := Block{
		Index: 0,
		Timestamp: time.Now(),
		Data: KeyValue{Key: "Genesis", Value: "Genesis block"},
		PreviosHash: "",
		Hash: "d754ed9f64ac293b10268157f283ee23256fb32a4f8dedb25c8446ca5bcb0bb3", //Hash of the genesis block (golang)
	}
	return block
}

func GenerateBlock(oldBlock Block, data string, hashCondition func([]byte) bool) (Block, error) {
	//Generate the key value pair
	newData, err := FromString(data)
	var newBlock Block
	t := time.Now()
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t
	newBlock.Data = newData
	newBlock.PreviosHash = oldBlock.Hash
	result, err := CalculateHash(newBlock, hashCondition)
	if err != nil {
		return newBlock, err
	}
	newBlock.Hash = result
	return newBlock, nil
}

func HashCondition(hash []byte) bool {
	//Implement this function
	//the hash condition for this blockchain is that the first 2 bytes of the hash must be 0
	for i := 0; i < 2; i++ {
		if hash[i] != 0 {
			return false
		}
	}
	return true
}

func CalculateHash(block Block, hashCondition func([]byte) bool) (string, error) {
	//Implement this function
	dataAsString, ok := ToString(block.Data)
	if !ok {
		//return an error message
		return "", errors.New("data cannot be parsed as string")
	}
	magicNumber := 0 //This is the magic number that will be used to mine the block
	for {
		payload := strconv.Itoa(block.Index) + block.Timestamp.String() + dataAsString + block.PreviosHash + strconv.Itoa(magicNumber) //Concatenate the data to be hashed
		h:= sha256.New() //Create a new sha256 hash
		h.Write([]byte(payload)) //Write the payload to the hash
		hashed := h.Sum(nil)
		if hashCondition(hashed) {
			return hex.EncodeToString(hashed), nil
		} else 
		{
			magicNumber++
		}
	}
}

func IsBlockValid(newBlock, previousBlock Block, hashCondition func([]byte) bool) bool {
	//Check if the new index is the spected
	if previousBlock.Index+1 != newBlock.Index {
		return false
	}
	//Check if the previous hash is the same as the hash of the previous block
	if previousBlock.Hash != newBlock.PreviosHash {
		return false
	}
	//Check if the hash of the new block is the same as the calculated hash
	result, err := CalculateHash(newBlock, hashCondition)
	if err != nil {
		return false
	}

	if result != newBlock.Hash {
		return false
	}
	return true
}