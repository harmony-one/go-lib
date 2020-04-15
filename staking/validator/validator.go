package validator

import (
	"strconv"

	"github.com/harmony-one/go-lib/accounts"
	"github.com/harmony-one/go-lib/crypto"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/harmony/numeric"
	hmyStaking "github.com/harmony-one/harmony/staking/types"
	"github.com/pkg/errors"
)

// Validator - represents the validator details
type Validator struct {
	RawShardID string `yaml:"shard_id"`
	ShardID    uint32 `yaml:"-"`
	Account    *accounts.Account
	Details    ValidatorDetails `yaml:"details"`
	Commission Commission       `yaml:"commission"`
	BLSKeys    []crypto.BLSKey  `yaml:"-"`
	Exists     bool

	RawMinimumSelfDelegation string      `yaml:"minimum_self_delegation"`
	MinimumSelfDelegation    numeric.Dec `yaml:"-"`

	RawMaximumTotalDelegation string      `yaml:"maximum_total_delegation"`
	MaximumTotalDelegation    numeric.Dec `yaml:"-"`

	RawAmount string      `yaml:"amount"`
	Amount    numeric.Dec `yaml:"-"`

	EligibilityStatus string `yaml:"eligibility-status"`
}

// ValidatorDetails - represents the validator details
type ValidatorDetails struct {
	Name            string `yaml:"name"`
	Identity        string `yaml:"identity"`
	Website         string `yaml:"website"`
	SecurityContact string `yaml:"security_contact"`
	Details         string `yaml:"details"`
}

// Commission - represents the validator commission settings
type Commission struct {
	RawRate string      `yaml:"rate"`
	Rate    numeric.Dec `yaml:"-"`

	RawMaxRate string      `yaml:"max_rate"`
	MaxRate    numeric.Dec `yaml:"-"`

	RawMaxChangeRate string      `yaml:"max_change_rate"`
	MaxChangeRate    numeric.Dec `yaml:"-"`
}

// Initialize - initializes and converts values for a given validator
func (validator *Validator) Initialize() error {
	// Set the shard id - will be used when generating bls keys
	if validator.RawShardID != "" {
		converted, err := strconv.ParseInt(validator.RawShardID, 10, 64)
		if err != nil {
			return err
		}
		validator.ShardID = uint32(converted)
	} else {
		validator.ShardID = 0
	}

	if validator.RawMinimumSelfDelegation != "" {
		decMinimumSelfDelegation, err := common.NewDecFromString(validator.RawMinimumSelfDelegation)
		if err != nil {
			return errors.Wrapf(err, "Validator: MinimumSelfDelegation")
		}
		validator.MinimumSelfDelegation = decMinimumSelfDelegation
	}

	if validator.RawMaximumTotalDelegation != "" {
		decMaximumTotalDelegation, err := common.NewDecFromString(validator.RawMaximumTotalDelegation)
		if err != nil {
			return errors.Wrapf(err, "Validator: MaximumTotalDelegation")
		}
		validator.MaximumTotalDelegation = decMaximumTotalDelegation
	}

	if validator.RawAmount != "" {
		decAmount, err := common.NewDecFromString(validator.RawAmount)
		if err != nil {
			return errors.Wrapf(err, "Validator: Amount")
		}
		validator.Amount = decAmount
	}

	// Initialize commission values
	if err := validator.Commission.Initialize(); err != nil {
		return err
	}

	return nil
}

// Initialize - initializes and converts values for a given test case
func (commission *Commission) Initialize() error {
	if commission.RawRate != "" {
		decRate, err := common.NewDecFromString(commission.RawRate)
		if err != nil {
			return errors.Wrapf(err, "Commission: Rate")
		}
		commission.Rate = decRate
	}

	if commission.RawMaxRate != "" {
		decMaxRate, err := common.NewDecFromString(commission.RawMaxRate)
		if err != nil {
			return errors.Wrapf(err, "Commission: MaxRate")
		}
		commission.MaxRate = decMaxRate
	}

	if commission.RawMaxChangeRate != "" {
		decMaxChangeRate, err := common.NewDecFromString(commission.RawMaxChangeRate)
		if err != nil {
			return errors.Wrapf(err, "Commission: MaxChangeRate")
		}
		commission.MaxChangeRate = decMaxChangeRate
	}

	return nil
}

// ToStakingDescription - convert validator details to a suitable format for staking txs
func (validator *Validator) ToStakingDescription() hmyStaking.Description {
	return hmyStaking.Description{
		Name:            validator.Details.Name,
		Identity:        validator.Details.Identity,
		Website:         validator.Details.Website,
		SecurityContact: validator.Details.SecurityContact,
		Details:         validator.Details.Details,
	}
}

// ToCommissionRates - convert validator commission rates to a suitable format for staking txs
func (validator *Validator) ToCommissionRates() hmyStaking.CommissionRates {
	return hmyStaking.CommissionRates{
		Rate:          validator.Commission.Rate,
		MaxRate:       validator.Commission.MaxRate,
		MaxChangeRate: validator.Commission.MaxChangeRate,
	}
}
