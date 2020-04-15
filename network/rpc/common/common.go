package common

import (
	"github.com/harmony-one/go-lib/network/rpc/sharding"
	commonTypes "github.com/harmony-one/go-lib/network/types/common"
	"github.com/harmony-one/go-lib/network/utils"
	goSDK_RPC "github.com/harmony-one/go-sdk/pkg/rpc"
	goSDK_Sharding "github.com/harmony-one/go-sdk/pkg/sharding"
)

// NewRPCClient - resolve the RPC/HTTP Messenger to use for remote commands using a node and a shardID
func NewRPCClient(node string, shardID uint32, shardingStructure []goSDK_Sharding.RPCRoutes, retry *commonTypes.Retry) (*goSDK_RPC.HTTPMessenger, []goSDK_Sharding.RPCRoutes, error) {
	if utils.IsLocalNode(node) {
		return goSDK_RPC.NewHTTPHandler(node), nil, nil
	}

	if shardingStructure == nil || len(shardingStructure) == 0 {
		structure, err := sharding.ShardingStructure(node, retry)
		if err != nil {
			return nil, nil, err
		}
		shardingStructure = structure
	}

	for _, shard := range shardingStructure {
		if uint32(shard.ShardID) == shardID {
			return goSDK_RPC.NewHTTPHandler(shard.HTTP), shardingStructure, nil
		}
	}

	return nil, nil, nil
}
