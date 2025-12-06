//go:build zksnark

package zkproof

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/sha2"
	"github.com/consensys/gnark/std/math/uints"
)

// SeedCommitmentCircuit represents the zkSNARK circuit for proving knowledge of a seed
// The circuit proves: "I know a seed whose SHA-256 hash, combined with participants hash, equals this commitment"
// Specifically: commitment = SHA256(hex(seed) + hex(participantsHash))
type SeedCommitmentCircuit struct {
	// Seed is the private witness (256 bits = 32 bytes)
	// This is the secret input that the prover knows but doesn't reveal
	Seed [32]frontend.Variable `gnark:",secret"`

	// ParticipantsHash is the public input (SHA-256 hash of participants file content)
	// This allows verifying which participants list was used
	ParticipantsHash [32]frontend.Variable `gnark:",public"`

	// Commitment is the public input (SHA-256 hash of hex(seed) + hex(participantsHash))
	// This is what everyone can see and verify against
	Commitment [32]frontend.Variable `gnark:",public"`
}

// Define declares the circuit constraints
// This method is called by gnark to build the constraint system
func (circuit *SeedCommitmentCircuit) Define(api frontend.API) error {
	// Initialize uints API for byte operations
	uapi, err := uints.New[uints.U32](api)
	if err != nil {
		return err
	}

	// Step 1: Convert seed (32 bytes) to hex string (64 bytes)
	seedHex := make([]uints.U8, 64)
	for i := 0; i < 32; i++ {
		b := uapi.ByteValueOf(circuit.Seed[i])
		// High nibble: b >> 4
		highNibble := api.Div(b.Val, 16)
		seedHex[i*2] = byteToHexChar(api, highNibble)
		// Low nibble: b - (highNibble * 16)
		lowNibble := api.Sub(b.Val, api.Mul(highNibble, 16))
		seedHex[i*2+1] = byteToHexChar(api, lowNibble)
	}

	// Step 2: Convert participantsHash (32 bytes) to hex string (64 bytes)
	participantsHex := make([]uints.U8, 64)
	for i := 0; i < 32; i++ {
		b := uapi.ByteValueOf(circuit.ParticipantsHash[i])
		// High nibble: b >> 4
		highNibble := api.Div(b.Val, 16)
		participantsHex[i*2] = byteToHexChar(api, highNibble)
		// Low nibble: b - (highNibble * 16)
		lowNibble := api.Sub(b.Val, api.Mul(highNibble, 16))
		participantsHex[i*2+1] = byteToHexChar(api, lowNibble)
	}

	// Step 3: Concatenate hex strings (64 + 64 = 128 bytes)
	combined := make([]uints.U8, 128)
	copy(combined[:64], seedHex)
	copy(combined[64:], participantsHex)

	// Step 4: Compute SHA-256 hash of the concatenated hex string
	hasher, err := sha2.New(api)
	if err != nil {
		return err
	}
	hasher.Write(combined)
	hash := hasher.Sum()

	// Step 5: Assert that the computed hash equals the public commitment
	// This is the main constraint: SHA256(hex(seed) + hex(participantsHash)) == commitment
	for i := 0; i < 32; i++ {
		api.AssertIsEqual(hash[i].Val, circuit.Commitment[i])
	}

	return nil
}

// byteToHexChar converts a nibble (0-15) to its ASCII hex character in-circuit
// 0-9 → '0'-'9' (ASCII 48-57)
// 10-15 → 'a'-'f' (ASCII 97-102)
func byteToHexChar(api frontend.API, nibbleVar frontend.Variable) uints.U8 {
	// Check if nibble < 10
	// If nibble < 10: result = nibble + '0' (48)
	// If nibble >= 10: result = nibble - 10 + 'a' (87)

	// Create a boolean: isDigit = (nibble < 10)
	isDigit := api.IsZero(api.Sub(api.Cmp(nibbleVar, 10), 1))

	// Compute both possibilities:
	// digitChar = nibble + 48
	digitChar := api.Add(nibbleVar, 48)
	// letterChar = nibble + 87 (which is 'a' - 10)
	letterChar := api.Add(nibbleVar, 87)

	// Select based on isDigit: if isDigit then digitChar else letterChar
	result := api.Select(isDigit, digitChar, letterChar)

	return uints.U8{Val: result}
}

// NewCircuit creates a new SeedCommitmentCircuit instance
func NewCircuit() *SeedCommitmentCircuit {
	return &SeedCommitmentCircuit{}
}
