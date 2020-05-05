package staking

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/harmony-one/go-lib/crypto"
	"github.com/harmony-one/go-lib/network"
	"github.com/harmony-one/go-lib/transactions"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/harmony/accounts"
	"github.com/harmony-one/harmony/accounts/keystore"
	"github.com/harmony-one/harmony/numeric"
	"github.com/harmony-one/harmony/shard"
	hmyStaking "github.com/harmony-one/harmony/staking/types"
)

// SendTx - generate the staking tx, sign it, encode the signature and send the actual tx data
func SendTx(
	keystore *keystore.KeyStore,
	account *accounts.Account,
	rpcClient *rpc.HTTPMessenger,
	chain *common.ChainID,
	gasLimit int64,
	gasPrice numeric.Dec,
	nonce uint64,
	keystorePassphrase string,
	node string,
	confirmationWaitTime int,
	payloadGenerator hmyStaking.StakeMsgFulfiller,
	logMessage string,
) (map[string]interface{}, error) {
	if keystore == nil || account == nil {
		return nil, errors.New("keystore account can't be nil - please make sure the account you want to use exists in the keystore")
	}
	stakingTx, calculatedGasLimit, err := GenerateStakingTransaction(gasLimit, gasPrice, nonce, payloadGenerator)
	if err != nil {
		return nil, err
	}

	signedTx, err := SignStakingTransaction(keystore, account, stakingTx, chain.Value)
	if err != nil {
		return nil, err
	}

	signature, err := transactions.EncodeSignature(signedTx)
	if err != nil {
		return nil, err
	}

	if logMessage != "" {
		logMessage = fmt.Sprintf("\n[Harmony SDK]: %s - %s\n\tGas Limit: %d\n\tGas Price: %f\n\tNonce: %d\n\tSignature: %v\n",
			time.Now().Format(network.LoggingTimeFormat),
			logMessage,
			calculatedGasLimit,
			gasPrice,
			nonce,
			signature,
		)
		fmt.Println(logMessage)
	}

	receiptHash, err := SendRawStakingTransaction(rpcClient, signature)
	if err != nil {
		return nil, err
	}

	if confirmationWaitTime > 0 {
		hash := receiptHash.(string)
		result, _ := transactions.WaitForTxConfirmation(rpcClient, node, "staking", hash, confirmationWaitTime)

		if result != nil {
			return result, nil
		}
	}

	result := make(map[string]interface{})
	result["transactionHash"] = receiptHash

	return result, nil
}

// GenerateStakingTransaction - generate a staking transaction
func GenerateStakingTransaction(gasLimit int64, gasPrice numeric.Dec, nonce uint64, payloadGenerator hmyStaking.StakeMsgFulfiller) (*hmyStaking.StakingTransaction, uint64, error) {
	directive, payload := payloadGenerator()
	isCreateValidator := (directive == hmyStaking.DirectiveCreateValidator)

	bytes, err := rlp.EncodeToBytes(payload)
	if err != nil {
		return nil, 0, err
	}

	data := base64.StdEncoding.EncodeToString(bytes)

	calculatedGasLimit, err := transactions.CalculateGasLimit(gasLimit, data, isCreateValidator)
	if err != nil {
		return nil, 0, err
	}

	stakingTx, err := hmyStaking.NewStakingTransaction(nonce, calculatedGasLimit, gasPrice.TruncateInt(), payloadGenerator)
	if err != nil {
		return nil, 0, err
	}

	return stakingTx, calculatedGasLimit, nil
}

// ProcessBlsKeys - separate bls keys to pub key and sig slices
func ProcessBlsKeys(blsKeys []crypto.BLSKey) (blsPubKeys []shard.BLSPublicKey, blsSigs []shard.BLSSignature) {
	blsPubKeys = make([]shard.BLSPublicKey, len(blsKeys))
	blsSigs = make([]shard.BLSSignature, len(blsKeys))

	for i, blsKey := range blsKeys {
		blsPubKeys[i] = *blsKey.ShardPublicKey
		blsSigs[i] = *blsKey.ShardSignature
	}

	return blsPubKeys, blsSigs
}

// SignStakingTransaction - sign a staking transaction
func SignStakingTransaction(keystore *keystore.KeyStore, account *accounts.Account, stakingTx *hmyStaking.StakingTransaction, chainID *big.Int) (*hmyStaking.StakingTransaction, error) {
	signedTransaction, err := keystore.SignStakingTx(*account, stakingTx, chainID)
	if err != nil {
		return nil, err
	}

	return signedTransaction, nil
}

// SendRawStakingTransaction - send the raw staking tx to the RPC endpoint
func SendRawStakingTransaction(rpcClient *rpc.HTTPMessenger, signature *string) (interface{}, error) {
	reply, err := rpcClient.SendRPC(rpc.Method.SendRawStakingTransaction, []interface{}{signature})
	if err != nil {
		return nil, err
	}

	receiptHash, _ := reply["result"]

	return receiptHash, nil
}

// NumericDecToBigIntAmount - convert a numeric.Dec amount to a converted big.Int amount
func NumericDecToBigIntAmount(amount numeric.Dec) (bigAmount *big.Int) {
	if !amount.IsNil() {
		amount = amount.Mul(transactions.OneAsDec)
		bigAmount = amount.RoundInt()
	}

	return bigAmount
}
