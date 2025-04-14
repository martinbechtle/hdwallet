package hdwallet

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	// maxPathDepth is the maximum allowed depth for a derivation path
	maxPathDepth = 10
)

var (
	// Valid path formats:
	// m/44'/60'/0'/0/0
	// m/44h/60h/0h/0/0
	// m/44H/60H/0H/0/0
	// 44'/60'/0'/0/0
	pathRegex = regexp.MustCompile(`^(m/)?(\d+'?H?h?/)*(\d+'?H?h?)$`)
)

// TODO verify this regexp

// PathComponent represents a single segment in a derivation path
type PathComponent struct {
	Index    uint32
	Hardened bool
}

// DerivationPath represents a complete BIP44 derivation path
type DerivationPath struct {
	Components []PathComponent
	Purpose    uint8  // BIP44 = 44, BIP49 = 49, BIP84 = 84, etc.
	CoinType   uint32 // Bitcoin = 0, Ethereum = 60, etc.
	Account    uint32
	Change     uint32 // 0 = external, 1 = internal (change)
	AddressIdx uint32
}

// ParseDerivationPath parses a derivation path string with a maximum depth specified by the maxPathDepth constant.
// Returns an error if the path is not formatted correctly or exceeds the max depth.
// Supports ', h and H as hardening suffixes.
// Example of valid derivation path: m/44'/60'/0'/0
func ParseDerivationPath(path string) (DerivationPath, error) {
	if !pathRegex.MatchString(path) {
		return DerivationPath{}, fmt.Errorf("invalid derivation path format")
	}

	// Remove the leading "m/" if present
	if strings.HasPrefix(path, "m/") {
		path = path[2:]
	}
	components := strings.Split(path, "/")
	if len(components) > maxPathDepth {
		return DerivationPath{}, fmt.Errorf("derivation path too deep (max %d levels)", maxPathDepth)
	}

	result := DerivationPath{
		Components: make([]PathComponent, len(components)),
	}
	// TODO double check this logic against known libraries and unit test

	for i, part := range components {
		// Check for hardened component (suffixed with ', h, or H)
		hardened := false
		if strings.HasSuffix(part, "'") || strings.HasSuffix(part, "h") || strings.HasSuffix(part, "H") {
			hardened = true
			part = part[:len(part)-1]
		}

		index, err := strconv.ParseUint(part, 10, 32)
		if err != nil {
			return DerivationPath{}, fmt.Errorf("invalid index at position %d: %s", i, part)
		}
		result.Components[i] = PathComponent{
			Index:    uint32(index),
			Hardened: hardened,
		}
	}

	// Extract BIP44 components if available
	// TODO doesn't seem right... does it mean we should limit the depth?
	if len(result.Components) >= 1 {
		if result.Components[0].Hardened {
			result.Purpose = uint8(result.Components[0].Index)
		}
	}
	if len(result.Components) >= 2 {
		if result.Components[1].Hardened {
			result.CoinType = result.Components[1].Index
		}
	}
	if len(result.Components) >= 3 {
		if result.Components[2].Hardened {
			result.Account = result.Components[2].Index
		}
	}
	if len(result.Components) >= 4 {
		result.Change = result.Components[3].Index
	}
	if len(result.Components) >= 5 {
		result.AddressIdx = result.Components[4].Index
	}

	return result, nil
}

// String returns the string representation of the derivation path
func (p DerivationPath) String() string {
	parts := make([]string, len(p.Components))
	for i, component := range p.Components {
		if component.Hardened {
			parts[i] = fmt.Sprintf("%d'", component.Index)
		} else {
			parts[i] = fmt.Sprintf("%d", component.Index)
		}
	}
	return "m/" + strings.Join(parts, "/")
}

// TODO probably best to have constants for these? Add EVM and SVM to begin with

// NewBIP44Path creates a standard BIP44 path
func NewBIP44Path(coinType, account, change, addressIdx uint32) DerivationPath {
	return DerivationPath{
		Components: []PathComponent{
			{Index: 44, Hardened: true},
			{Index: coinType, Hardened: true},
			{Index: account, Hardened: true},
			{Index: change, Hardened: false},
			{Index: addressIdx, Hardened: false},
		},
		Purpose:    44,
		CoinType:   coinType,
		Account:    account,
		Change:     change,
		AddressIdx: addressIdx,
	}
}

// NewSolanaPath creates a Solana-style path (all hardened)
func NewSolanaPath(account, change, addressIdx uint32) DerivationPath {
	return DerivationPath{
		Components: []PathComponent{
			{Index: 44, Hardened: true},
			{Index: 501, Hardened: true}, // Solana coin type
			{Index: account, Hardened: true},
			{Index: change, Hardened: true},
			{Index: addressIdx, Hardened: true},
		},
		Purpose:    44,
		CoinType:   501,
		Account:    account,
		Change:     change,
		AddressIdx: addressIdx,
	}
}

// IsValid checks if the derivation path follows BIP44 standards
func (p DerivationPath) IsValid() bool {
	// Basic validation
	if len(p.Components) == 0 {
		return true // Master key is valid
	}

	// Check if purpose is hardened
	if len(p.Components) >= 1 && !p.Components[0].Hardened {
		return false
	}

	// Check if coin type is hardened
	if len(p.Components) >= 2 && !p.Components[1].Hardened {
		return false
	}

	// Check if account is hardened
	if len(p.Components) >= 3 && !p.Components[2].Hardened {
		return false
	}

	// Further validation based on specific BIP44 rules could be added here
	return true
}

// IsBIP44 checks if this is a standard BIP44 path
func (p DerivationPath) IsBIP44() bool {
	if len(p.Components) != 5 {
		return false
	}

	// m/44'/coinType'/account'/change/address
	return p.Components[0].Index == 44 &&
		p.Components[0].Hardened &&
		p.Components[1].Hardened &&
		p.Components[2].Hardened &&
		!p.Components[3].Hardened &&
		!p.Components[4].Hardened
}

// RequiresHardenedDerivation checks if this path requires curves that support
// hardened derivation (like ed25519)
func (p DerivationPath) RequiresHardenedDerivation() bool {
	for _, component := range p.Components {
		if component.Hardened {
			return true
		}
	}
	return false
}

// AllComponentsHardened checks if all components in the path are hardened
func (p DerivationPath) AllComponentsHardened() bool {
	for _, component := range p.Components {
		if !component.Hardened {
			return false
		}
	}
	return true
}
