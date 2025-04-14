package hdwallet

import (
	"reflect"
	"testing"
)

func TestParseDerivationPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    DerivationPath
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
			want: DerivationPath{
				Components: []PathComponent{
					{Index: 44, Hardened: true},
					{Index: 60, Hardened: true},
					{Index: 0, Hardened: true},
					{Index: 0, Hardened: false},
					{Index: 0, Hardened: false},
				},
				Purpose:    44,
				CoinType:   60,
				Account:    0,
				Change:     0,
				AddressIdx: 0,
			},
		},
		{
			name: "valid path without m prefix",
			path: "44'/60'/0'/0/0",
			want: DerivationPath{
				Components: []PathComponent{
					{Index: 44, Hardened: true},
					{Index: 60, Hardened: true},
					{Index: 0, Hardened: true},
					{Index: 0, Hardened: false},
					{Index: 0, Hardened: false},
				},
				Purpose:    44,
				CoinType:   60,
				Account:    0,
				Change:     0,
				AddressIdx: 0,
			},
		},
		{
			name: "hardened with h suffix",
			path: "m/44h/60h/0h/0/0",
			want: DerivationPath{
				Components: []PathComponent{
					{Index: 44, Hardened: true},
					{Index: 60, Hardened: true},
					{Index: 0, Hardened: true},
					{Index: 0, Hardened: false},
					{Index: 0, Hardened: false},
				},
				Purpose:    44,
				CoinType:   60,
				Account:    0,
				Change:     0,
				AddressIdx: 0,
			},
		},
		{
			name: "hardened with H suffix",
			path: "m/44H/60H/0H/0/0",
			want: DerivationPath{
				Components: []PathComponent{
					{Index: 44, Hardened: true},
					{Index: 60, Hardened: true},
					{Index: 0, Hardened: true},
					{Index: 0, Hardened: false},
					{Index: 0, Hardened: false},
				},
				Purpose:    44,
				CoinType:   60,
				Account:    0,
				Change:     0,
				AddressIdx: 0,
			},
		},
		{
			name: "partial path",
			path: "m/44'/60'/0'",
			want: DerivationPath{
				Components: []PathComponent{
					{Index: 44, Hardened: true},
					{Index: 60, Hardened: true},
					{Index: 0, Hardened: true},
				},
				Purpose:  44,
				CoinType: 60,
				Account:  0,
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
			got, err := ParseDerivationPath(tt.path)
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
		path DerivationPath
		want string
	}{
		{
			name: "empty path",
			path: DerivationPath{Components: []PathComponent{}},
			want: "m/",
		},
		{
			name: "BIP44 path",
			path: DerivationPath{
				Components: []PathComponent{
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
			path: DerivationPath{
				Components: []PathComponent{
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
				t.Errorf("DerivationPath.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDerivationPath_IsValid(t *testing.T) {
	tests := []struct {
		name string
		path DerivationPath
		want bool
	}{
		{
			name: "empty path",
			path: DerivationPath{Components: []PathComponent{}},
			want: true,
		},
		{
			name: "valid BIP44 path",
			path: DerivationPath{
				Components: []PathComponent{
					{Index: 44, Hardened: true},
					{Index: 60, Hardened: true},
					{Index: 0, Hardened: true},
					{Index: 0, Hardened: false},
					{Index: 0, Hardened: false},
				},
			},
			want: true,
		},
		{
			name: "invalid - purpose not hardened",
			path: DerivationPath{
				Components: []PathComponent{
					{Index: 44, Hardened: false},
					{Index: 60, Hardened: true},
					{Index: 0, Hardened: true},
					{Index: 0, Hardened: false},
					{Index: 0, Hardened: false},
				},
			},
			want: false,
		},
		{
			name: "invalid - coin type not hardened",
			path: DerivationPath{
				Components: []PathComponent{
					{Index: 44, Hardened: true},
					{Index: 60, Hardened: false},
					{Index: 0, Hardened: true},
					{Index: 0, Hardened: false},
					{Index: 0, Hardened: false},
				},
			},
			want: false,
		},
		{
			name: "invalid - account not hardened",
			path: DerivationPath{
				Components: []PathComponent{
					{Index: 44, Hardened: true},
					{Index: 60, Hardened: true},
					{Index: 0, Hardened: false},
					{Index: 0, Hardened: false},
					{Index: 0, Hardened: false},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.path.IsValid(); got != tt.want {
				t.Errorf("DerivationPath.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDerivationPath_IsBIP44(t *testing.T) {
	tests := []struct {
		name string
		path DerivationPath
		want bool
	}{
		{
			name: "valid BIP44",
			path: DerivationPath{
				Components: []PathComponent{
					{Index: 44, Hardened: true},
					{Index: 60, Hardened: true},
					{Index: 0, Hardened: true},
					{Index: 0, Hardened: false},
					{Index: 0, Hardened: false},
				},
			},
			want: true,
		},
		{
			name: "wrong purpose",
			path: DerivationPath{
				Components: []PathComponent{
					{Index: 49, Hardened: true},
					{Index: 60, Hardened: true},
					{Index: 0, Hardened: true},
					{Index: 0, Hardened: false},
					{Index: 0, Hardened: false},
				},
			},
			want: false,
		},
		{
			name: "wrong hardening pattern",
			path: DerivationPath{
				Components: []PathComponent{
					{Index: 44, Hardened: true},
					{Index: 60, Hardened: true},
					{Index: 0, Hardened: true},
					{Index: 0, Hardened: true}, // Should be non-hardened
					{Index: 0, Hardened: false},
				},
			},
			want: false,
		},
		{
			name: "too short",
			path: DerivationPath{
				Components: []PathComponent{
					{Index: 44, Hardened: true},
					{Index: 60, Hardened: true},
					{Index: 0, Hardened: true},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.path.IsBIP44(); got != tt.want {
				t.Errorf("DerivationPath.IsBIP44() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDerivationPath_RequiresHardenedDerivation(t *testing.T) {
	tests := []struct {
		name string
		path DerivationPath
		want bool
	}{
		{
			name: "with hardened components",
			path: DerivationPath{
				Components: []PathComponent{
					{Index: 44, Hardened: true},
					{Index: 60, Hardened: true},
					{Index: 0, Hardened: true},
					{Index: 0, Hardened: false},
					{Index: 0, Hardened: false},
				},
			},
			want: true,
		},
		{
			name: "no hardened components",
			path: DerivationPath{
				Components: []PathComponent{
					{Index: 44, Hardened: false},
					{Index: 60, Hardened: false},
					{Index: 0, Hardened: false},
					{Index: 0, Hardened: false},
					{Index: 0, Hardened: false},
				},
			},
			want: false,
		},
		{
			name: "empty path",
			path: DerivationPath{Components: []PathComponent{}},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.path.RequiresHardenedDerivation(); got != tt.want {
				t.Errorf("DerivationPath.RequiresHardenedDerivation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDerivationPath_AllComponentsHardened(t *testing.T) {
	tests := []struct {
		name string
		path DerivationPath
		want bool
	}{
		{
			name: "all hardened",
			path: DerivationPath{
				Components: []PathComponent{
					{Index: 44, Hardened: true},
					{Index: 501, Hardened: true},
					{Index: 0, Hardened: true},
					{Index: 0, Hardened: true},
					{Index: 0, Hardened: true},
				},
			},
			want: true,
		},
		{
			name: "mixed hardening",
			path: DerivationPath{
				Components: []PathComponent{
					{Index: 44, Hardened: true},
					{Index: 60, Hardened: true},
					{Index: 0, Hardened: true},
					{Index: 0, Hardened: false},
					{Index: 0, Hardened: false},
				},
			},
			want: false,
		},
		{
			name: "empty path",
			path: DerivationPath{Components: []PathComponent{}},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.path.AllComponentsHardened(); got != tt.want {
				t.Errorf("DerivationPath.AllComponentsHardened() = %v, want %v", got, tt.want)
			}
		})
	}
}
