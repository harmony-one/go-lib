package network

import (
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/harmony/numeric"
	"github.com/pkg/errors"
)

// Gas - represents the gas settings
type Gas struct {
	RawCost  string      `json:"cost" yaml:"cost"`
	Cost     numeric.Dec `json:"-" yaml:"-"`
	Limit    int64       `json:"limit" yaml:"limit"`
	RawPrice string      `json:"price" yaml:"price"`
	Price    numeric.Dec `json:"-" yaml:"-"`
}

// Initialize - convert the raw values to their appropriate numeric.Dec values
func (gas *Gas) Initialize() error {
	if gas.RawCost != "" {
		decCost, err := common.NewDecFromString(gas.RawCost)
		if err != nil {
			return errors.Wrapf(err, "Gas: Cost")
		}
		gas.Cost = decCost
	}

	if gas.RawPrice != "" {
		decPrice, err := common.NewDecFromString(gas.RawPrice)
		if err != nil {
			return errors.Wrapf(err, "Gas: Price")
		}
		gas.Price = decPrice
	}

	if gas.Limit == 0 {
		gas.Limit = -1
	}

	if gas.Price.IsNil() || gas.Price.IsZero() || gas.Price.IsNegative() {
		gas.Price = numeric.NewDec(1)
	}

	return nil
}
