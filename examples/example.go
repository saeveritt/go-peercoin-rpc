package main

import (
	"log"

	"github.com/saeveritt/go-peercoin-rpc/config"
	"github.com/saeveritt/go-peercoin-rpc/jsonrpc"
)

func main() {
	// Load the configuration file.
	config, err := config.LoadConfig()
	if err != nil {
		log.Printf("Encountered an Error: %v\n", err)
		return
	}
	// Create a new jsonrpc client.
	client := jsonrpc.NewRPCClient(config)
	// Call the getblockchaininfo method.
	resp, err := client.GetBlockchainInfo()
	if err != nil {
		log.Printf("Encountered an Error: %v\n", err)
		return
	}
	// Print the response.
	log.Printf("Response: %v\n", resp)
	log.Printf("GetChainWork: %v", resp.GetChainwork())

}
