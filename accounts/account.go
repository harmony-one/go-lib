package accounts

import (
	"regexp"

	"github.com/btcsuite/btcd/btcec"
	networkTypes "github.com/harmony-one/go-lib/network/types/network"
	goSDKAccount "github.com/harmony-one/go-sdk/pkg/account"
	goSDKAddress "github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/store"
	goSDKStore "github.com/harmony-one/go-sdk/pkg/store"
	hmyAccounts "github.com/harmony-one/harmony/accounts"
	hmyKeystore "github.com/harmony-one/harmony/accounts/keystore"
	"github.com/harmony-one/harmony/numeric"
)

var (
	addressRegex = regexp.MustCompile(`"address":"(?P<address>[a-z0-9]+)"`)
)

// Account - represents an account
type Account struct {
	Name       string `json:"name" yaml:"name"`
	Address    string `json:"address" yaml:"address"`
	Passphrase string `json:"passphrase" yaml:"passphrase"`
	Nonce      uint64
	Balance    numeric.Dec
	Keystore   *hmyKeystore.KeyStore
	Account    *hmyAccounts.Account
	PrivateKey *btcec.PrivateKey
	PublicKey  *btcec.PublicKey
	Unlocked   bool
}

// Unlock - unlocks a given account's keystore
func (account *Account) Unlock() (err error) {
	if !account.Unlocked {
		if account.Keystore == nil || account.Account == nil {
			account.Keystore, account.Account, err = goSDKStore.UnlockedKeystore(account.Address, account.Passphrase)
			if err != nil {
				return err
			}

			account.Unlocked = true
		} else {
			if err := account.Keystore.Unlock(*account.Account, account.Passphrase); err != nil {
				return err
			}

			account.Unlocked = true
		}
	}

	return nil
}

// ExportKeystore - exports a given account/keystore as JSON
func (account *Account) ExportKeystore(passphrase string) ([]byte, error) {
	if err := account.Unlock(); err != nil {
		return nil, err
	}

	for _, acc := range account.Keystore.Accounts() {
		keyJSON, err := account.Keystore.Export(hmyAccounts.Account{Address: acc.Address}, passphrase, passphrase)
		if err != nil {
			return nil, err
		}
		return keyJSON, nil
	}

	return nil, nil
}

// GetAllShardBalances - checks the balances in all shards for a given network, mode and address
func (account *Account) GetAllShardBalances(net *networkTypes.Network) (map[uint32]numeric.Dec, error) {
	return net.GetAllShardBalances(account.Address)
}

// GetShardBalance - gets the balance for a given network, mode, address and shard
func (account *Account) GetShardBalance(net *networkTypes.Network, shardID uint32) (numeric.Dec, error) {
	return net.GetShardBalance(account.Address, shardID)
}

// GetTotalBalance - gets the total balance across all shards for a given network, mode and address
func (account *Account) GetTotalBalance(net *networkTypes.Network) (numeric.Dec, error) {
	return net.GetTotalBalance(account.Address)
}

// DoesNamedAccountExist - wrapper around store.DoesNamedAccountExist(name)
func DoesNamedAccountExist(name string) bool {
	return store.DoesNamedAccountExist(name)
}

// ImportPrivateKeyAccount - imports a private key into the keystore
func ImportPrivateKeyAccount(privateKey string, keyName string, address string, passphrase string) (acc Account, err error) {
	acc = Account{}

	if DoesNamedAccountExist(keyName) {
		tempAcc := FindAccountByName(keyName)
		if tempAcc.Address != "" {
			acc = tempAcc
			acc.Passphrase = passphrase
			return acc, nil
		}
	} else if DoesAddressExistInKeystore(address) {
		tempAcc := FindAccountByAddress(address)
		if tempAcc.Address != "" {
			acc = tempAcc
			acc.Passphrase = passphrase
			return acc, nil
		}
	} else {
		//fmt.Println(fmt.Sprintf("Proceeding to import private key %s - name: %s, address: %s", privateKey, keyName, address))

		err = acc.ImportFromPrivateKey(privateKey, keyName, passphrase)
		if acc.Name == "" || err != nil {
			return acc, err
		}

		//fmt.Println(fmt.Sprintf("Successfully imported private key %s - name %s, address: %s", privateKey, keyName, address))
	}

	return acc, nil
}

