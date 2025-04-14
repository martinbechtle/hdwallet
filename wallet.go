package hdwallet

import (
	"fmt"
)

type Curve string

const (
	CurveSecp256k1 Curve = "secp256k1"
	CurveEd25519   Curve = "ed25519"
)

type Wallet struct {
	Seed []byte
}

// Erase can be used for security reasons in apps that need to instantiate a wallet temporarily and then avoid keeping
// sensitive bytes in memory
func (w *Wallet) Erase() {
	eraseBytes(w.Seed)
}

func (w *Wallet) Derive(seed []byte, coinType uint32, curve string) (privateKey []byte, err error) {
	switch curve {
	case "secp256k1":
		// BIP44 path: m/44'/coinType'/0'/0/0
		//path := fmt.Sprintf("m/44'/%d'/0'/0/0", coinType)
		//return DeriveSecp256k1Key(seed, path)
		return nil, fmt.Errorf("secp256k1 not supported yet")

	case "ed25519":
		// BIP44 path for ed25519: m/44'/coinType'/0'/0'/0'
		// Note: All segments are hardened
		//path := fmt.Sprintf("m/44'/%d'/0'/0'/0'", coinType)
		//return DeriveEd25519Key(seed, path)
		return nil, fmt.Errorf("ed25519 not supported yet")

	default:
		return nil, fmt.Errorf("unsupported curve: %s", curve)
	}
}
