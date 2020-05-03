package validator

import (
	"encoding/json"
	"fmt"

	"github.com/harmony-one/go-lib/transactions"
	goSdkRPC "github.com/harmony-one/go-sdk/pkg/rpc"
)

// AllInformation - retrieve all validator information
func AllInformation(node string, fetchAllPages bool) ([]RPCValidatorResult, error) {
	validatorResults := []RPCValidatorResult{}

	if fetchAllPages {
		page := 0
		for {
			validatorPagedResults, err := allInformationRequest(node, page)
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
		validatorPagedResults, err := allInformationRequest(node, 0)
		if err != nil {
			return validatorResults, err
		}
		validatorResults = append(validatorResults, validatorPagedResults...)
	}

	return validatorResults, nil
}

func allInformationRequest(node string, page int) ([]RPCValidatorResult, error) {
	response := RPCValidatorInfosWrapper{}
	results := []RPCValidatorResult{}

	bytes, err := goSdkRPC.RawRequest(goSdkRPC.Method.GetAllValidatorInformation, node, transactions.ParamsWrapper{page})
	if err != nil {
		return results, err
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

// All - retrieves all validators
func All(rpcClient *goSdkRPC.HTTPMessenger) (addresses []string, err error) {
	reply, err := rpcClient.SendRPC(goSdkRPC.Method.GetAllValidatorAddresses, transactions.ParamsWrapper{})
	if err != nil {
		return nil, err
	}

	return processAddressResponse(reply), nil
}

// Exists - checks if a given validator exists
func Exists(rpcClient *goSdkRPC.HTTPMessenger, validatorAddress string) bool {
	allValidators, err := All(rpcClient)
	if err == nil && len(allValidators) > 0 {
		for _, address := range allValidators {
			if address == validatorAddress {
				return true
			}
		}
	}

	return false
}

// AllElected - retrieves all active validators
func AllElected(rpcClient *goSdkRPC.HTTPMessenger) ([]string, error) {
	reply, err := rpcClient.SendRPC(goSdkRPC.Method.GetElectedValidatorAddresses, transactions.ParamsWrapper{})
	if err != nil {
		return nil, err
	}

	return processAddressResponse(reply), nil
}

// ElectedExists - checks if a given elected validator exists
func ElectedExists(rpcClient *goSdkRPC.HTTPMessenger, validatorAddress string) bool {
	allActiveValidators, err := AllElected(rpcClient)
	if err == nil && len(allActiveValidators) > 0 {
		for _, address := range allActiveValidators {
			if address == validatorAddress {
				return true
			}
		}
	}

	return false
}

func processAddressResponse(reply map[string]interface{}) (addresses []string) {
	for _, address := range reply["result"].([]interface{}) {
		addresses = append(addresses, address.(string))
	}

	return addresses
}

// Information - get the validator information for a given address
func Information(node string, validatorAddress string) (RPCValidatorResult, error) {
	response := RPCValidatorInfoWrapper{}
	result := RPCValidatorResult{}

	bytes, err := goSdkRPC.RawRequest(goSdkRPC.Method.GetValidatorInformation, node, transactions.ParamsWrapper{validatorAddress})
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
