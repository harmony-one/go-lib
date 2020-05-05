package validator

import (
	"encoding/json"
	"fmt"

	goSdkRPC "github.com/harmony-one/go-sdk/pkg/rpc"
)

// AllInformationForBlock - get all validator info for a given block
func AllInformationForBlock(node string, blockNumber int, fetchAllPages bool) ([]RPCValidatorResult, error) {
	validatorResults := []RPCValidatorResult{}

	if fetchAllPages {
		page := 0
		for {
			validatorPagedResults, err := allInformationRequest(node, page, blockNumber)
			if err != nil {
				return validatorResults, err
			}
			if len(validatorPagedResults) > 0 {
				validatorResults = append(validatorResults, validatorPagedResults...)
			} else {
				break
			}
			page++
		}
	} else {
		validatorPagedResults, err := allInformationRequest(node, 0, blockNumber)
		if err != nil {
			return validatorResults, err
		}
		validatorResults = append(validatorResults, validatorPagedResults...)
	}

	return validatorResults, nil
}

// AllInformation - retrieve all validator information
func AllInformation(node string, fetchAllPages bool) ([]RPCValidatorResult, error) {
	validatorResults := []RPCValidatorResult{}

	if fetchAllPages {
		page := 0
		for {
			validatorPagedResults, err := allInformationRequest(node, page, -1)
			if err != nil {
				return validatorResults, err
			}
			if len(validatorPagedResults) > 0 {
				validatorResults = append(validatorResults, validatorPagedResults...)
			} else {
				break
			}
			page++
		}
	} else {
		validatorPagedResults, err := allInformationRequest(node, 0, -1)
		if err != nil {
			return validatorResults, err
		}
		validatorResults = append(validatorResults, validatorPagedResults...)
	}

	return validatorResults, nil
}

func allInformationRequest(node string, page int, blockNumber int) ([]RPCValidatorResult, error) {
	response := RPCValidatorInfosWrapper{}
	results := []RPCValidatorResult{}
	var bytes []byte
	var err error

	if blockNumber >= 0 {
		hexBlockNumberParam := fmt.Sprintf("0x%x", blockNumber)
		bytes, err = goSdkRPC.RawRequest(goSdkRPC.Method.GetAllValidatorInformationByBlockNumber, node, []interface{}{page, hexBlockNumberParam})
		if err != nil {
			return results, err
		}
	} else {
		bytes, err = goSdkRPC.RawRequest(goSdkRPC.Method.GetAllValidatorInformation, node, []interface{}{page})
		if err != nil {
			return results, err
		}
	}

	json.Unmarshal(bytes, &response)

	if response.Error.Message != "" {
		return results, fmt.Errorf("%s (%d)", response.Error.Message, response.Error.Code)
	}

	results = response.Result

	validatorResults := []RPCValidatorResult{}
	for _, result := range results {
		if err := result.Initialize(); err != nil {
			return results, nil
		}
		validatorResults = append(validatorResults, result)
	}

	return validatorResults, nil
}

// Information - get the validator information for a given address
func Information(node string, validatorAddress string) (RPCValidatorResult, error) {
	response := RPCValidatorInfoWrapper{}
	result := RPCValidatorResult{}

	bytes, err := goSdkRPC.RawRequest(goSdkRPC.Method.GetValidatorInformation, node, []interface{}{validatorAddress})
	if err != nil {
		return result, err
	}

	json.Unmarshal(bytes, &response)

	if response.Error.Message != "" {
		return result, fmt.Errorf("%s (%d)", response.Error.Message, response.Error.Code)
	}

	result = response.Result
	result.Initialize()

	return result, nil
}
