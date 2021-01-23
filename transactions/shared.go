package transactions

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	eth_hexutil "github.com/ethereum/go-ethereum/common/hexutil"
	eth_rlp "github.com/ethereum/go-ethereum/rlp"
	"github.com/harmony-one/go-lib/network"
	"github.com/harmony-one/go-lib/rpc"
	goSdkRPC "github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/harmony/common/denominations"
	"github.com/harmony-one/harmony/core"
	"github.com/harmony-one/harmony/numeric"
)

// Copied from harmony-one/harmony/internal/params/protocol_params.go
const (
	// TxGas ...
	TxGas uint64 = 21000 // Per transaction not creating a contract. NOTE: Not payable on data of calls between transactions.
	// TxGasContractCreation ...
	TxGasContractCreation uint64 = 53000 // Per transaction that creates a contract. NOTE: Not payable on data of calls between transactions.
	// TxGasValidatorCreation ...
	TxGasValidatorCreation uint64 = 5300000 // Per transaction that creates a new validator. NOTE: Not payable on data of calls between transactions.
)

// Copied from harmony-one/go-sdk/pkg/sharding/structure.go
var (
	// NanoAsDec - Nano denomination in numeric.Dec
	NanoAsDec = numeric.NewDec(denominations.Nano)
	// OneAsDec - One denomination in numeric.Dec
	OneAsDec = numeric.NewDec(denominations.One)
)

// Transaction - represents an executed test case transaction
type Transaction struct {
	FromAddress     string
	FromShardID     uint32
	ToAddress       string
	ToShardID       uint32
	Data            string
	Amount          numeric.Dec
	GasPrice        int64
	Timeout         int
	TransactionHash string
	Success         bool
	Response        map[string]interface{}
	Error           error
}

// ToTransaction - converts a raw tx response map to a typed Transaction type
func ToTransaction(fromAddress string, fromShardID uint32, toAddress string, toShardID uint32, rawTx map[string]interface{}, err error) Transaction {
	if err != nil {
		return Transaction{Error: err}
	}

	var tx Transaction

	txHash := rawTx["transactionHash"].(string)

	if txHash != "" {
		success := IsTransactionSuccessful(rawTx)

		tx = Transaction{
			FromAddress:     fromAddress,
			FromShardID:     fromShardID,
			ToAddress:       toAddress,
			ToShardID:       toShardID,
			TransactionHash: txHash,
			Success:         success,
			Response:        rawTx,
		}
	}

	return tx
}

// EncodeSignature - RLP encodes a given transaction signature as a hex signature
func EncodeSignature(tx interface{}) (*string, error) {
	enc, err := eth_rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}

	hexSignature := eth_hexutil.Encode(enc)
	signature := &hexSignature

	return signature, nil
}

// SendRawTransaction - sends a raw signed transaction via RPC
func SendRawTransaction(rpcClient *goSdkRPC.HTTPMessenger, signature *string) (interface{}, error) {
	reply, err := rpcClient.SendRPC(goSdkRPC.Method.SendRawTransaction, []interface{}{signature})
	if err != nil {
		return nil, err
	}

	receiptHash, _ := reply["result"]

	return receiptHash, nil
}

// WaitForTxConfirmation - waits a given amount of seconds defined by timeout to try to receive a finalized transaction
func WaitForTxConfirmation(rpcClient *goSdkRPC.HTTPMessenger, node string, txType string, receiptHash string, timeout int) (map[string]interface{}, error) {
	var failures []rpc.Failure

	if timeout > 0 {
		for {
			if timeout < 0 {
				return nil, nil
			}

			if txType == "transaction" {
				failures, _ = rpc.TransactionFailures(node)
				if err := handleTransactionError(receiptHash, failures); err != nil {
					if network.Verbose {
						fmt.Println(fmt.Sprintf("\n[Harmony SDK]: %s - transaction error occurred for tx %s, error: %s", time.Now().Format(network.LoggingTimeFormat), receiptHash, err.Error()))
						fmt.Println("")
					}
					return nil, err
				}
			} else if txType == "staking" {
				failures, _ = rpc.StakingFailures(node)
				if err := handleTransactionError(receiptHash, failures); err != nil {
					if network.Verbose {
						fmt.Println(fmt.Sprintf("\n[Harmony SDK]: %s - staking error occurred for tx %s, error: %s", time.Now().Format(network.LoggingTimeFormat), receiptHash, err.Error()))
						fmt.Println("")
					}
					return nil, err
				}
			}

			response, err := GetTransactionReceipt(rpcClient, receiptHash)
			if err != nil {
				return nil, err
			}

			if response != nil {
				return response, nil
			}

			time.Sleep(time.Second * 1)
			timeout = timeout - 1
		}
	}

	return nil, nil
}

func handleTransactionError(receiptHash string, failures []rpc.Failure) error {
	failure, failed := rpc.FailureOccurredForTransaction(failures, receiptHash)
	if failed {
		return errors.New(failure.ErrorMessage)
	}

	return nil
}

// CalculateGasLimit - calculates the proper gas limit for a given gas limit and input data
func CalculateGasLimit(gasLimit int64, inputData string, isValidatorCreation bool) (calculatedGasLimit uint64, err error) {
	// -1 means that the gas limit has not been specified by the user and that it should be automatically calculated based on the tx data
	if gasLimit == -1 {
		if len(inputData) > 0 {
			base64InputData, err := base64.StdEncoding.DecodeString(inputData)
			if err != nil {
				return 0, err
			}

			calculatedGasLimit, err = core.IntrinsicGas(base64InputData, false, true, true, isValidatorCreation)
			if err != nil {
				return 0, err
			}

			if calculatedGasLimit == 0 {
				return 0, errors.New("calculated gas limit is 0 - this shouldn't be possible")
			}
		} else {
			calculatedGasLimit = TxGasContractCreation
		}
	} else {
		calculatedGasLimit = uint64(gasLimit)
	}

	return calculatedGasLimit, nil
}

// BumpGasPrice - bumps the gas price by the required percentage, as defined by core.DefaultTxPoolConfig.PriceBump
func BumpGasPrice(gasPrice numeric.Dec) numeric.Dec {
	//return gasPrice.Add(numeric.NewDec(1).Quo(OneAsDec))
	return gasPrice.Mul(numeric.NewDec(100 + int64(core.DefaultTxPoolConfig.PriceBump)).Quo(numeric.NewDec(100)))
}

// GetTransactionReceipt - retrieves the transaction info/data for a transaction
func GetTransactionReceipt(rpcClient *goSdkRPC.HTTPMessenger, receiptHash interface{}) (map[string]interface{}, error) {
	response, err := rpcClient.SendRPC(goSdkRPC.Method.GetTransactionReceipt, []interface{}{receiptHash})
	if err != nil {
		return nil, err
	}

	if response["result"] != nil {
		return response["result"].(map[string]interface{}), nil
	}

	return nil, nil
}

// IsTransactionSuccessful - checks if a transaction is successful given a transaction response
func IsTransactionSuccessful(txResponse map[string]interface{}) (success bool) {
	txStatus, ok := txResponse["status"].(string)

	if txStatus != "" && ok {
		success = (txStatus == "0x1")
	}

	return success
}

// GenerateTxData - generates tx data based on a given byte size
func GenerateTxData(char string, byteSize int) string {
	buffer := new(bytes.Buffer)

	for i := 0; i < byteSize; i++ {
		buffer.Write([]byte(char))
	}

	return buffer.String()
}
