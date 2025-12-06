//go:build zksnark

package zkproof

import (
	"crypto/sha256"
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/consensys/gnark/test/unsafekzg"
)

// SetupKeysWithCS contains the proving and verification keys plus the constraint system
type SetupKeysWithCS struct {
	ProvingKey       plonk.ProvingKey
	VerifyingKey     plonk.VerifyingKey
	ConstraintSystem constraint.ConstraintSystem
}

// Setup performs the trusted setup for the PLONK proving system
// This needs to be done once for the circuit
// Returns the proving and verification keys
func Setup() (*SetupKeysWithCS, error) {
	// Create the circuit definition
	circuit := NewCircuit()

	// Compile the circuit into a constraint system (SCS = Sparse Constraint System)
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), scs.NewBuilder, circuit)
	if err != nil {
		return nil, fmt.Errorf("failed to compile circuit: %w", err)
	}

	// Generate unsafe SRS (Structured Reference String) for testing
	// In production, use a proper trusted setup ceremony
	srs, srsLagrange, err := unsafekzg.NewSRS(ccs)
	if err != nil {
		return nil, fmt.Errorf("failed to generate SRS: %w", err)
	}

	// Run the PLONK setup to generate proving and verification keys
	// This is the "trusted setup" phase
	pk, vk, err := plonk.Setup(ccs, srs, srsLagrange)
	if err != nil {
		return nil, fmt.Errorf("failed to setup PLONK: %w", err)
	}

	return &SetupKeysWithCS{
		ProvingKey:       pk,
		VerifyingKey:     vk,
		ConstraintSystem: ccs,
	}, nil
}

// GenerateProof creates a zkSNARK proof that the prover knows the seed
// that, combined with participants hash, produces the given commitment
func GenerateProof(seed []byte, participantsContent []byte, keys *SetupKeysWithCS) (*ProofData, error) {
	// Validate seed length (must be 32 bytes = 256 bits)
	if len(seed) != 32 {
		return nil, fmt.Errorf("seed must be exactly 32 bytes, got %d", len(seed))
	}

	// Compute commitment and participants hash using the new formula
	commitment, participantsHash, err := ComputeCommitment(seed, participantsContent)
	if err != nil {
		return nil, fmt.Errorf("failed to compute commitment: %w", err)
	}

	// Create the witness (assignment of values to circuit variables)
	assignment := &SeedCommitmentCircuit{}

	// Assign the private input (seed)
	for i := 0; i < 32; i++ {
		assignment.Seed[i] = seed[i]
	}

	// Assign the public input (participants hash)
	for i := 0; i < 32; i++ {
		assignment.ParticipantsHash[i] = participantsHash[i]
	}

	// Assign the public input (commitment)
	for i := 0; i < 32; i++ {
		assignment.Commitment[i] = commitment[i]
	}

	// Create the witness from the assignment
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return nil, fmt.Errorf("failed to create witness: %w", err)
	}

	// Generate the proof using PLONK
	proof, err := plonk.Prove(keys.ConstraintSystem, keys.ProvingKey, witness)
	if err != nil {
		return nil, fmt.Errorf("failed to generate proof: %w", err)
	}

	return &ProofData{
		Proof:            proof,
		VerifyingKey:     keys.VerifyingKey,
		Commitment:       commitment,
		ParticipantsHash: participantsHash,
		Curve:            ecc.BN254,
	}, nil
}

// ComputeCommitment computes the commitment using the formula:
// commitment = SHA256(hex(seed) + hex(SHA256(participantsContent)))
func ComputeCommitment(seed []byte, participantsContent []byte) ([]byte, []byte, error) {
	if len(seed) != 32 {
		return nil, nil, fmt.Errorf("seed must be exactly 32 bytes, got %d", len(seed))
	}

	// Step 1: Compute SHA-256 hash of participants content
	participantsHash := sha256.Sum256(participantsContent)

	// Step 2: Convert seed to hex string
	seedHex := fmt.Sprintf("%x", seed)

	// Step 3: Convert participantsHash to hex string
	participantsHashHex := fmt.Sprintf("%x", participantsHash[:])

	// Step 4: Concatenate hex strings
	combined := seedHex + participantsHashHex

	// Step 5: Compute SHA-256 of concatenated hex string
	commitment := sha256.Sum256([]byte(combined))

	return commitment[:], participantsHash[:], nil
}
