package transactions

import (
	"github.com/harmony-one/harmony/numeric"
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
