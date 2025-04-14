package hdwallet

import (
	"crypto/ecdsa"
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

// DeriveECDSA can be used to derive ECDSA keys for secp256k1 signatures
func (w *Wallet) DeriveECDSA(path DerivationPath) (*ecdsa.PrivateKey, error) {
	masterKey, err := hdkeychain.NewMaster(w.Seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, fmt.Errorf("failed to generate hdkeychain master key: %w", err)
	}
	key := masterKey
	for _, index := range path.Indices() {
		key, err = key.Derive(index)
		if err != nil {
			return nil, fmt.Errorf("failed to derive key: %w", err)
		}
	}
	privateKey, err := key.ECPrivKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get EC private key: %w", err)
	}
	return privateKey.ToECDSA(), nil
}
