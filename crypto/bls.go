package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"

	bls_core "github.com/harmony-one/bls/ffi/go/bls"
	"github.com/harmony-one/harmony/crypto/bls"
	"github.com/harmony-one/harmony/crypto/hash"
	"github.com/harmony-one/harmony/staking/types"
)

// BLSKey - represents a BLS key
type BLSKey struct {
	PrivateKey    *bls_core.SecretKey
	PrivateKeyHex string
	PublicKey     *bls_core.PublicKey
	PublicKeyHex  string

	ShardPublicKey *bls.SerializedPublicKey
	ShardSignature *bls.SerializedSignature
}

// GenerateBlsKey - generates a new bls key and returns its private and public keys as hex strings
func GenerateBlsKey(message string) (BLSKey, error) {
	privateKey := bls.RandPrivateKey()
	privateKeyHex := privateKey.SerializeToHexStr()

	publicKey := privateKey.GetPublicKey()
	publicKeyHex := publicKey.SerializeToHexStr()

	key := BLSKey{
		PrivateKey:    privateKey,
		PrivateKeyHex: privateKeyHex,
		PublicKey:     publicKey,
		PublicKeyHex:  publicKeyHex,
	}

	err := key.Initialize(message)
	if err != nil {
		return BLSKey{}, err
	}

	return key, nil
}

// Initialize - generates a BLS Signature based on a given BLS key
func (blsKey *BLSKey) Initialize(message string) error {
	if err := blsKey.AssignShardSignature(message); err != nil {
		return err
	}

	if err := blsKey.AssignShardPublicKey(); err != nil {
		return err
	}

	return nil
}

// AssignShardSignature - signs a given message using the BLSKey and assigns ShardSignature
func (blsKey *BLSKey) AssignShardSignature(message string) error {
	var sig bls.SerializedSignature

	if message == "" {
		message = types.BLSVerificationStr
	}

	msgHash := hash.Keccak256([]byte(message))
	signature := blsKey.PrivateKey.SignHash(msgHash[:])

	bytes := signature.Serialize()
	if len(bytes) != bls.BLSSignatureSizeInBytes {
		return errors.New("bls key length is not 96 bytes")
	}

	copy(sig[:], bytes)
	blsKey.ShardSignature = &sig

	return nil
}

// AssignShardPublicKey - converts a regular pub key to a shardPubKey and assigns ShardPublicKey
func (blsKey *BLSKey) AssignShardPublicKey() error {
	shardPubKey := new(bls.SerializedPublicKey)
	err := shardPubKey.FromLibBLSPublicKey(blsKey.PublicKey)
	if err != nil {
		return errors.New("couldn't convert bls.PublicKey -> shard.BLSPublicKey")
	}

	blsKey.ShardPublicKey = shardPubKey

	return nil
}

// Encrypt - encrypts a BLSKey with a given passphrase
func (blsKey *BLSKey) Encrypt(passphrase string) (string, error) {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(blsKey.PrivateKeyHex), nil)

	return hex.EncodeToString(ciphertext), nil
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
