package accounts

import (
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"

	"github.com/btcsuite/btcd/btcec"
	mapset "github.com/deckarep/golang-set"
	goSDKAccount "github.com/harmony-one/go-sdk/pkg/account"
	"github.com/harmony-one/go-sdk/pkg/address"
	goSDKAddress "github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/keys"
	"github.com/harmony-one/go-sdk/pkg/mnemonic"
	goSDKStore "github.com/harmony-one/go-sdk/pkg/store"
	"github.com/harmony-one/harmony/accounts/keystore"
)

// CreateNewLocalAccount - generates a new local keystore account
func (account *Account) CreateNewLocalAccount(candidate *goSDKAccount.Creation) (err error) {
	ks := goSDKStore.FromAccountName(candidate.Name)
	if candidate.Mnemonic == "" {
		candidate.Mnemonic = mnemonic.Generate()
	}

	// Hardcoded index of 0 here.
	privateKey, _ := keys.FromMnemonicSeedAndPassphrase(candidate.Mnemonic, 0)
	acc, err := ks.ImportECDSA(privateKey.ToECDSA(), candidate.Passphrase)
	if err != nil {
		return err
	}

	account.Name = candidate.Name
	account.Address = goSDKAddress.ToBech32(acc.Address)
	account.Passphrase = candidate.Passphrase
	account.Keystore = ks
	account.Account = &acc
	//account.PrivateKey = privateKey
	//account.PublicKey = publicKey

	return nil
}

// ImportFromPrivateKey allows import of an ECDSA private key
func (account *Account) ImportFromPrivateKey(privateKey, name, passphrase string) error {
	privateKey = strings.TrimPrefix(privateKey, "0x")

	if name == "" {
		name = generateName() + "-imported"
		for goSDKStore.DoesNamedAccountExist(name) {
			name = generateName() + "-imported"
		}
	} else if goSDKStore.DoesNamedAccountExist(name) {
		return fmt.Errorf("account %s already exists", name)
	}

	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return err
	}
	if len(privateKeyBytes) != common.Secp256k1PrivateKeyBytesLength {
		return common.ErrBadKeyLength
	}

	// btcec.PrivKeyFromBytes only returns a secret key and public key
	btcecPrivateKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privateKeyBytes)
	ks := goSDKStore.FromAccountName(name)
	acc, err := ks.ImportECDSA(btcecPrivateKey.ToECDSA(), passphrase)
	if err != nil {
		return err
	}

	account.Name = name
	account.Address = goSDKAddress.ToBech32(acc.Address)
	account.Passphrase = passphrase
	account.Keystore = ks
	account.Account = &acc
	//account.PrivateKey = btcecPrivateKey
	//account.PublicKey = btcecPublicKey

	return err
}

// ImportKeyStore imports a keystore along with a password
func (account *Account) ImportKeyStore(keyPath, name, passphrase string) error {
	keyPath, err := filepath.Abs(keyPath)
	if err != nil {
		return err
	}

	keyJSON, readError := ioutil.ReadFile(keyPath)
	if readError != nil {
		return readError
	}

	if name == "" {
		name = generateName() + "-imported"
		for goSDKStore.DoesNamedAccountExist(name) {
			name = generateName() + "-imported"
		}
	} else if goSDKStore.DoesNamedAccountExist(name) {
		return fmt.Errorf("account %s already exists", name)
	}

	key, err := keystore.DecryptKey(keyJSON, passphrase)
	if err != nil {
		return err
	}

	b32 := address.ToBech32(key.Address)
	ks := goSDKStore.FromAddress(b32)
	if ks == nil {
		uDir, _ := homedir.Dir()
		newPath := filepath.Join(uDir, common.DefaultConfigDirName, common.DefaultConfigAccountAliasesDirName, name, filepath.Base(keyPath))

		err = writeToFile(newPath, string(keyJSON))
		if err != nil {
			return err
		}

		account.Keystore, account.Account, err = goSDKStore.UnlockedKeystore(b32, passphrase)
		if err != nil {
			return err
		}
	}

	if account.Account == nil {
		acc, err := ks.ImportECDSA(key.PrivateKey, passphrase)
		if err != nil {
			return err
		}

		account.Account = &acc
	}

	account.Name = name
	account.Address = b32
	account.Passphrase = passphrase

	return nil
}

func generateName() string {
	words := strings.Split(mnemonic.Generate(), " ")
	existingAccounts := mapset.NewSet()
	for a := range goSDKStore.LocalAccounts() {
		existingAccounts.Add(a)
	}
	foundName := false
	acct := ""
	i := 0
	for {
		if foundName {
			break
		}
		if i == len(words)-1 {
			words = strings.Split(mnemonic.Generate(), " ")
		}
		candidate := words[i]
		if !existingAccounts.Contains(candidate) {
			foundName = true
			acct = candidate
			break
		}
	}
	return acct
}

func writeToFile(path string, data string) error {
	currDir, _ := os.Getwd()
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	os.MkdirAll(filepath.Dir(path), 0777)
	os.Chdir(filepath.Dir(path))
	file, err := os.Create(filepath.Base(path))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	os.Chdir(currDir)
	return file.Sync()
}
