package selector

import (
	"testing"
)

func TestGenerateSeed(t *testing.T) {
	seed, err := GenerateSeed()
	if err != nil {
		t.Fatalf("GenerateSeed() error = %v", err)
	}

	if len(seed) != SeedBytes {
		t.Errorf("GenerateSeed() returned %d bytes, want %d", len(seed), SeedBytes)
	}

	// Test that multiple calls produce different seeds
	seed2, err := GenerateSeed()
	if err != nil {
		t.Fatalf("GenerateSeed() second call error = %v", err)
	}

	// Compare as strings to get readable output if they're the same
	if string(seed) == string(seed2) {
		t.Error("GenerateSeed() produced identical seeds on consecutive calls (very unlikely)")
	}
}

func TestParseSeed_Valid(t *testing.T) {
	// Valid 64-character hex string (32 bytes)
	validHex := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	seed, err := ParseSeed(validHex)
	if err != nil {
		t.Fatalf("ParseSeed() with valid hex error = %v", err)
	}

	if len(seed) != SeedBytes {
		t.Errorf("ParseSeed() returned %d bytes, want %d", len(seed), SeedBytes)
	}
}

func TestParseSeed_InvalidHex(t *testing.T) {
	tests := []struct {
		name    string
		hexStr  string
		wantErr bool
	}{
		{
			name:    "non-hex characters",
			hexStr:  "gggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg",
			wantErr: true,
		},
		{
			name:    "odd number of hex characters",
			hexStr:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcde",
			wantErr: true,
		},
		{
			name:    "spaces in hex",
			hexStr:  "0123456789abcdef 123456789abcdef0123456789abcdef0123456789abcdef",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseSeed(tt.hexStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSeed() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseSeed_WrongLength(t *testing.T) {
	tests := []struct {
		name   string
		hexStr string
	}{
		{
			name:   "too short",
			hexStr: "0123456789abcdef",
		},
		{
			name:   "too long",
			hexStr: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef00",
		},
		{
			name:   "empty",
			hexStr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseSeed(tt.hexStr)
			if err == nil {
				t.Error("ParseSeed() expected error for wrong length, got nil")
			}
		})
	}
}

func TestSeedToInt64(t *testing.T) {
	seed := make([]byte, SeedBytes)
	for i := range seed {
		seed[i] = byte(i)
	}

	val1 := SeedToInt64(seed)

	// Same seed should produce same int64
	val2 := SeedToInt64(seed)
	if val1 != val2 {
		t.Errorf("SeedToInt64() not deterministic: %d != %d", val1, val2)
	}

	// Different seed should (very likely) produce different int64
	seed[0] = 255
	val3 := SeedToInt64(seed)
	if val1 == val3 {
		t.Error("SeedToInt64() produced same value for different seeds")
	}
}

func TestSeedToInt64_AllZeros(t *testing.T) {
	seed := make([]byte, SeedBytes)
	// All zeros is a valid seed
	val := SeedToInt64(seed)

	// Should be deterministic
	val2 := SeedToInt64(seed)
	if val != val2 {
		t.Errorf("SeedToInt64() not deterministic for all-zero seed: %d != %d", val, val2)
	}
}
