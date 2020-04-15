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
	Validator      RPCValidator                  `json:"validator,omitempty" yaml:"validator,omitempty"`
	SigningPercent RPCCurrentEpochSigningPercent `json:"current-epoch-signing-percent,omitempty" yaml:"current-epoch-signing-percent,omitempty"`
	VotingPower    RPCCurrentEpochVotingPower    `json:"current-epoch-voting-power,omitempty" yaml:"current-epoch-voting-power,omitempty"`
}

// RPCCurrentEpochSigningPercent - the current epoch signing percentage
type RPCCurrentEpochSigningPercent struct {
	CurrentEpochSigned uint32 `json:"current-epoch-signed,omitempty" yaml:"current-epoch-signed,omitempty"`
	CurrentEpochToSign uint32 `json:"current-epoch-to-sign,omitempty" yaml:"current-epoch-to-sign,omitempty"`

	RawPercentage string      `json:"percentage,omitempty" yaml:"percentage,omitempty"`
	Percentage    numeric.Dec `json:"-" yaml:"-"`
}

// RPCCurrentEpochVotingPower - the current epoch voting power
type RPCCurrentEpochVotingPower struct {
	RawEffectiveStake string      `json:"effective-stake,omitempty" yaml:"effective-stake,omitempty"`
	EffectiveStake    numeric.Dec `json:"-" yaml:"-"`

	ShardID uint32 `json:"shard-id,omitempty" yaml:"shard-id,omitempty"`

	RawVotingPowerAdjusted string      `json:"voting-power-adjusted,omitempty" yaml:"voting-power-adjusted,omitempty"`
	VotingPowerAdjusted    numeric.Dec `json:"-" yaml:"-"`

	RawVotingPowerRaw string      `json:"voting-power-raw,omitempty" yaml:"voting-power-raw,omitempty"`
	VotingPowerRaw    numeric.Dec `json:"-" yaml:"-"`
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
	if err := validatorResult.Validator.Initialize(); err != nil {
		return err
	}
	if err := validatorResult.SigningPercent.Initialize(); err != nil {
		return err
	}
	if err := validatorResult.VotingPower.Initialize(); err != nil {
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

	if len(validatorInfo.Delegations) > 0 {
		for _, del := range validatorInfo.Delegations {
			if err := del.Initialize(); err != nil {
				return err
			}
		}
	}

	return nil
}

// Initialize - initialize and convert values for a given RPCValidatorResult struct
func (signingPercent *RPCCurrentEpochSigningPercent) Initialize() error {
	if signingPercent.RawPercentage != "" {
		decPercentage, err := common.NewDecFromString(signingPercent.RawPercentage)
		if err != nil {
			return errors.Wrapf(err, "SigningPercent: Percentage")
		}
		signingPercent.Percentage = decPercentage
	}

	return nil
}

// Initialize - initialize and convert values for a given RPCValidatorResult struct
func (votingPower *RPCCurrentEpochVotingPower) Initialize() error {
	if votingPower.RawEffectiveStake != "" {
		decEffectiveStake, err := common.NewDecFromString(votingPower.RawEffectiveStake)
		if err != nil {
			return errors.Wrapf(err, "VotingPower: EffectiveStake")
		}
		votingPower.EffectiveStake = decEffectiveStake
	}

	if votingPower.RawVotingPowerAdjusted != "" {
		decVotingPowerAdjusted, err := common.NewDecFromString(votingPower.RawVotingPowerAdjusted)
		if err != nil {
			return errors.Wrapf(err, "VotingPower: VotingPowerAdjusted")
		}
		votingPower.VotingPowerAdjusted = decVotingPowerAdjusted
	}

	if votingPower.RawVotingPowerRaw != "" {
		decVotingPowerRaw, err := common.NewDecFromString(votingPower.RawVotingPowerRaw)
		if err != nil {
			return errors.Wrapf(err, "VotingPower: VotingPowerRaw")
		}
		votingPower.VotingPowerRaw = decVotingPowerRaw
	}

	return nil
}
