package jsonrpc

import (
	"fmt"
	"sync/atomic"

	"github.com/google/uuid"
)

type RPCRequest struct {
	JSONRPC    string      `json:"jsonrpc"`
	ID         string      `json:"id"`
	Method     string      `json:"method"`
	Parameters interface{} `json:"params"`
}

func NewRequest(method string, params interface{}) *RPCRequest {
	request := &RPCRequest{
		Method:     method,
		Parameters: params,
		JSONRPC:    jsonrpcVersion,
		ID:         newID(),
	}
	return request
}

var UseIntegerID = false

var integerID = new(atomic.Uint64)

var useFixedID = false

const defaultFixedID = 1

func newID() string {
	if useFixedID {
		return fmt.Sprint(defaultFixedID)
	}
	if UseIntegerID {
		return fmt.Sprint(integerID.Add(1))
	}
	return uuid.New().String()
}

type RPCResponse[T any] struct {
	JSONRPC string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  T      `json:"result"`
}
