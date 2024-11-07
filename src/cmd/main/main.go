package main
import (
	"net/http"
	//"strconv"
	"github.com/gin-gonic/gin"
	"bufio"
	"os"
	"time"
	"basicBlockchainKeyValueDB/pkg/blockchain"
	"net"
	"fmt"
)

func loadNodes(filename string) error {
	addrss, err := net.InterfaceAddrs() //get all network interfaces
	if err != nil { //Check if there is an error getting the network interfaces
		return err //If there is an error return it
	}

	for _, address := range addrss { //Iterate over the network interfaces
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() { //Check if the address is not the loopback address
			if ipnet.IP.To4() != nil { //Check if the address is an IPv4 address
				selfAddress = ipnet.IP.String() //Set the self address to the address
				break
			}
		}
	}

	if selfAddress == "" { //Check if the self address is empty
		return nil //If it is empty return nil
	}
	
	file, err := os.Open(filename) //Open the file
	if err != nil { //Check if there is an error opening the file
		return err //If there is an error return it
	}
	defer file.Close() //Close the file when the function ends
	scanner := bufio.NewScanner(file) //Create a scanner for the file
	for scanner.Scan() { //Iterate over the lines of the file
		data := scanner.Text() //Get the line
		if data == selfAddress { //Check if the line is the self address
			continue //If it is continue to the iteration
		}
		peers = append(peers, blockchain.Peer{Address: scanner.Text()}) //Append the address to the peers slice
	}

	if err := scanner.Err(); err != nil { //Check if there is an error reading the file
		return err //If there is an error return it
	}
	return nil //If there is no error return nil
}

//Inter node communication

func getBlockchain(c *gin.Context) {
	/*
	GetBlockchain
	Get the current blockchain in the node
	---
	definitions:
		Block:
			type: object
			properties:
				index:
					type: integer
					description: Index of the block in the blockchain
					example: 0
				timestamp:
					type: string
					description: Date of the block
					example: "2018-03-01T00:00:00Z"
				data:
					type: string
					description: Data of the block
					example: "Hello World"
				previousHash:
					type: string
					description: Hash of the previous block
					example: "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
				hash:
					type: string
					description: Hash of the block
					example: "88d4266fd4e6338d13b845fcf289579d209c897823b9217da3e161936f031589"
	responses:
		200:
			description: A list of blocks
			schema:
				type: array
				items:
					$ref: '#/definitions/Block'
		204:
			description: No content
	*/

	if len(blockchainArray) == 0 { //Check if the blockchain is empty
		c.JSON(http.StatusNoContent, nil) //If, then return a 204 status with no content
	}else{
		c.JSON(http.StatusOK, blockchainArray) //Else return the blockchain
	}
}

func addBlock(c *gin.Context) {
	/*
	AddBlock
	Verify and if valid add a new block to the blockchain in the node
	---
	parameters:
		- name: block
		  in: body
		  description: Block object
		  required: true
		  schema:
			$ref: '#/definitions/Block'
	responses:
		201:
			description: Block added
		
		400:
			description: Invalid genesis block
		
		401:
			description: Invalid hash
	*/
	var newBlock blockchain.Block //Create a new block
	err := c.BindJSON(&newBlock) //Bind the JSON to the block
	if err != nil { //Check if there is an error binding the JSON
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid block: Does not match Schema"}) //If there is an error return a 400 status with the error
		return
	}
	if len(blockchainArray) == 0 { //Check if the blockchain is empty
		if newBlock.Index != 0 { //Check if the index of the new block is 0
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genesis block"}) //If it is not return a 400 status with the error
			return
		}
	} else {
		previousBlock := blockchainArray[len(blockchainArray)-1] //Get the last block of the blockchain
		if !blockchain.IsBlockValid(newBlock, previousBlock, blockchain.HashCondition) { //Check if the new block is valid
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid hash"}) //If it is not return a 401 status with the error
			return
		}
	}
	blockchainArray = append(blockchainArray, newBlock) //Append the new block to the blockchain
	c.JSON(http.StatusCreated, newBlock) //Return a 201 status with the new block
}

