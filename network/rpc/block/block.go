package block

import (
	"github.com/harmony-one/go-sdk/pkg/rpc"
)

// GetCurrentEpoch - returns the block header current epoch
func GetCurrentEpoch(node string) (uint32, error) {
	params := []interface{}{}
	blockReply, err := rpc.Request(rpc.Method.GetLatestBlockHeader, node, params)
	if err != nil {
		return 0, err
	}
	return uint32(blockReply["result"].(map[string]interface{})["epoch"].(float64)), nil
}
