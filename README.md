# Crypto HDWallet

Small utility for Go crypto projects.

Historically, ECDSA key derivation for secp256k has been available through [btcd](https://github.com/btcsuite/btcd), but
the ecosystem has been lacking good support for ed25519 key derivation.

This library adds support for both, which should allow Go projects to support key derivation for most blockchain
protocols out there.

## Usage example

### EVM

```go
package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"log"

	"github.com/martinbechtle/hdwallet"
)

func main() {
	mnemonic := "elder sign digital common crumble else express festival menu surge price lawsuit"
	wallet, err := hdwallet.NewWalletFromMnemonic(mnemonic, "")
	if err != nil {
		log.Fatal(err)
	}
	derivedKey, err := wallet.DeriveECDSA(hdwallet.EvmDefaultDerivationPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(hex.EncodeToString(derivedKey.D.Bytes()))
}
```

You can then get the public portion of the key and use ethereum-specific libraries to infer the address and produce
signatures.

### Solana

```go
derivedKey, err := wallet.DeriveEd25519(hdwallet.SolanaDefaultDerivationPath)
if err != nil {
    log.Fatal(err)
}
log.Println(hex.EncodeToString(derivedKey))
```