// Peer to Useer communication

func appendToChain(c *gin.Context) {
	/*
	Append
	Append a new block to the blockchain
	---
	parameters:
		- in: body
		  name: data
		  description: Data to be added to the blockchain
		  required: true
		  type: string
	responses:
		201:
			description: Block added
		400:
			description: Invalid data
		401:
			description: No consensus
		500:
			description: Error calculating hash
	*/
	var data struct {
		Data string `json:"data"` //Create a struct to hold the data
	}
	err := c.BindJSON(&data) //Bind the JSON to the data struct
	if err != nil { //Check if there is an error binding the JSON
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"}) //If there is an error return a 400 status with the error
		return
	}
	previousBlock := blockchainArray[len(blockchainArray)-1] //Get the last block of the blockchain
	newBlock := blockchain.Block{ //Create a new block
		Index:       previousBlock.Index + 1, //Set the index of the new block
		Timestamp:   time.Now(), //Set the timestamp of the new block
		Data:        data.Data, //Set the data of the new block
		PreviosHash: previousBlock.Hash, //Set the previous hash of the new block
	}
	result, err := blockchain.CalculateHash(newBlock, blockchain.HashCondition) //Calculate the hash of the new block
	if err != nil { //Check if there is an error calculating the hash
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error calculating hash"}) //If there is an error return a 500 status with the error
		return
	}
	newBlock.Hash = result //Set the hash of the new block
	//Broadcast the new block to the peers
	err = blockchain.BroadcastBlock(newBlock, peers) //Broadcast the new block to the peers
	if err != nil { //Check if there is an error broadcasting the block
		if err.Error() == "No consensus" { //Check if the error is "No consensus"
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No consensus"}) //If it is return a 401 status with the error
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error broadcasting block"}) //If it is not return a 500 status with the error
		return
	}
	blockchainArray = append(blockchainArray, newBlock) //Append the new block to the blockchain
	c.JSON(http.StatusCreated, "Block added") //Return a 201 status with the message "Block added"

}

func fetchChain(c *gin.Context) {
	/*
	FetchChain
	Return the current blockchain of the node
	---
	parameters:
		- in: query
		  name: node
		  description: Node to fetch the blockchain from
		  required: true
		  type: string
	responses:
		200:
			description: Blockchain
			schema:
				type: array
				items:
					$ref: '#/definitions/Block'
		500:
			description: Internal server error
	*/
	c.JSON(http.StatusOK, blockchainArray) //Return the blockchain
}

// Info endpoints

func ping(c *gin.Context) {
	/*
	Ping
	Return pong
	---
	responses:
		200:
			description: Pong
	*/
	c.String(http.StatusOK, fmt.Sprintf("Pong from %s", selfAddress)) //Return a 200 status with the message "Pong from" and the self address
}

//Router
func setupRouter() *gin.Engine {
	router := gin.Default() //Create a new router
	//Inter node communication endpoints
	router.GET("/blockchain", getBlockchain) //return the blockchain for Inter node communication
	router.POST("/block", addBlock) //Add a block for Inter node communication
	//Peer to User communication endpoints
	router.POST("/append", appendToChain) //Append a block for Peer to User communication
	router.GET("/chain", fetchChain) //Fetch the blockchain for Peer to User communication
	//Info endpoints
	router.GET("/", ping) //Ping the node
	return router //Return the router
}

var blockchainArray []blockchain.Block
var peers []blockchain.Peer
var selfAddress string

func main() {
	path := "nodes.txt" //Set the default path
	if len(os.Args) == 2 { //Check if there is an argument
		path = os.Args[1] //Get the path from the arguments
	}
	err := loadNodes(path) //Load the nodes
	if err != nil { //Check if there is an error loading the nodes
		panic(err) //If there is an error panic
	}
	router := setupRouter() //Setup the router
	router.Run(":8080") //Run the router on port 8080
}