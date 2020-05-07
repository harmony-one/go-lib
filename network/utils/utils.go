package utils

import (
	"fmt"
	"net"
	"strings"

	"github.com/harmony-one/go-sdk/pkg/common"
)

// IsLocalNode - check what node to use
func IsLocalNode(node string) bool {
	removePrefix := strings.TrimPrefix(node, "http://")
	removePrefix = strings.TrimPrefix(removePrefix, "https://")
	possibleIP := strings.Split(removePrefix, ":")[0]
	isLocalNode := (possibleIP == "localhost" || net.ParseIP(string(possibleIP)) != nil)
	return isLocalNode
}

// IdentifyNetworkChainID - identifies a chain id given a network name
func IdentifyNetworkChainID(network string) (chain *common.ChainID, err error) {
	network = NormalizedNetworkName(network)

	if network == "dryrun" {
		return &common.Chain.MainNet, nil
	}

	return common.StringToChainID(network)
}

// GenerateNodeAddress - generates a node address given a network, mode and a shardID
func GenerateNodeAddress(network string, mode string, shardID uint32) (node string) {
	node = ToNodeAddress(network, shardID)

	if strings.ToLower(mode) == "local" {
		node = "http://localhost:9500"
	}

	return node
}

// NormalizedNetworkName - return a normalized network name
func NormalizedNetworkName(network string) string {
	switch strings.ToLower(network) {
	case "local", "localnet":
		return "localnet"
	case "dev", "devnet", "pga":
		return "devnet"
	case "staking", "openstaking", "os", "ostn", "pangaea":
		return "pangaea"
	case "partner", "partnernet", "pstn":
		return "partner"
	case "stress", "stresstest", "stn", "stressnet":
		return "stressnet"
	case "testnet", "p", "b":
		return "testnet"
	case "dryrun", "dry":
		return "dryrun"
	case "mainnet", "main", "t":
		return "mainnet"
	default:
		return ""
	}
}

// ToNodeAddress - generates a node address based on a given network name
func ToNodeAddress(network string, shardID uint32) (node string) {
	network = NormalizedNetworkName(network)

	switch network {
	case "localnet":
		node = fmt.Sprintf("http://localhost:950%d", shardID)
	case "devnet":
		node = fmt.Sprintf("https://api.s%d.pga.hmny.io", shardID)
	case "pangaea":
		node = fmt.Sprintf("https://api.s%d.os.hmny.io", shardID)
	case "partner":
		node = fmt.Sprintf("https://api.s%d.ps.hmny.io", shardID)
	case "stressnet":
		node = fmt.Sprintf("https://api.s%d.stn.hmny.io", shardID)
	case "testnet":
		node = fmt.Sprintf("https://api.s%d.b.hmny.io", shardID)
	case "dryrun":
		node = fmt.Sprintf("https://api.s%d.dry.hmny.io", shardID)
	case "mainnet":
		node = fmt.Sprintf("https://api.s%d.t.hmny.io", shardID)
	default:
		node = fmt.Sprintf("http://localhost:950%d", shardID)
	}

	return node
}

// ResolveStartingNode - resolve the starting node to use for resolving the sharding structure
func ResolveStartingNode(network string, mode string, shardID uint32, nodes []string) string {
	node := ""
	switch mode {
	case "api":
		node = GenerateNodeAddress(network, mode, shardID)
	case "custom":
		if len(nodes) > 0 {
			node = nodes[0]
		} else {
			node = GenerateNodeAddress(network, mode, shardID)
		}
	default:
		node = GenerateNodeAddress(network, mode, shardID)
	}

	return node
}
