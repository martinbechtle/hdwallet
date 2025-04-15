package hdwallet

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"encoding/binary"
	"fmt"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/tyler-smith/go-bip39"
)

type Wallet struct {
	Seed []byte
}

// NewWallet creates a new keychain that allows deriving keys from a seed.
// Don't forget to run Wallet.Erase after using, for security.
func NewWallet(seed []byte) *Wallet {
	return &Wallet{
		Seed: seed,
	}
}

// NewWalletFromMnemonic creates a new keychain that allows deriving keys from a seed.
// The seed is determined by parsing a bip39 mnemonic and an optional passphrase.
// Don't forget to run Wallet.Erase after using, for security.
func NewWalletFromMnemonic(mnemonic string, passphrase string) (*Wallet, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic")
	}
	seed := bip39.NewSeed(mnemonic, passphrase)
	return NewWallet(seed), nil
}

// Erase can be used for security reasons in apps that need to instantiate a wallet temporarily and then avoid keeping
// sensitive bytes in memory.
func (w *Wallet) Erase() {
	eraseBytes(w.Seed)
}

// DeriveECDSA can be used to deriveEd25519Child ECDSA keys for secp256k1 signatures
func (w *Wallet) DeriveECDSA(path DerivationPath) (*ecdsa.PrivateKey, error) {
	masterKey, err := hdkeychain.NewMaster(w.Seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, fmt.Errorf("failed to generate hdkeychain master key: %w", err)
	}
	key := masterKey
	for _, index := range path.Indices() {
		key, err = key.Derive(index)
		if err != nil {
			return nil, fmt.Errorf("failed to deriveEd25519Child key: %w", err)
		}
	}
	privateKey, err := key.ECPrivKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get EC private key: %w", err)
	}
	return privateKey.ToECDSA(), nil
}

func (w *Wallet) DeriveEd25519(path DerivationPath) (ed25519.PrivateKey, error) {
	if len(w.Seed) != 64 {
		return nil, fmt.Errorf("seed must be 64 bytes long")
	}
	sum := sha512hmac([]byte("ed25519 seed"), w.Seed)
	derivedSeed := sum[:32]
	chain := sum[32:]

	for _, index := range path.Indices() {
		derivedSeed, chain = deriveEd25519Child(derivedSeed, chain, index)
	}
	return ed25519.NewKeyFromSeed(derivedSeed), nil
}

func deriveEd25519Child(parentKey []byte, parentChainCode []byte, index uint32) ([]byte, []byte) {
	// Data to be hashed: 0x00 || parentKey || index (big endian)
	data := make([]byte, 1+32+4)
	data[0] = 0x00 // Important: leading zero byte
	copy(data[1:], parentKey)
	binary.BigEndian.PutUint32(data[1+32:], index)

	// Calculate HMAC hash and split the result
	I := sha512hmac(parentChainCode, data)
	IL := I[:32]
	IR := I[32:]
	return IL, IR
}
