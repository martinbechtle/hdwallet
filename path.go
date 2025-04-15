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
	EvmDefaultDerivationPath    = MustParseDerivationPath("m/44'/60'/0'/0/0")
	SolanaDefaultDerivationPath = MustParseDerivationPath("m/44'/501'/0'/0'")

	// Valid path formats:
	// m/44'/60'/0'/0/0
	// m/44h/60h/0h/0/0
	// m/44H/60H/0H/0/0
	// 44'/60'/0'/0/0
	pathRegex = regexp.MustCompile(`^(m/)?(\d+'?H?h?/)*(\d+'?H?h?)$`)
)

// PathComponent represents a single segment in a derivation path
type PathComponent struct {
	Index    uint32
	Hardened bool
}

// DerivationPath represents a complete BIP44 derivation path
type DerivationPath struct {
	Components []PathComponent
}

type BipComponents struct {
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
	path = strings.TrimPrefix(path, "m/")
	components := strings.Split(path, "/")
	if len(components) > maxPathDepth {
		return DerivationPath{}, fmt.Errorf("derivation path too deep (max %d levels)", maxPathDepth)
	}

	result := DerivationPath{
		Components: make([]PathComponent, len(components)),
	}

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
	return result, nil
}

// MustParseDerivationPath same as ParseDerivationPath but panics on error
func MustParseDerivationPath(path string) DerivationPath {
	result, err := ParseDerivationPath(path)
	if err != nil {
		panic(err)
	}
	return result
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

func (p DerivationPath) Indices() []uint32 {
	indices := make([]uint32, len(p.Components))
	for i, component := range p.Components {
		index := component.Index
		if component.Hardened {
			index |= 0x80000000
		}
		indices[i] = index
	}
	return indices
}

func (p DerivationPath) BipComponents() BipComponents {
	result := BipComponents{}
	if len(p.Components) >= 1 {
		if p.Components[0].Hardened {
			result.Purpose = uint8(p.Components[0].Index)
		}
	}
	if len(p.Components) >= 2 {
		if p.Components[1].Hardened {
			result.CoinType = p.Components[1].Index
		}
	}
	if len(p.Components) >= 3 {
		if p.Components[2].Hardened {
			result.Account = p.Components[2].Index
		}
	}
	if len(p.Components) >= 4 {
		result.Change = p.Components[3].Index
	}
	if len(p.Components) >= 5 {
		result.AddressIdx = p.Components[4].Index
	}
	return result
}
