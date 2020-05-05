package rpc

// RPCError - error returned from the RPC endpoint
type RPCError struct {
	Code    int    `json:"code,omitempty" yaml:"code,omitempty"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
}
