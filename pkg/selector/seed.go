package selector

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
)

const SeedBytes = 32 // 256 bits

// GenerateSeed generates a random 256-bit (32 byte) seed from /dev/urandom
func GenerateSeed() ([]byte, error) {
	seed := make([]byte, SeedBytes)
	f, err := os.Open("/dev/urandom")
	if err != nil {
		return nil, fmt.Errorf("cannot open /dev/urandom: %w", err)
	}
	defer f.Close()

	n, err := f.Read(seed)
	if err != nil {
		return nil, fmt.Errorf("cannot read from /dev/urandom: %w", err)
	}
	if n != SeedBytes {
		return nil, fmt.Errorf("insufficient random data: got %d bytes, need %d", n, SeedBytes)
	}

	return seed, nil
}

// ParseSeed converts a hex string to a seed byte slice and validates the length
func ParseSeed(hexString string) ([]byte, error) {
	seed, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, fmt.Errorf("invalid seed format: must be hexadecimal: %w", err)
	}

	if len(seed) != SeedBytes {
		return nil, fmt.Errorf("seed must be %d hex characters (%d bytes), got %d characters (%d bytes)",
			SeedBytes*2, SeedBytes, len(hexString), len(seed))
	}

	return seed, nil
}

// SeedToInt64 converts a seed byte slice to an int64 for use with math/rand
// Uses SHA-256 hash to ensure good distribution even if seed has patterns
func SeedToInt64(seed []byte) int64 {
	hash := sha256.Sum256(seed)
	return int64(binary.BigEndian.Uint64(hash[:8]))
}
