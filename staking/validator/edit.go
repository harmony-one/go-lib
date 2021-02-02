package validator

import (
	"fmt"
	"strings"

	"github.com/harmony-one/go-lib/crypto"
	"github.com/harmony-one/go-lib/network"
	"github.com/harmony-one/go-lib/staking"
	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/harmony/accounts"
	"github.com/harmony-one/harmony/accounts/keystore"
	"github.com/harmony-one/harmony/crypto/bls"
	"github.com/harmony-one/harmony/numeric"
	"github.com/harmony-one/harmony/staking/effective"
	hmyStaking "github.com/harmony-one/harmony/staking/types"
)

// Edit - edits the details for an existing validator
func Edit(
	keystore *keystore.KeyStore,
	account *accounts.Account,
	rpcClient *rpc.HTTPMessenger,
	chain *common.ChainID,
	validatorAddress string,
	description hmyStaking.Description,
	commissionRate *numeric.Dec,
	minimumSelfDelegation numeric.Dec,
	maximumTotalDelegation numeric.Dec,
	blsKeyToRemove *crypto.BLSKey,
	blsKeyToAdd *crypto.BLSKey,
	status string,
	gasLimit int64,
	gasPrice numeric.Dec,
	nonce uint64,
	keystorePassphrase string,
	node string,
	timeout int,
) (map[string]interface{}, error) {
	statusEnum := determineEposStatus(status)

	payloadGenerator, err := editTransactionGenerator(validatorAddress, description, commissionRate, minimumSelfDelegation, maximumTotalDelegation, blsKeyToRemove, blsKeyToAdd, statusEnum)
	if err != nil {
		return nil, err
	}

	var logMessage string
	if network.Verbose {
		logMessage = fmt.Sprintf("Generating a new edit validator transaction:\n\tValidator Address: %s\n\tValidator Name: %s\n\tValidator Identity: %s\n\tValidator Website: %s\n\tValidator Security Contact: %s\n\tValidator Details: %s\n\tCommission Rate: %v\n\tMinimum Self Delegation: %f\n\tMaximum Total Delegation: %f\n\tRemove BLS key: %v\n\tAdd BLS key: %v\n\tStatus: %v",
			validatorAddress,
			description.Name,
			description.Identity,
			description.Website,
			description.SecurityContact,
			description.Details,
			commissionRate,
			minimumSelfDelegation,
			maximumTotalDelegation,
			blsKeyToRemove,
			blsKeyToAdd,
			statusEnum,
		)
	}

	return staking.SendTx(keystore, account, rpcClient, chain, gasLimit, gasPrice, nonce, keystorePassphrase, node, timeout, payloadGenerator, logMessage)
}

func determineEposStatus(status string) (statusEnum effective.Eligibility) {
	switch strings.ToLower(status) {
	case "active":
		return effective.Active
	case "inactive":
		return effective.Inactive
	default:
		return effective.Nil
	}
}

func editTransactionGenerator(
	validatorAddress string,
	stakingDescription hmyStaking.Description,
	commissionRate *numeric.Dec,
	minimumSelfDelegation numeric.Dec,
	maximumTotalDelegation numeric.Dec,
	blsKeyToRemove *crypto.BLSKey,
	blsKeyToAdd *crypto.BLSKey,
	statusEnum effective.Eligibility,
) (hmyStaking.StakeMsgFulfiller, error) {
	var shardBlsKeyToRemove *bls.SerializedPublicKey
	if blsKeyToRemove != nil {
		shardBlsKeyToRemove = blsKeyToRemove.ShardPublicKey
	}

	var shardBlsKeyToAdd *bls.SerializedPublicKey
	var shardBlsKeyToAddSig *bls.SerializedSignature
	if blsKeyToAdd != nil {
		shardBlsKeyToAdd = blsKeyToAdd.ShardPublicKey
		shardBlsKeyToAddSig = blsKeyToAdd.ShardSignature
	}

	bigMinimumSelfDelegation := staking.NumericDecToBigIntAmount(minimumSelfDelegation)
	bigMaximumTotalDelegation := staking.NumericDecToBigIntAmount(maximumTotalDelegation)

	payloadGenerator := func() (hmyStaking.Directive, interface{}) {
		return hmyStaking.DirectiveEditValidator, hmyStaking.EditValidator{
			ValidatorAddress:   address.Parse(validatorAddress),
			Description:        stakingDescription,
			CommissionRate:     commissionRate,
			MinSelfDelegation:  bigMinimumSelfDelegation,
			MaxTotalDelegation: bigMaximumTotalDelegation,
			SlotKeyToRemove:    shardBlsKeyToRemove,
			SlotKeyToAdd:       shardBlsKeyToAdd,
			SlotKeyToAddSig:    shardBlsKeyToAddSig,
			EPOSStatus:         statusEnum,
		}
	}

	return payloadGenerator, nil
}

// EditStatus - edits the validator status for an existing validator
func EditStatus(
	keystore *keystore.KeyStore,
	account *accounts.Account,
	rpcClient *rpc.HTTPMessenger,
	chain *common.ChainID,
	validatorAddress string,
	status string,
	gasLimit int64,
	gasPrice numeric.Dec,
	nonce uint64,
	keystorePassphrase string,
	node string,
	timeout int,
) (map[string]interface{}, error) {
	statusEnum := determineEposStatus(status)

	payloadGenerator := editValidatorStatusGenerator(validatorAddress, statusEnum)

	var logMessage string
	if network.Verbose {
		logMessage = fmt.Sprintf("Generating a new edit validator status transaction:\n\tValidator Address: %s\n\tStatus: %v",
			validatorAddress,
			statusEnum,
		)
	}

	return staking.SendTx(keystore, account, rpcClient, chain, gasLimit, gasPrice, nonce, keystorePassphrase, node, timeout, payloadGenerator, logMessage)
}

func editValidatorStatusGenerator(
	validatorAddress string,
	statusEnum effective.Eligibility,
) hmyStaking.StakeMsgFulfiller {
	payloadGenerator := func() (hmyStaking.Directive, interface{}) {
		return hmyStaking.DirectiveEditValidator, hmyStaking.EditValidator{
			ValidatorAddress: address.Parse(validatorAddress),
			EPOSStatus:       statusEnum,
		}
	}

	return payloadGenerator
}
