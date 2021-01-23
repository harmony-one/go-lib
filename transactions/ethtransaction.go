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

// SendEthTransaction - send eth transactions
func SendEthTransaction(keystore *keystore.KeyStore, account *accounts.Account, rpcClient *goSdkRPC.HTTPMessenger, chain *common.ChainID, fromAddress string, toAddress string, amount numeric.Dec, gasLimit int64, gasPrice numeric.Dec, nonce uint64, inputData string, keystorePassphrase string, node string, timeout int) (map[string]interface{}, error) {
	if keystore == nil || account == nil {
		return nil, libErrors.ErrMissingAccount
	}

	signedTx, err := GenerateAndSignEthTransaction(
		keystore,
		account,
		chain,
		fromAddress,
		toAddress,
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

// GenerateAndSignEthTransaction - generates and signs a transaction based on the supplied tx params and keystore/account
func GenerateAndSignEthTransaction(
	keystore *keystore.KeyStore,
	account *accounts.Account,
	chain *common.ChainID,
	fromAddress string,
	toAddress string,
	amount numeric.Dec,
	gasLimit int64,
	gasPrice numeric.Dec,
	nonce uint64,
	inputData string,
) (tx *types.EthTransaction, err error) {
	generatedTx, err := GenerateEthTransaction(fromAddress, toAddress, amount, gasLimit, gasPrice, nonce, inputData)
	if err != nil {
		return nil, err
	}

	tx, err = SignEthTransaction(keystore, account, generatedTx, chain.Value)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// GenerateEthTransaction - generate a new transaction
func GenerateEthTransaction(
	fromAddress string,
	toAddress string,
	amount numeric.Dec,
	gasLimit int64,
	gasPrice numeric.Dec,
	nonce uint64,
	inputData string,
) (tx *types.EthTransaction, err error) {
	calculatedGasLimit, err := CalculateGasLimit(gasLimit, inputData, false)
	if err != nil {
		return nil, err
	}

	if network.Verbose {
		fmt.Println(fmt.Sprintf("\n[Harmony SDK]: %s - Generating a new transaction:\n\tReceiver address: %s\n\tAmount: %f\n\tNonce: %d\n\tGas limit: %d\n\tGas price: %f\n\tData length (bytes): %d\n",
			time.Now().Format(network.LoggingTimeFormat),
			toAddress,
			amount,
			nonce,
			calculatedGasLimit,
			gasPrice,
			len(inputData)),
		)
	}

	tx = transaction.NewEthTransaction(
		nonce,
		calculatedGasLimit,
		address.Parse(toAddress),
		amount.Mul(OneAsDec),
		gasPrice.Mul(NanoAsDec),
		[]byte(inputData),
	)

	return tx, nil
}

// SignEthTransaction - signs a transaction using a given keystore / account
func SignEthTransaction(keystore *keystore.KeyStore, account *accounts.Account, tx *types.EthTransaction, chainID *big.Int) (*types.EthTransaction, error) {
	signedTransaction, err := keystore.SignEthTx(*account, tx, chainID)
	if err != nil {
		return nil, err
	}

	return signedTransaction, nil
}

// AttachEthSigningData - attaches the signing data to the tx - necessary for e.g. propagating txs directly via p2p
func AttachEthSigningData(chainID *big.Int, tx *types.EthTransaction) (*types.EthTransaction, error) {
	signer := types.NewEIP155Signer(chainID)

	_, err := types.Sender(signer, tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