// ImportKeystoreAccount - imports an existing keystore file into the keystore
func ImportKeystoreAccount(keyFile string, keyName string, address string, passphrase string) (acc Account, err error) {
	acc = Account{}

	if DoesNamedAccountExist(keyName) {
		tempAcc := FindAccountByName(keyName)
		if tempAcc.Address != "" {
			acc = tempAcc
			acc.Passphrase = passphrase
			return acc, nil
		}
	} else if DoesAddressExistInKeystore(address) {
		tempAcc := FindAccountByAddress(address)
		if tempAcc.Address != "" {
			acc = tempAcc
			acc.Passphrase = passphrase
			return acc, nil
		}
	} else {
		//fmt.Println(fmt.Sprintf("Proceeding to import keyfile with path %s - name: %s, address: %s", keyFile, keyName, address))

		err = acc.ImportKeyStore(keyFile, keyName, passphrase)
		if acc.Name == "" || err != nil {
			return acc, err
		}

		//fmt.Println(fmt.Sprintf("Successfully imported keyfile with path: %s, name %s, address: %s", keyFile, keyName, address))
	}

	return acc, nil
}

// FindAccountByName - finds the account address associated with a given key store name
func FindAccountByName(targetName string) Account {
	for _, name := range store.LocalAccounts() {
		if name == targetName {
			ks := store.FromAccountName(name)
			allAccounts := ks.Accounts()
			for _, account := range allAccounts {
				return Account{
					Name:    name,
					Address: goSDKAddress.ToBech32(account.Address),
				}
			}
		}
	}

	return Account{}
}

// FindAccountByAddress - finds the account name associated with a given address
func FindAccountByAddress(addr string) Account {
	for _, name := range store.LocalAccounts() {
		ks := store.FromAccountName(name)
		allAccounts := ks.Accounts()
		for _, account := range allAccounts {
			formattedAddress := goSDKAddress.ToBech32(account.Address)
			if formattedAddress == addr {
				return Account{
					Name:    name,
					Address: formattedAddress,
				}
			}
		}
	}

	return Account{}
}

// FindAccountAddressByName - finds the account address associated with a given key store name
func FindAccountAddressByName(targetName string) string {
	for _, name := range store.LocalAccounts() {
		if name == targetName {
			ks := store.FromAccountName(name)
			allAccounts := ks.Accounts()
			for _, account := range allAccounts {
				return goSDKAddress.ToBech32(account.Address)
			}
		}
	}

	return ""
}

// FindAccountNameByAddress - finds the account name associated with a given address
func FindAccountNameByAddress(addr string) string {
	for _, name := range store.LocalAccounts() {
		ks := store.FromAccountName(name)
		allAccounts := ks.Accounts()
		for _, account := range allAccounts {
			if goSDKAddress.ToBech32(account.Address) == addr {
				return name
			}
		}
	}

	return ""
}

// DoesAddressExistInKeystore - checks if a given address exists in the keystore
func DoesAddressExistInKeystore(targetAddress string) bool {
	exists := false

	for _, name := range store.LocalAccounts() {
		ks := store.FromAccountName(name)
		allAccounts := ks.Accounts()
		for _, account := range allAccounts {
			if targetAddress == goSDKAddress.ToBech32(account.Address) {
				return true
			}
		}
	}

	return exists
}

// GenerateAccount - generates a new account using the specified name and passphrase
func GenerateAccount(name string, passphrase string) (acc Account, err error) {
	acc = Account{}
	accountExists := store.DoesNamedAccountExist(name)

	if !accountExists {
		accCreation := goSDKAccount.Creation{
			Name:            name,
			Passphrase:      passphrase,
			Mnemonic:        "",
			HdAccountNumber: nil,
			HdIndexNumber:   nil,
		}

		err = acc.CreateNewLocalAccount(&accCreation)
		if err != nil {
			return acc, err
		}
	}

	return acc, nil
}
