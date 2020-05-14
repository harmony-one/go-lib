package delegation

import (
	"math/big"

	"github.com/harmony-one/harmony/common/denominations"
	"github.com/harmony-one/harmony/numeric"
)

// DelegationInfoWrapper - wrapper for the GetValidatorInformation RPC method
type DelegationInfoWrapper struct {
	ID      string           `json:"id" yaml:"id"`
	JSONRPC string           `json:"jsonrpc" yaml:"jsonrpc"`
	Result  []DelegationInfo `json:"result" yaml:"result"`
}

// DelegationInfo - the actual delegation info
type DelegationInfo struct {
	Undelegations    []UndelegationInfo `json:"Undelegations,omitempty" yaml:"Undelegations,omitempty"`
	ValidatorAddress string             `json:"validator_address,omitempty" yaml:"validator_address,omitempty"`
	DelegatorAddress string             `json:"delegator_address,omitempty" yaml:"delegator_address,omitempty"`
	RawAmount        *big.Int           `json:"amount,omitempty" yaml:"amount,omitempty"`
	Amount           numeric.Dec        `json:"-" yaml:"-"`
	RawReward        *big.Int           `json:"reward,omitempty" yaml:"reward,omitempty"`
	Reward           numeric.Dec        `json:"-" yaml:"-"`
}

// UndelegationInfo - represents the info for a given undelegation
type UndelegationInfo struct {
	RawAmount *big.Int    `json:"Amount" yaml:"Amount"`
	Amount    numeric.Dec `json:"-" yaml:"-"`
	Epoch     int         `json:"Epoch" yaml:"Epoch"`
}

// InitializeDelegationInfos - initializes a DelegationInfo slice
func InitializeDelegationInfos(delegations []DelegationInfo) ([]DelegationInfo, error) {
	if len(delegations) > 0 {
		rawDelegations := delegations
		delegations = []DelegationInfo{}

		for _, del := range rawDelegations {
			if err := del.Initialize(); err != nil {
				return delegations, err
			}

			delegations = append(delegations, del)
		}
	}

	return delegations, nil
}

// Initialize - initialize and convert values for a given ValidatorInfo struct
func (delegationInfo *DelegationInfo) Initialize() error {
	if delegationInfo.RawAmount != nil {
		decAmount := numeric.NewDecFromBigInt(delegationInfo.RawAmount)
		delegationInfo.Amount = decAmount.Quo(numeric.NewDec(denominations.One))
	}

	if delegationInfo.RawReward != nil {
		decReward := numeric.NewDecFromBigInt(delegationInfo.RawReward)
		delegationInfo.Reward = decReward.Quo(numeric.NewDec(denominations.One))
	}

	undelegations, err := InitializeUndelegationInfos(delegationInfo.Undelegations)
	if err != nil {
		return err
	}
	delegationInfo.Undelegations = undelegations

	return nil
}

// InitializeUndelegationInfos - initializes an UndelegationInfo slice
func InitializeUndelegationInfos(undelegations []UndelegationInfo) ([]UndelegationInfo, error) {
	if len(undelegations) > 0 {
		rawUndelegations := undelegations
		undelegations = []UndelegationInfo{}

		for _, undel := range rawUndelegations {
			if err := undel.Initialize(); err != nil {
				return undelegations, err
			}

			undelegations = append(undelegations, undel)
		}
	}

	return undelegations, nil
}

// Initialize - initialize and convert values for a given ValidatorInfo struct
func (undelegationInfo *UndelegationInfo) Initialize() error {
	if undelegationInfo.RawAmount != nil {
		decAmount := numeric.NewDecFromBigInt(undelegationInfo.RawAmount)
		undelegationInfo.Amount = decAmount.Quo(numeric.NewDec(denominations.One))
	}
	return nil
}
