package validator

import (
	"fmt"

	"github.com/harmony-one/go-lib/crypto"
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

// Create - creates a validator
func Create(
	keystore *keystore.KeyStore,
	account *accounts.Account,
	rpcClient *rpc.HTTPMessenger,
	chain *common.ChainID,
	validatorAddress string,
	description hmyStaking.Description,
	commissionRates hmyStaking.CommissionRates,
	minimumSelfDelegation numeric.Dec,
	maximumTotalDelegation numeric.Dec,
	blsKeys []crypto.BLSKey,
	amount numeric.Dec,
	gasLimit int64,
	gasPrice numeric.Dec,
	nonce uint64,
	keystorePassphrase string,
	node string,
	confirmationWaitTime int,
) (map[string]interface{}, error) {
	payloadGenerator, err := createTransactionGenerator(validatorAddress, description, commissionRates, minimumSelfDelegation, maximumTotalDelegation, blsKeys, amount)
	if err != nil {
		return nil, err
	}

	var logMessage string
	if network.Verbose {
		logMessage = fmt.Sprintf("Generating a new create validator transaction:\n\tValidator Address: %s\n\tValidator Name: %s\n\tValidator Identity: %s\n\tValidator Website: %s\n\tValidator Security Contact: %s\n\tValidator Details: %s\n\tCommission Rate: %f\n\tCommission Max Rate: %f\n\tCommission Max Change Rate: %d\n\tMinimum Self Delegation: %f\n\tMaximum Total Delegation: %f\n\tBls Public Keys: %v\n\tAmount: %f",
			validatorAddress,
			description.Name,
			description.Identity,
			description.Website,
			description.SecurityContact,
			description.Details,
			commissionRates.Rate,
			commissionRates.MaxRate,
			commissionRates.MaxChangeRate,
			minimumSelfDelegation,
			maximumTotalDelegation,
			blsKeys,
			amount,
		)
	}

	return staking.SendTx(keystore, account, rpcClient, chain, gasLimit, gasPrice, nonce, keystorePassphrase, node, confirmationWaitTime, payloadGenerator, logMessage)
}

func createTransactionGenerator(
	validatorAddress string,
	stakingDescription hmyStaking.Description,
	stakingCommissionRates hmyStaking.CommissionRates,
	minimumSelfDelegation numeric.Dec,
	maximumTotalDelegation numeric.Dec,
	blsKeys []crypto.BLSKey,
	amount numeric.Dec,
) (hmyStaking.StakeMsgFulfiller, error) {
	blsPubKeys, blsSigs := staking.ProcessBlsKeys(blsKeys)

	bigAmount := staking.NumericDecToBigIntAmount(amount)
	bigMinimumSelfDelegation := staking.NumericDecToBigIntAmount(minimumSelfDelegation)
	bigMaximumTotalDelegation := staking.NumericDecToBigIntAmount(maximumTotalDelegation)

	payloadGenerator := func() (hmyStaking.Directive, interface{}) {
		return hmyStaking.DirectiveCreateValidator, hmyStaking.CreateValidator{
			ValidatorAddress:   address.Parse(validatorAddress),
			Description:        stakingDescription,
			CommissionRates:    stakingCommissionRates,
			MinSelfDelegation:  bigMinimumSelfDelegation,
			MaxTotalDelegation: bigMaximumTotalDelegation,
			SlotPubKeys:        blsPubKeys,
			SlotKeySigs:        blsSigs,
			Amount:             bigAmount,
		}
	}

	return payloadGenerator, nil
}
