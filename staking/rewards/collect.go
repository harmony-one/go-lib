package rewards

import (
	"fmt"

	"github.com/harmony-one/go-lib/network"
	"github.com/harmony-one/go-lib/staking"
	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/harmony/accounts"
	"github.com/harmony-one/harmony/accounts/keystore"
	"github.com/harmony-one/harmony/numeric"
	hmyStaking "github.com/harmony-one/harmony/staking/types"
)

// CollectRewards - collects rewards for a given delegator
func CollectRewards(
	keystore *keystore.KeyStore,
	account *accounts.Account,
	rpcClient *rpc.HTTPMessenger,
	chain *common.ChainID,
	delegatorAddress string,
	gasLimit int64,
	gasPrice numeric.Dec,
	nonce uint64,
	keystorePassphrase string,
	node string,
	timeout int,
) (map[string]interface{}, error) {
	payloadGenerator, err := createCollectRewardsTransactionGenerator(delegatorAddress)
	if err != nil {
		return nil, err
	}

	var logMessage string
	if network.Verbose {
		logMessage = fmt.Sprintf("Generating a new collect rewards transaction:\n\tDelegator Address: %s",
			delegatorAddress,
		)
	}

	return staking.SendTx(keystore, account, rpcClient, chain, gasLimit, gasPrice, nonce, keystorePassphrase, node, timeout, payloadGenerator, logMessage)
}

func createCollectRewardsTransactionGenerator(delegatorAddress string) (hmyStaking.StakeMsgFulfiller, error) {
	payloadGenerator := func() (hmyStaking.Directive, interface{}) {
		return hmyStaking.DirectiveCollectRewards, hmyStaking.CollectRewards{
			address.Parse(delegatorAddress),
		}
	}

	return payloadGenerator, nil
}
