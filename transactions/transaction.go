package transactions

import (
	"fmt"
	"math/big"
	"time"

	libErrors "github.com/harmony-one/go-lib/errors"
	"github.com/harmony-one/go-lib/network"
	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/common"
	goSdkRPC "github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/go-sdk/pkg/transaction"
	"github.com/harmony-one/harmony/accounts"
	"github.com/harmony-one/harmony/accounts/keystore"
	"github.com/harmony-one/harmony/core/types"
	"github.com/harmony-one/harmony/numeric"
)

// SendTransaction - send transactions
func SendTransaction(keystore *keystore.KeyStore, account *accounts.Account, rpcClient *goSdkRPC.HTTPMessenger, chain *common.ChainID, fromAddress string, fromShardID uint32, toAddress string, toShardID uint32, amount numeric.Dec, gasLimit int64, gasPrice numeric.Dec, nonce uint64, inputData string, keystorePassphrase string, node string, timeout int) (map[string]interface{}, error) {
	if keystore == nil || account == nil {
		return nil, libErrors.ErrMissingAccount
	}

	signedTx, err := GenerateAndSignTransaction(
		keystore,
		account,
		chain,
		fromAddress,
		fromShardID,
		toAddress,
		toShardID,
		amount,
		gasLimit,
		gasPrice,
		nonce,
		inputData,
	)
	if err != nil {
		return nil, err
	}

	signature, err := EncodeSignature(signedTx)
	if err != nil {
		return nil, err
	}

	if network.Verbose {
		fmt.Printf("\n[Harmony SDK]: %s - signed transaction using chain: %s (id: %d), signature: %v\n", time.Now().Format(network.LoggingTimeFormat), chain.Name, chain.Value, signature)
		json, _ := signedTx.MarshalJSON()
		fmt.Printf("%s\n", common.JSONPrettyFormat(string(json)))
		fmt.Printf("\n[Harmony SDK]: %s - sending transaction using node: %s, chain: %s (id: %d), signature: %v, timeout: %d\n\n", time.Now().Format(network.LoggingTimeFormat), node, chain.Name, chain.Value, signature, timeout)
	}

	receiptHash, err := SendRawTransaction(rpcClient, signature)
	if err != nil {
		return nil, err
	}

	if timeout > 0 {
		hash := receiptHash.(string)
		result, err := WaitForTxConfirmation(rpcClient, node, "transaction", hash, timeout)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	result := make(map[string]interface{})
	result["transactionHash"] = receiptHash

	return result, nil
}

// GenerateAndSignTransaction - generates and signs a transaction based on the supplied tx params and keystore/account
func GenerateAndSignTransaction(
	keystore *keystore.KeyStore,
	account *accounts.Account,
	chain *common.ChainID,
	fromAddress string,
	fromShardID uint32,
	toAddress string,
	toShardID uint32,
	amount numeric.Dec,
	gasLimit int64,
	gasPrice numeric.Dec,
	nonce uint64,
	inputData string,
) (tx *types.Transaction, err error) {
	generatedTx, err := GenerateTransaction(fromAddress, fromShardID, toAddress, toShardID, amount, gasLimit, gasPrice, nonce, inputData)
	if err != nil {
		return nil, err
	}

	tx, err = SignTransaction(keystore, account, generatedTx, chain.Value)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// GenerateTransaction - generate a new transaction
func GenerateTransaction(
	fromAddress string,
	fromShardID uint32,
	toAddress string,
	toShardID uint32,
	amount numeric.Dec,
	gasLimit int64,
	gasPrice numeric.Dec,
	nonce uint64,
	inputData string,
) (tx *types.Transaction, err error) {
	calculatedGasLimit, err := CalculateGasLimit(gasLimit, inputData, false)
	if err != nil {
		return nil, err
	}

	if network.Verbose {
		fmt.Println(fmt.Sprintf("\n[Harmony SDK]: %s - Generating a new transaction:\n\tReceiver address: %s\n\tFrom shard: %d\n\tTo shard: %d\n\tAmount: %f\n\tNonce: %d\n\tGas limit: %d\n\tGas price: %f\n\tData length (bytes): %d\n",
			time.Now().Format(network.LoggingTimeFormat),
			toAddress,
			fromShardID,
			toShardID,
			amount,
			nonce,
			calculatedGasLimit,
			gasPrice,
			len(inputData)),
		)
	}

	tx = transaction.NewTransaction(
		nonce,
		calculatedGasLimit,
		address.Parse(toAddress),
		fromShardID,
		toShardID,
		amount.Mul(OneAsDec),
		gasPrice.Mul(NanoAsDec),
		[]byte(inputData),
	)

	return tx, nil
}

// SignTransaction - signs a transaction using a given keystore / account
func SignTransaction(keystore *keystore.KeyStore, account *accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	signedTransaction, err := keystore.SignTx(*account, tx, chainID)
	if err != nil {
		return nil, err
	}

	return signedTransaction, nil
}

// AttachSigningData - attaches the signing data to the tx - necessary for e.g. propagating txs directly via p2p
func AttachSigningData(chainID *big.Int, tx *types.Transaction) (*types.Transaction, error) {
	signer := types.NewEIP155Signer(chainID)

	_, err := types.Sender(signer, tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
