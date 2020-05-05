package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/harmony-one/go-lib/utils"
	goSdkRPC "github.com/harmony-one/go-sdk/pkg/rpc"
)

// BlockWrapper - wrapper for the GetBlockByNumber RPC method
type BlockWrapper struct {
	ID      string    `json:"id" yaml:"id"`
	JSONRPC string    `json:"jsonrpc" yaml:"jsonrpc"`
	Result  BlockInfo `json:"result" yaml:"result"`
	Error   RPCError  `json:"error,omitempty" yaml:"error,omitempty"`
}

// BlockInfo - block info
type BlockInfo struct {
	BlockNumber  uint64    `json:"-" yaml:"-"`
	Difficulty   int       `json:"difficulty,omitempty" yaml:"difficulty,omitempty"`
	ExtraData    string    `json:"extraData,omitempty" yaml:"extraData,omitempty"`
	Hash         string    `json:"hash,omitempty" yaml:"hash,omitempty"`
	Nonce        uint32    `json:"nonce,omitempty" yaml:"nonce,omitempty"`
	RawTimestamp string    `json:"timestamp,omitempty" yaml:"timestamp,omitempty"`
	Timestamp    time.Time `json:"-" yaml:"-"`
}

// RPCGenericSingleHexResponse - wrapper for RPC calls returning a single result in a hex format
type RPCGenericSingleHexResponse struct {
	ID      string `json:"id" yaml:"id"`
	JSONRPC string `json:"jsonrpc" yaml:"jsonrpc"`
	Result  string `json:"result" yaml:"result"`
}

// GetBlockByNumber - retrieve information for a specific block
func GetBlockByNumber(blockNumber uint64, includeTransactions bool, node string) (BlockInfo, error) {
	response := BlockWrapper{}
	result := BlockInfo{}
	blockNum := fmt.Sprintf("0x%x", blockNumber)

	bytes, err := goSdkRPC.RawRequest(goSdkRPC.Method.GetBlockByNumber, node, []interface{}{blockNum, includeTransactions})
	if err != nil {
		return result, err
	}

	json.Unmarshal(bytes, &response)

	if response.Error.Message != "" {
		return result, fmt.Errorf("%s (%d)", response.Error.Message, response.Error.Code)
	}

	result = response.Result
	result.BlockNumber = blockNumber
	if err = result.Initialize(); err != nil {
		return result, err
	}

	return result, nil
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

// Initialize - initialize and convert values for a given BlockInfo struct
func (blockInfo *BlockInfo) Initialize() error {
	if blockInfo.RawTimestamp != "" {
		unixTime, err := utils.HexToDecimal(blockInfo.RawTimestamp)
		if err != nil {
			return err
		}

		blockInfo.Timestamp = time.Unix(int64(unixTime), 0).UTC()
	}

	return nil
}
