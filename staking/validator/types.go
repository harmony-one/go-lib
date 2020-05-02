package validator

import (
	"math/big"

	"github.com/pkg/errors"

	"github.com/harmony-one/go-lib/staking/delegation"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/harmony/common/denominations"
	"github.com/harmony-one/harmony/numeric"
)

// RPCValidatorInfoWrapper - wrapper for the GetValidatorInformation RPC method
type RPCValidatorInfoWrapper struct {
	ID      string             `json:"id" yaml:"id"`
	JSONRPC string             `json:"jsonrpc" yaml:"jsonrpc"`
	Result  RPCValidatorResult `json:"result" yaml:"result"`
	Error   RPCError           `json:"error,omitempty" yaml:"error,omitempty"`
}

// RPCError - error returned from the RPC endpoint
type RPCError struct {
	Code    int    `json:"code,omitempty" yaml:"code,omitempty"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
}

// RPCValidatorResult - the actual result
type RPCValidatorResult struct {
	CurrentEpochPerformance RPCCurrentEpochPerformance `json:"current-epoch-performance,omitempty" yaml:"current-epoch-performance,omitempty"`
	Validator               RPCValidator               `json:"validator,omitempty" yaml:"validator,omitempty"`
	CurrentlyInCommittee    bool                       `json:"currently-in-committee,omitempty" yaml:"currently-in-committee,omitempty"`
	EposStatus              string                     `json:"epos-status,omitempty" yaml:"epos-status,omitempty"`
	RawTotalDelegation      *big.Int                   `json:"total-delegation,omitempty" yaml:"total-delegation,omitempty"`
	TotalDelegation         numeric.Dec                `json:"-" yaml:"-"`
	Lifetime                RPCValidatorLifetime       `json:"lifetime,omitempty" yaml:"lifetime,omitempty"`
}

// RPCCurrentEpochPerformance - the current epoch performance
type RPCCurrentEpochPerformance struct {
	CurrentEpochSigned               uint32      `json:"current-epoch-signed,omitempty" yaml:"current-epoch-signed,omitempty"`
	CurrentEpochToSign               uint32      `json:"current-epoch-to-sign,omitempty" yaml:"current-epoch-to-sign,omitempty"`
	RawCurrentEpochSigningPercentage string      `json:"current-epoch-signing-percentage,omitempty" yaml:"current-epoch-signing-percentage,omitempty"`
	CurrentEpochSigningPercentage    numeric.Dec `json:"-" yaml:"-"`
}

// RPCValidatorLifetime - validator lifetime rewards
type RPCValidatorLifetime struct {
	RawAPR string      `json:"apr,omitempty" yaml:"apr,omitempty"`
	APR    numeric.Dec `json:"-" yaml:"-"`

	RawRewardAccumulated *big.Int    `json:"reward-accumulated,omitempty" yaml:"reward-accumulated,omitempty"`
	RewardAccumulated    numeric.Dec `json:"-" yaml:"-"`

	Blocks RPCValidatorBlockStatistics `json:"blocks,omitempty" yaml:"blocks,omitempty"`
}

// RPCValidatorBlockStatistics - validator block statistics
type RPCValidatorBlockStatistics struct {
	Signed int `json:"signed,omitempty" yaml:"signed,omitempty"`
	ToSign int `json:"to-sign,omitempty" yaml:"to-sign,omitempty"`
}

// RPCValidator - the actual validator info
type RPCValidator struct {
	Address               string                      `json:"address,omitempty" yaml:"address,omitempty"`
	BLSPublicKeys         []string                    `json:"bls-public-keys,omitempty" yaml:"bls-public-keys,omitempty"`
	CreationHeight        uint32                      `json:"creation-height,omitempty" yaml:"creation-height,omitempty"`
	RawMaxTotalDelegation *big.Int                    `json:"max-total-delegation,omitempty" yaml:"max-total-delegation,omitempty"`
	MaxTotalDelegation    numeric.Dec                 `json:"-" yaml:"-"`
	RawMinSelfDelegation  *big.Int                    `json:"min-self-delegation,omitempty" yaml:"min-self-delegation,omitempty"`
	MinSelfDelegation     numeric.Dec                 `json:"-" yaml:"-"`
	Name                  string                      `json:"name,omitempty" yaml:"name,omitempty"`                         // name
	Identity              string                      `json:"identity,omitempty" yaml:"identity,omitempty"`                 // optional identity signature (ex. UPort or Keybase)
	Website               string                      `json:"website,omitempty" yaml:"website,omitempty"`                   // optional website link
	SecurityContact       string                      `json:"security-contact,omitempty" yaml:"security-contact,omitempty"` // optional security contact info
	Details               string                      `json:"details,omitempty" yaml:"details,omitempty"`                   // optional details
	RawRate               string                      `json:"rate,omitempty" yaml:"rate,omitempty"`
	Rate                  numeric.Dec                 `json:"-" yaml:"-"`
	RawMaxChangeRate      string                      `json:"max-change-rate,omitempty" yaml:"max-change-rate,omitempty"`
	MaxChangeRate         numeric.Dec                 `json:"-" yaml:"-"`
	RawMaxRate            string                      `json:"max-rate,omitempty" yaml:"max-rate,omitempty"`
	MaxRate               numeric.Dec                 `json:"-" yaml:"-"`
	EligibilityStatus     string                      `json:"epos-eligibility-status,omitempty" yaml:"epos-eligibility-status,omitempty"`
	LastEpochInCommittee  uint32                      `json:"last-epoch-in-committee,omitempty" yaml:"last-epoch-in-committee,omitempty"`
	UpdateHeight          uint32                      `json:"update-height,omitempty" yaml:"update-height,omitempty"`
	Availability          RPCValidatorAvailability    `json:"availability,omitempty" yaml:"availability,omitempty"`
	Delegations           []delegation.DelegationInfo `json:"delegations,omitempty" yaml:"delegations,omitempty"`
}

// RPCValidatorAvailability - validator availability
type RPCValidatorAvailability struct {
	BlocksSigned uint32 `json:"num-blocks-signed,omitempty" yaml:"num-blocks-signed,omitempty"`
	BlocksToSign uint32 `json:"num-blocks-to-sign,omitempty" yaml:"num-blocks-to-sign,omitempty"`
}

// Initialize - initialize and convert values for a given RPCValidatorResult struct
func (validatorResult *RPCValidatorResult) Initialize() error {
	if validatorResult.RawTotalDelegation != nil {
		decTotalDelegation := numeric.NewDecFromBigInt(validatorResult.RawTotalDelegation)
		validatorResult.TotalDelegation = decTotalDelegation.Quo(numeric.NewDec(denominations.One))
	}

	if err := validatorResult.Validator.Initialize(); err != nil {
		return err
	}

	if err := validatorResult.CurrentEpochPerformance.Initialize(); err != nil {
		return err
	}

	if err := validatorResult.Lifetime.Initialize(); err != nil {
		return err
	}

	return nil
}

// Initialize - initialize and convert values for a given ValidatorInfo struct
func (validatorInfo *RPCValidator) Initialize() error {
	if validatorInfo.RawMaxTotalDelegation != nil {
		decMaxTotalDelegation := numeric.NewDecFromBigInt(validatorInfo.RawMaxTotalDelegation)
		validatorInfo.MaxTotalDelegation = decMaxTotalDelegation.Quo(numeric.NewDec(denominations.One))
	}

	if validatorInfo.RawMinSelfDelegation != nil {
		decMinSelfDelegation := numeric.NewDecFromBigInt(validatorInfo.RawMinSelfDelegation)
		validatorInfo.MinSelfDelegation = decMinSelfDelegation.Quo(numeric.NewDec(denominations.One))
	}

	if validatorInfo.RawRate != "" {
		decRate, err := common.NewDecFromString(validatorInfo.RawRate)
		if err != nil {
			return errors.Wrapf(err, "ValidatorInfo: Rate")
		}
		validatorInfo.Rate = decRate
	}

	if validatorInfo.RawMaxChangeRate != "" {
		decMaxChangeRate, err := common.NewDecFromString(validatorInfo.RawMaxChangeRate)
		if err != nil {
			return errors.Wrapf(err, "ValidatorInfo: MaxChangeRate")
		}
		validatorInfo.MaxChangeRate = decMaxChangeRate
	}

	if validatorInfo.RawMaxRate != "" {
		decMaxRate, err := common.NewDecFromString(validatorInfo.RawMaxRate)
		if err != nil {
			return errors.Wrapf(err, "ValidatorInfo: MaxRate")
		}
		validatorInfo.MaxRate = decMaxRate
	}

	delegations, err := delegation.InitializeDelegationInfos(validatorInfo.Delegations)
	if err != nil {
		return err
	}
	validatorInfo.Delegations = delegations

	return nil
}

// Initialize - initialize and convert values for a given RPCValidatorResult struct
func (epochPerformance *RPCCurrentEpochPerformance) Initialize() error {
	if epochPerformance.RawCurrentEpochSigningPercentage != "" {
		decPercentage, err := common.NewDecFromString(epochPerformance.RawCurrentEpochSigningPercentage)
		if err != nil {
			return errors.Wrapf(err, "EpochPerformance: CurrentEpochSigningPercentage")
		}
		epochPerformance.CurrentEpochSigningPercentage = decPercentage
	}

	return nil
}

// Initialize - initialize and convert values for a given RPCValidatorResult struct
func (lifetime *RPCValidatorLifetime) Initialize() error {
	if lifetime.RawAPR != "" {
		decAPR, err := common.NewDecFromString(lifetime.RawAPR)
		if err != nil {
			return errors.Wrapf(err, "ValidatorLifetime: APR")
		}
		lifetime.APR = decAPR
	}

	if lifetime.RawRewardAccumulated != nil {
		decRewardAccumulated := numeric.NewDecFromBigInt(lifetime.RawRewardAccumulated)
		lifetime.RewardAccumulated = decRewardAccumulated.Quo(numeric.NewDec(denominations.One))
	}

	return nil
}
