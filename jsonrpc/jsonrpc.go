package jsonrpc

import "context"

const (
	jsonrpcVersion = "2.0"
)

type JsonRpcClinet interface {
	СallForInto(ctx context.Context, out interface{}, method string, params []interface{}) error
}

type RpcClient struct {
}

func (c *RpcClient) CallForInto(ctx context.Context, out interface{}, method string, params []interface{}) error {

	request := NewRequest(method, params)

	return Call(ctx, request)
}

func Call[T any](ctx context.Context, request *RPCRequest) (*RPCResponse[T], error) {
	return nil, nil
}
