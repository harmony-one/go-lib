package validator

import (
	goSdkRPC "github.com/harmony-one/go-sdk/pkg/rpc"
)

// All - retrieves all validators
func All(rpcClient *goSdkRPC.HTTPMessenger) (addresses []string, err error) {
	reply, err := rpcClient.SendRPC(goSdkRPC.Method.GetAllValidatorAddresses, []interface{}{})
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
	reply, err := rpcClient.SendRPC(goSdkRPC.Method.GetElectedValidatorAddresses, []interface{}{})
	if err != nil {
		return nil, err
	}

	return processAddressResponse(reply), nil
}

func processAddressResponse(reply map[string]interface{}) (addresses []string) {
	for _, address := range reply["result"].([]interface{}) {
		addresses = append(addresses, address.(string))
	}

	return addresses
}
