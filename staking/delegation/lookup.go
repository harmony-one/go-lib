package delegation

import (
	"encoding/json"

	"github.com/harmony-one/go-lib/transactions"
	goSdkRPC "github.com/harmony-one/go-sdk/pkg/rpc"
)

// ByValidator - get delegations by validator
func ByValidator(node string, address string) ([]DelegationInfo, error) {
	return lookupDelegation(node, goSdkRPC.Method.GetDelegationsByValidator, address)
}

// ByDelegator - get delegations by delegator
func ByDelegator(node string, address string) ([]DelegationInfo, error) {
	return lookupDelegation(node, goSdkRPC.Method.GetDelegationsByDelegator, address)
}

func lookupDelegation(node string, rpcMethod string, address string) ([]DelegationInfo, error) {
	response := DelegationInfoWrapper{}
	delegationInfo := []DelegationInfo{}

	bytes, err := goSdkRPC.RawRequest(rpcMethod, node, transactions.ParamsWrapper{address})
	if err != nil {
		return delegationInfo, err
	}

	json.Unmarshal(bytes, &response)
	delegationInfo = response.Result
	if len(delegationInfo) > 0 {
		for _, info := range delegationInfo {
			info.Initialize()
		}
	}

	return delegationInfo, nil
}
