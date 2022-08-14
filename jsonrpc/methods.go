package jsonrpc

import (
	"encoding/json"
	"github.com/saeveritt/go-peercoin-rpc/proto"
	"log"
)

func (client *RPCClient) GetBlockchainInfo() (*proto.GetBlockChainInfoResult, error) {
	resp, err := client.Call("getblockchaininfo", nil)
	if err != nil {
		// log error
		log.Printf("Encountered error: %v", err)
		return nil, err
	}
	result := proto.GetBlockChainInfoResult{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		// log error
		log.Printf("Encountered error while unmarshalling data: %v", err)
		return nil, err
	}
	return &result, err
}
