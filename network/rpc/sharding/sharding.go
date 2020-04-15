package sharding

import (
	"time"

	commonTypes "github.com/harmony-one/go-lib/network/types/common"
	"github.com/harmony-one/go-sdk/pkg/sharding"
	"github.com/pkg/errors"
)

// ShardingStructure - retrieve the sharding structure based on a given node
func ShardingStructure(node string, retry *commonTypes.Retry) (routes []sharding.RPCRoutes, err error) {
	if retry != nil && retry.Attempts > 0 {
		ret := *retry
		for {
			ret.Attempts--

			routes, err = shardingStructure(node)
			if err == nil {
				return routes, nil
			}

			if ret.Attempts > 0 {
				time.Sleep(time.Second * time.Duration(ret.Wait))
			} else {
				break
			}
		}
	} else {
		routes, err = shardingStructure(node)
	}

	return routes, err
}

func shardingStructure(node string) (routes []sharding.RPCRoutes, err error) {
	routes, err = sharding.Structure(node)
	if err != nil {
		return nil, errors.Wrapf(err, "network.ShardingStructure")
	}
	return routes, err
}
