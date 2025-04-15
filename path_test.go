package hdwallet_test

import (
	"reflect"
	"testing"

	"github.com/martinbechtle/hdwallet"
	"github.com/stretchr/testify/require"
)

func TestParseDerivationPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    hdwallet.DerivationPath
		wantErr string
	}{
		{
			name:    "invalid format",
			path:    "not-valid",
			wantErr: "invalid derivation path format",
		},
		{
			name:    "invalid index",
			path:    "m/44'/abc'/0'/0/0",
			wantErr: "invalid derivation path format",
		},
		{
			name: "valid BIP44 path",
			path: "m/44'/60'/0'/0/0",
			want: hdwallet.DerivationPath{
				Components: []hdwallet.PathComponent{
					{Index: 44, Hardened: true},
					{Index: 60, Hardened: true},
					{Index: 0, Hardened: true},
					{Index: 0, Hardened: false},
					{Index: 0, Hardened: false},
				},
			},
		},
		{
			name: "valid path without m prefix",
			path: "44'/60'/0'/0/0",
			want: hdwallet.DerivationPath{
				Components: []hdwallet.PathComponent{
					{Index: 44, Hardened: true},
					{Index: 60, Hardened: true},
					{Index: 0, Hardened: true},
					{Index: 0, Hardened: false},
					{Index: 0, Hardened: false},
				},
			},
		},
		{
			name: "hardened with h suffix",
			path: "m/44h/60h/0h/0/0",
			want: hdwallet.DerivationPath{
				Components: []hdwallet.PathComponent{
					{Index: 44, Hardened: true},
					{Index: 60, Hardened: true},
					{Index: 0, Hardened: true},
					{Index: 0, Hardened: false},
					{Index: 0, Hardened: false},
				},
			},
		},
		{
			name: "hardened with H suffix",
			path: "m/44H/60H/0H/0/0",
			want: hdwallet.DerivationPath{
				Components: []hdwallet.PathComponent{
					{Index: 44, Hardened: true},
					{Index: 60, Hardened: true},
					{Index: 0, Hardened: true},
					{Index: 0, Hardened: false},
					{Index: 0, Hardened: false},
				},
			},
		},
		{
			name: "partial path",
			path: "m/44'/60'/0'",
			want: hdwallet.DerivationPath{
				Components: []hdwallet.PathComponent{
					{Index: 44, Hardened: true},
					{Index: 60, Hardened: true},
					{Index: 0, Hardened: true},
				},
			},
		},
		{
			name:    "invalid format",
			path:    "x/44'/60'/0'/0/0",
			wantErr: "invalid derivation path format",
		},
		{
			name:    "too deep",
			path:    "m/1/2/3/4/5/6/7/8/9/10/11",
			wantErr: "derivation path too deep (max 10 levels)",
		},
		{
			name:    "invalid index (int32 overflow)",
			path:    "m/44'/4294967296'/0'/0/0",
			wantErr: "invalid index at position 1: 4294967296",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hdwallet.ParseDerivationPath(tt.path)
			if tt.wantErr == "" && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseDerivationPath() = %v, want %v", got, tt.want)
				return
			}
			if tt.wantErr != "" && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("ParseDerivationPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDerivationPath_String(t *testing.T) {
	tests := []struct {
		name string
		path hdwallet.DerivationPath
		want string
	}{
		{
			name: "empty path",
			path: hdwallet.DerivationPath{Components: []hdwallet.PathComponent{}},
			want: "m/",
		},
		{
			name: "BIP44 path",
			path: hdwallet.DerivationPath{
				Components: []hdwallet.PathComponent{
					{Index: 44, Hardened: true},
					{Index: 60, Hardened: true},
					{Index: 0, Hardened: true},
					{Index: 0, Hardened: false},
					{Index: 0, Hardened: false},
				},
			},
			want: "m/44'/60'/0'/0/0",
		},
		{
			name: "partial path",
			path: hdwallet.DerivationPath{
				Components: []hdwallet.PathComponent{
					{Index: 44, Hardened: true},
					{Index: 60, Hardened: true},
				},
			},
			want: "m/44'/60'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.path.String(); got != tt.want {
				t.Errorf("hdwallet.DerivationPath.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDerivationPath_Indices(t *testing.T) {
	derivationPath, err := hdwallet.ParseDerivationPath("m/44'/60'/0'/0/0")
	require.NoError(t, err)
	indices := derivationPath.Indices()
	require.Equal(t, []uint32{2147483692, 2147483708, 2147483648, 0, 0}, indices) // first three are hardened
}

func TestDerivationPath_BipComponents(t *testing.T) {
	path := hdwallet.DerivationPath{
		Components: []hdwallet.PathComponent{
			{Index: 44, Hardened: true},
			{Index: 60, Hardened: true},
			{Index: 0, Hardened: true},
			{Index: 0, Hardened: false},
			{Index: 0, Hardened: false},
		},
	}
	components := path.BipComponents()
	expected := hdwallet.BipComponents{
		Purpose:    44,
		CoinType:   60,
		Account:    0,
		Change:     0,
		AddressIdx: 0,
	}
	require.Equal(t, expected, components)
}
