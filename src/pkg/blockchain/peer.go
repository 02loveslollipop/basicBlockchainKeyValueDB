package blockchain
import (
	"bytes"
	"encoding/json"
	"net/http"
	"errors"
)

// Peer abstraction in the blockchain
type Peer struct {
	Address string //Address of the peer
}

func BroadcastBlock(newBlock Block, peers []Peer) error {
    blockData, err := json.Marshal(newBlock)
    if err != nil {
        return err
    }

    for _, peer := range peers {
        resp, err := http.Post("http://"+peer.Address+":8080/block", "application/json", bytes.NewBuffer(blockData)) //TODO: Change the port to a variable
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