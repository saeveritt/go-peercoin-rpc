package jsonrpc

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/saeveritt/go-peercoin-rpc/config"
)

// RPCRequest struct with json fields
type RPCRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []byte `json:"params,omitempty"`
	Id      int    `json:"id"`
}

// RPCResponse struct with json fields
type RPCResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   *RPCError   `json:"error"`
	Id      int         `json:"id"`
}

// RPCError struct with json fields
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

// RPCNotification struct with json fields
type RPCNotification struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []byte `json:"params,omitempty"`
}

// BatchResponse struct with array of RPCResponse structs
type BatchResponse struct {
	Responses []RPCResponse `json:"responses"`
}

// RPCClient struct with json fields
type RPCClient struct {
	Host            string
	Port            int
	httpClient      *http.Client
	customHeaders   map[string]string
	nextId          int
	autoIncrementId bool
	idMutex         sync.Mutex
}

// NewRPCClient that accepts a Config and returns a RPCClient
func NewRPCClient(config *config.Config) *RPCClient {
	client := &RPCClient{
		Host:            config.Host,
		Port:            config.Port,
		httpClient:      http.DefaultClient,
		customHeaders:   make(map[string]string),
		nextId:          0,
		autoIncrementId: true,
	}
	client.SetBasicAuth(config)
	return client
}

// SetBasicAuth function for RPCClient that accepts a Config and sets a custom Authorization header with the username and password base64 encoded
func (client *RPCClient) SetBasicAuth(config *config.Config) {
	//set authorization header with username and password base64 encoded
	auth := config.Username + ":" + config.Password
	client.customHeaders["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

// NewRPCRequestObject accepts a method and params and returns a RPCRequest
func (client *RPCClient) NewRPCRequestObject(method string, params []byte) *RPCRequest {
	// lock client id mutex
	client.idMutex.Lock()

	rpcRequest := RPCRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		Id:      client.nextId,
	}
	// if client.autoIncrementId is true, increment nextId
	if client.autoIncrementId {
		client.nextId++
	}
	// unlock client id mutex
	client.idMutex.Unlock()

	// set rpcRequest params to empty array if params is nil
	if rpcRequest.Params == nil {
		rpcRequest.Params = []byte{}
	}
	return &rpcRequest
}

// NewRPCNotificationObject that accepts a method and params and returns a RPCNotification struct
func (client *RPCClient) NewRPCNotificationObject(method string, params []byte) *RPCNotification {
	// set rpcNotification params to empty array if params is nil
	if params == nil {
		params = []byte{}
	}
	return &RPCNotification{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
	}
}

// Call function for RPCClient that accepts a method string and params and returns a RPCResponse or error
func (client *RPCClient) Call(method string, params []byte) ([]byte, error) {
	// create http request with method and params
	req, err := client.newRequest(false, method, params)
	if err != nil {
		return nil, err
	}
	// create http response
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// create rpcResponse object from http response body
	rpcResponse := new(RPCResponse)
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&rpcResponse)
	if err != nil {
		// log error and return error
		log.Printf("Encountered an Error: %v\n", err)
		return nil, err
	}
	result, err := json.Marshal(rpcResponse.Result)
	if err != nil {
		// log error and return error
		log.Printf("Encountered an Error: %v\n", err)
		return nil, err
	}

	return result, nil
}

// newRequest function for RPCClient that takes a notification bool, method string, and params and returns a http Request or error
func (client *RPCClient) newRequest(notification bool, method string, params []byte) (*http.Request, error) {
	// create rpcRequest object to be either a RPCRequest or RPCNotification
	var rpcRequest interface{}
	// if notification is true, create rpcNotification object
	if notification {
		rpcNotification := client.NewRPCNotificationObject(method, params)
		rpcRequest = rpcNotification

	} else {
		rpcRequest = client.NewRPCRequestObject(method, params)
	}

	// marshal rpcRequest object to json
	jsonRpcReq, err := json.Marshal(rpcRequest)
	if err != nil {
		// log error and return error
		log.Printf("Error marshalling json: %v", err)
		return nil, err
	}
	// convert client port to string
	// create http request with jsonRpcReq as body
	req, err := http.NewRequest("POST", "http://"+client.Host+":"+strconv.Itoa(client.Port), bytes.NewBuffer(jsonRpcReq))
	if err != nil {
		return nil, err
	}
	// set custom headers
	for k, v := range client.customHeaders {
		req.Header.Set(k, v)
	}
	// log request
	log.Printf("%v", req)

	return req, nil
}
