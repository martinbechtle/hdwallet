package hdwallet_test

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"

	"github.com/martinbechtle/hdwallet"
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

func TestWalletDerive_secp256k(t *testing.T) {
	wallet, err := hdwallet.NewWalletFromMnemonic(validMnemonic, "")
	require.NoError(t, err)

	derivedKey, err := wallet.DeriveECDSA(hdwallet.EvmDefaultDerivationPath)
	require.NoError(t, err)

	privateKeyBytes := derivedKey.D.Bytes()
	require.Equal(t, "695e0bbbe84b57d4f7d3d14c0e05a406f42cd73bb15d0161f6261c1ce6ddd1ec", hex.EncodeToString(privateKeyBytes))
}

func TestWalletDerive_ed25519(t *testing.T) {
	t.Run("12 word seed", func(t *testing.T) {
		mnemonic := "elder sign digital common crumble else express festival menu surge price lawsuit"
		wallet, err := hdwallet.NewWalletFromMnemonic(mnemonic, "")
		require.NoError(t, err)
		require.Equal(t, "15572a9ce08615becf26dcc73da42cec3d20853ce1027b65f22050fa439e205321fa5559e076548d8e856538fcf599a245cf7edac4db5d3c93c04dc21131babb", hex.EncodeToString(wallet.Seed))

		derivedKey, err := wallet.DeriveEd25519(hdwallet.SolanaDefaultDerivationPath)
		require.NoError(t, err)

		pubKey, ok := derivedKey.Public().(ed25519.PublicKey)
		require.True(t, ok)
		require.Equal(t, "33b021058c8da2a09734396e93c10403b0213219a679dba6051be68adec3a9ad", hex.EncodeToString(pubKey))
	})
	t.Run("24 word seed", func(t *testing.T) {
		mnemonic := "seminar gadget common sing coral blood turkey quit bike veteran glimpse invite setup million vapor eight left detail donkey gun train olympic sad alone"
		wallet, err := hdwallet.NewWalletFromMnemonic(mnemonic, "")
		require.NoError(t, err)

		derivedKey, err := wallet.DeriveEd25519(hdwallet.SolanaDefaultDerivationPath)
		require.NoError(t, err)

		pubKey, ok := derivedKey.Public().(ed25519.PublicKey)
		require.True(t, ok)
		require.Equal(t, "376024eec8627dced706cc72efd4b4bf4a5050788af96eacbb983a539893b4e0", hex.EncodeToString(pubKey))
	})
}
