package hdwallet_test

import (
	"encoding/hex"
	"testing"

	"github.com/MartinBechtle/hdwallet"
	"github.com/stretchr/testify/require"
)

var (
	validMnemonic = "night faint enjoy portion appear movie busy waste report circle giant hungry"
)

func TestNewWallet(t *testing.T) {
	_, err := hdwallet.NewWalletFromMnemonic("invalid mnemonic", "")
	require.EqualError(t, err, "invalid mnemonic")

	_, err = hdwallet.NewWalletFromMnemonic(validMnemonic, "")
	require.NoError(t, err)
}

func TestWalletDerive_EVM(t *testing.T) {
	wallet, err := hdwallet.NewWalletFromMnemonic(validMnemonic, "")
	require.NoError(t, err)

	derivedKey, err := wallet.DeriveECDSA(hdwallet.EvmDefaultDerivationPath)
	require.NoError(t, err)

	privateKeyBytes := derivedKey.D.Bytes()
	require.Equal(t, "695e0bbbe84b57d4f7d3d14c0e05a406f42cd73bb15d0161f6261c1ce6ddd1ec", hex.EncodeToString(privateKeyBytes))
}
