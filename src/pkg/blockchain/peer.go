package blockchain
import (
	"bytes"
	"encoding/json"
	"net/http"
	"errors"
)

type BroadcastBlockRequest struct {
    Block Block `json:"block"`
    MagicNumber int `json:"magicNumber"`
}

// Peer abstraction in the blockchain
type Peer struct {
	Address string //Address of the peer
}

func BroadcastBlock(newBlock Block,magicNumber int, peers []Peer) error {
    payload := BroadcastBlockRequest{
        Block: newBlock,
        MagicNumber: magicNumber,
    }
    blockData, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    for _, peer := range peers {
        
        resp, err := http.Post("http://"+peer.Address+":8080/block", "application/json", bytes.NewBuffer(blockData))
        if err != nil {
            return err
        }
		//check if the response status code is not 201
		if resp.StatusCode != http.StatusCreated {
			return errors.New("no consensus")
		}

    }

    return nil
}