package delegation

import (
	"github.com/pkg/errors"

	"github.com/harmony-one/go-sdk/pkg/common"
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
	DelegatorAddress string             `json:"delegator-address,omitempty" yaml:"delegator-address,omitempty"`
	ValidatorAddress string             `json:"validator-address,omitempty" yaml:"validator-address,omitempty"`
	RawAmount        string             `json:"amount,omitempty" yaml:"amount,omitempty"`
	Amount           numeric.Dec        `json:"-" yaml:"-"`
	RawReward        string             `json:"reward,omitempty" yaml:"reward,omitempty"`
	Reward           numeric.Dec        `json:"-" yaml:"-"`
}

// UndelegationInfo - represents the info for a given undelegation
type UndelegationInfo struct {
	RawAmount string      `json:"Amount" yaml:"Amount"`
	Amount    numeric.Dec `json:"-" yaml:"-"`
	Epoch     int         `json:"Epoch" yaml:"Epoch"`
}

// Initialize - initialize and convert values for a given ValidatorInfo struct
func (delegationInfo *DelegationInfo) Initialize() error {
	if delegationInfo.RawAmount != "" {
		decAmount, err := common.NewDecFromString(delegationInfo.RawAmount)
		if err != nil {
			return errors.Wrapf(err, "DelegationInfo: Amount")
		}
		delegationInfo.Amount = decAmount.Quo(numeric.NewDec(denominations.One))
	}

	if delegationInfo.RawReward != "" {
		decReward, err := common.NewDecFromString(delegationInfo.RawReward)
		if err != nil {
			return errors.Wrapf(err, "DelegationInfo: Reward")
		}
		delegationInfo.Reward = decReward.Quo(numeric.NewDec(denominations.One))
	}

	if len(delegationInfo.Undelegations) > 0 {
		for _, undel := range delegationInfo.Undelegations {
			if err := undel.Initialize(); err != nil {
				return err
			}
		}
	}

	return nil
}

// Initialize - initialize and convert values for a given ValidatorInfo struct
func (undelegationInfo *UndelegationInfo) Initialize() error {
	if undelegationInfo.RawAmount != "" {
		decAmount, err := common.NewDecFromString(undelegationInfo.RawAmount)
		if err != nil {
			return errors.Wrapf(err, "UndelegationInfo: Amount")
		}
		undelegationInfo.Amount = decAmount.Quo(numeric.NewDec(denominations.One))
	}
	return nil
}
