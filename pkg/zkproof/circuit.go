//go:build zksnark

package zkproof

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/sha2"
	"github.com/consensys/gnark/std/math/uints"
)

// SeedCommitmentCircuit represents the zkSNARK circuit for proving knowledge of a seed
// The circuit proves: "I know a seed whose SHA-256 hash equals this commitment"
type SeedCommitmentCircuit struct {
	// Seed is the private witness (256 bits = 32 bytes)
	// This is the secret input that the prover knows but doesn't reveal
	Seed [32]frontend.Variable `gnark:",secret"`

	// Commitment is the public input (SHA-256 hash of the seed)
	// This is what everyone can see and verify against
	Commitment [32]frontend.Variable `gnark:",public"`
}

// Define declares the circuit constraints
// This method is called by gnark to build the constraint system
func (circuit *SeedCommitmentCircuit) Define(api frontend.API) error {
	// Initialize SHA-256 hasher
	// This creates the circuit gadget for computing SHA-256
	hasher, err := sha2.New(api)
	if err != nil {
		return err
	}

	// Convert seed to U8 array for hashing
	// The SHA-256 hasher expects []uints.U8
	uapi, err := uints.New[uints.U32](api)
	if err != nil {
		return err
	}

	seedBytes := make([]uints.U8, 32)
	for i := 0; i < 32; i++ {
		seedBytes[i] = uapi.ByteValueOf(circuit.Seed[i])
	}

	// Compute SHA-256 hash of the seed in-circuit
	// This is where most of the constraints come from (~50,000 constraints)
	hasher.Write(seedBytes)
	hash := hasher.Sum()

	// Assert that the computed hash equals the public commitment
	// This is the main constraint: hash(seed) == commitment
	for i := 0; i < 32; i++ {
		api.AssertIsEqual(hash[i].Val, circuit.Commitment[i])
	}

	return nil
}

// NewCircuit creates a new SeedCommitmentCircuit instance
func NewCircuit() *SeedCommitmentCircuit {
	return &SeedCommitmentCircuit{}
}
