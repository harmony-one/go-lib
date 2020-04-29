package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	goSdkRPC "github.com/harmony-one/go-sdk/pkg/rpc"
)

// RPCGenericSingleHexResponse - wrapper for RPC calls returning a single result in a hex format
type RPCGenericSingleHexResponse struct {
	ID      string `json:"id" yaml:"id"`
	JSONRPC string `json:"jsonrpc" yaml:"jsonrpc"`
	Result  string `json:"result" yaml:"result"`
}

// GetCurrentBlockNumber - get the current block number for a given node
func GetCurrentBlockNumber(node string) (uint64, error) {
	response := RPCGenericSingleHexResponse{}

	bytes, err := goSdkRPC.RawRequest(goSdkRPC.Method.BlockNumber, node, []interface{}{})
	if err != nil {
		return 0, err
	}

	json.Unmarshal(bytes, &response)

	return handleBlockResult(response)
}

// GetTransactionCountByBlockNumber - get the transaction count for a given block
func GetTransactionCountByBlockNumber(blockNumber uint64, node string) (uint64, error) {
	response := RPCGenericSingleHexResponse{}

	param := fmt.Sprintf("0x%x", blockNumber)

	bytes, err := goSdkRPC.RawRequest(goSdkRPC.Method.GetBlockTransactionCountByNumber, node, []interface{}{param})
	if err != nil {
		return 0, err
	}

	json.Unmarshal(bytes, &response)

	return handleBlockResult(response)
}

func handleBlockResult(response RPCGenericSingleHexResponse) (uint64, error) {
	if response.Result == "" {
		return uint64(0), errors.New("Empty result")
	}
	res, _ := big.NewInt(0).SetString(response.Result[2:], 16)
	return res.Uint64(), nil
}
