package nonces

import (
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/go-sdk/pkg/transaction"
)

// CurrentNonce - get a specific nonce from input or from the network
func CurrentNonce(rpcClient *rpc.HTTPMessenger, address string) uint64 {
	return transaction.GetNextNonce(address, rpcClient)
}
