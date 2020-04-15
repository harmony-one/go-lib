package network

import (
	"fmt"
	"sync"

	"github.com/harmony-one/go-lib/network/rpc/balances"
	commonRPC "github.com/harmony-one/go-lib/network/rpc/common"
	"github.com/harmony-one/go-lib/network/rpc/nonces"
	commonTypes "github.com/harmony-one/go-lib/network/types/common"
	"github.com/harmony-one/go-lib/network/utils"
	goSDK_common "github.com/harmony-one/go-sdk/pkg/common"
	goSDK_RPC "github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/go-sdk/pkg/sharding"
	goSDK_sharding "github.com/harmony-one/go-sdk/pkg/sharding"
	"github.com/harmony-one/harmony/numeric"
)

// Network - represents a network configuration
type Network struct {
	Name              string
	Mode              string
	Node              string // Node - override any other node settings, if Node is set it will be used as the node address everywhere
	ChainID           *goSDK_common.ChainID
	Retry             commonTypes.Retry
	Shards            map[uint32]Shard
	ShardCount        int
	ShardingStructure []goSDK_sharding.RPCRoutes
	Mutex             sync.Mutex
}

// Shard - represents a shard configuration
type Shard struct {
	Node      string
	RPCClient *goSDK_RPC.HTTPMessenger
}

// Initialize - initializes a given network
func (network *Network) Initialize() {
	network.ShardCount = len(network.ShardingStructure)
	network.SetChainID()
}

// SetChainID - sets the chain id for a given network
func (network *Network) SetChainID() {
	chainID, err := utils.IdentifyNetworkChainID(network.Name)
	if chainID == nil || err != nil {
		network.ChainID = &goSDK_common.Chain.TestNet
	} else {
		network.ChainID = chainID
	}
}

// ShardsToMap - convert the shards to a map[uint32]string
func (network *Network) ShardsToMap() (shards map[uint32]string) {
	shards = make(map[uint32]string)

	for shardID, shard := range network.Shards {
		shards[shardID] = shard.Node
	}

	return shards
}

// NodeAddress - generates a node address given the network's name and mode + the supplied shardID
func (network *Network) NodeAddress(shardID uint32) string {
	if network.Node != "" {
		return network.Node
	}

	shard, ok := network.Shards[shardID]
	if ok && shard.Node != "" {
		return shard.Node
	}

	generated := utils.GenerateNodeAddress(network.Name, network.Mode, shardID)
	network.Shards[shardID] = Shard{Node: generated}

	return network.Shards[shardID].Node
}

// IdentifyChainID - identifies a chain id given a network name
func (network *Network) IdentifyChainID() (chain *goSDK_common.ChainID, err error) {
	return utils.IdentifyNetworkChainID(network.Name)
}

// GetAllShardBalances - checks the balances in all shards for a given network, mode and address
func (network *Network) GetAllShardBalances(address string) (map[uint32]numeric.Dec, error) {
	return balances.GetAllShardBalances(address, network.ShardsToMap(), &network.Retry)
}

// GetShardBalance - gets the balance for a given network, mode, address and shard
func (network *Network) GetShardBalance(address string, shardID uint32) (numeric.Dec, error) {
	return balances.GetShardBalance(address, shardID, network.ShardsToMap(), &network.Retry)
}

// GetTotalBalance - gets the total balance across all shards for a given network, mode and address
func (network *Network) GetTotalBalance(address string) (numeric.Dec, error) {
	return balances.GetTotalBalance(address, network.ShardsToMap(), &network.Retry)
}

// CurrentNonce - gets the current nonce for a given network, mode and address
func (network *Network) CurrentNonce(address string, shardID uint32) uint64 {
	rpcClient, err := network.RPCClient(shardID)
	if err != nil {
		return 0
	}
	return nonces.CurrentNonce(rpcClient, address)
}

// SetShardingStructure - sets the sharding structure for the network
func (network *Network) SetShardingStructure(shardingStructure []goSDK_sharding.RPCRoutes) {
	if network.ShardingStructure == nil && len(network.ShardingStructure) == 0 && len(shardingStructure) > 0 {
		network.ShardingStructure = shardingStructure
	}
}

// RPCClient - resolve the RPC/HTTP Messenger to use for remote commands
func (network *Network) RPCClient(shardID uint32) (*goSDK_RPC.HTTPMessenger, error) {
	if len(network.Node) > 0 {
		client, shardingStructure, err := commonRPC.NewRPCClient(network.Node, shardID, network.ShardingStructure, &network.Retry)
		network.SetShardingStructure(shardingStructure)
		return client, err
	}

	shard, ok := network.Shards[shardID]
	if ok && shard.RPCClient != nil {
		return shard.RPCClient, nil
	}

	node := network.NodeAddress(shardID)
	client, shardingStructure, err := commonRPC.NewRPCClient(node, shardID, network.ShardingStructure, &network.Retry)
	network.SetShardingStructure(shardingStructure)
	network.Mutex.Lock()
	network.Shards[shardID] = generateShardType(node, shardID, shardingStructure)
	network.Mutex.Unlock()

	return client, err
}

// GenerateShardSetup - generate the shard setup based on a given node
func GenerateShardSetup(node string, network string, mode string, nodes []string) (shards map[uint32]Shard, shardingStructure []goSDK_sharding.RPCRoutes, err error) {
	shards = make(map[uint32]Shard)

	shardingStructure, err = sharding.Structure(node)
	if err != nil {
		return shards, shardingStructure, fmt.Errorf("can't connect to the %s network using node %s - make sure the network is online and that you haven't gotten rate-limited", network, node)
	}

	if mode == "api" {
		for i := 0; i < len(shardingStructure); i++ {
			shardNode := utils.GenerateNodeAddress(network, mode, uint32(i))
			shards[uint32(i)] = generateShardType(shardNode, uint32(i), shardingStructure)
		}
	} else {
		if len(nodes) == len(shardingStructure) {
			for i, customNode := range nodes {
				shards[uint32(i)] = generateShardType(customNode, uint32(i), shardingStructure)
			}
		} else {
			return shards, shardingStructure, fmt.Errorf("the node count for the nodes you've specified (%d) doens't match the expected node count (%d) for the network %s", len(nodes), len(shardingStructure), network, node)
		}
	}

	return shards, shardingStructure, nil
}

func generateShardType(node string, shardID uint32, shardingStructure []goSDK_sharding.RPCRoutes) Shard {
	shard := Shard{Node: node}
	rpcClient, _, err := commonRPC.NewRPCClient(node, shardID, shardingStructure, nil)
	if err == nil {
		shard.RPCClient = rpcClient
	}

	return shard
}
