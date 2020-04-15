package rpc

import (
	"encoding/json"

	goSdkRPC "github.com/harmony-one/go-sdk/pkg/rpc"
)

// FailureWrapper - wrapper for the GetCurrentTransactionErrorSink / GetCurrentStakingErrorSink RPC methods
type FailureWrapper struct {
	ID       string    `json:"id" yaml:"id"`
	JSONRPC  string    `json:"jsonrpc" yaml:"jsonrpc"`
	Failures []Failure `json:"result" yaml:"result"`
}

// Failure - the actual failure
type Failure struct {
	DirectiveKind   string `json:"directive-kind,omitempty" yaml:"directive-kind,omitempty"`
	ErrorMessage    string `json:"error-message" yaml:"error-message"`
	TimeAtRejection uint32 `json:"time-at-rejection" yaml:"time-at-rejection"`
	TxHashID        string `json:"tx-hash-id" yaml:"tx-hash-id"`
}

// AllFailures - get both transaction and staking failures at the same time
func AllFailures(node string) (failures []Failure, err error) {
	txFailures, err := TransactionFailures(node)
	if err != nil {
		return failures, err
	}
	failures = append(failures, txFailures...)

	stakingFailures, err := StakingFailures(node)
	if err != nil {
		return failures, err
	}
	failures = append(failures, stakingFailures...)

	return failures, nil
}

// TransactionFailures - get the transaction failures from the given node
func TransactionFailures(node string) ([]Failure, error) {
	response := FailureWrapper{}
	failures := []Failure{}

	bytes, err := goSdkRPC.RawRequest(goSdkRPC.Method.GetCurrentTransactionErrorSink, node, []interface{}{})
	if err != nil {
		return failures, err
	}

	json.Unmarshal(bytes, &response)
	failures = response.Failures

	return failures, nil
}

// StakingFailures - get the staking failures from the given node
func StakingFailures(node string) ([]Failure, error) {
	response := FailureWrapper{}
	failures := []Failure{}

	bytes, err := goSdkRPC.RawRequest(goSdkRPC.Method.GetCurrentStakingErrorSink, node, []interface{}{})
	if err != nil {
		return failures, err
	}

	json.Unmarshal(bytes, &response)
	failures = response.Failures

	return failures, nil
}

// FailureOccurredForTransaction - checks if a failure has occurred for a given transaction
func FailureOccurredForTransaction(failures []Failure, receiptHash string) (Failure, bool) {
	for _, failure := range failures {
		if failure.TxHashID == receiptHash {
			return failure, true
		}
	}

	return Failure{}, false
}
