package delegation

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

// Delegate - delegate to a validator
func Delegate(
	keystore *keystore.KeyStore,
	account *accounts.Account,
	rpcClient *rpc.HTTPMessenger,
	chain *common.ChainID,
	delegatorAddress string,
	validatorAddress string,
	amount numeric.Dec,
	gasLimit int64,
	gasPrice numeric.Dec,
	nonce uint64,
	keystorePassphrase string,
	node string,
	timeout int,
) (map[string]interface{}, error) {
	payloadGenerator, err := createDelegationTransactionGenerator(delegatorAddress, validatorAddress, amount)
	if err != nil {
		return nil, err
	}

	var logMessage string
	if network.Verbose {
		logMessage = fmt.Sprintf("Generating a new delegation transaction:\n\tDelegator Address: %s\n\tValidator Address: %s\n\tAmount: %f",
			delegatorAddress,
			validatorAddress,
			amount,
		)
	}

	return staking.SendTx(keystore, account, rpcClient, chain, gasLimit, gasPrice, nonce, keystorePassphrase, node, timeout, payloadGenerator, logMessage)
}

func createDelegationTransactionGenerator(delegatorAddress string, validatorAddress string, amount numeric.Dec) (hmyStaking.StakeMsgFulfiller, error) {
	payloadGenerator := func() (hmyStaking.Directive, interface{}) {
		return hmyStaking.DirectiveDelegate, hmyStaking.Delegate{
			address.Parse(delegatorAddress),
			address.Parse(validatorAddress),
			staking.NumericDecToBigIntAmount(amount),
		}
	}

	return payloadGenerator, nil
}
