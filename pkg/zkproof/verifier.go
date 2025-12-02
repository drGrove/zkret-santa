//go:build zksnark

package zkproof

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/frontend"
)

// VerifyProof verifies a zkSNARK proof against a commitment
// Returns true if the proof is valid, false otherwise
func VerifyProof(proofData *ProofData) (bool, error) {
	// Validate inputs
	if proofData == nil {
		return false, fmt.Errorf("proof data is nil")
	}
	if len(proofData.Commitment) != 32 {
		return false, fmt.Errorf("commitment must be 32 bytes, got %d", len(proofData.Commitment))
	}

	// Create the public witness (only the commitment, not the seed)
	publicAssignment := &SeedCommitmentCircuit{}

	// Assign the public input (commitment)
	for i := 0; i < 32; i++ {
		publicAssignment.Commitment[i] = proofData.Commitment[i]
	}

	// Create the public witness
	publicWitness, err := frontend.NewWitness(publicAssignment, ecc.BN254.ScalarField(), frontend.PublicOnly())
	if err != nil {
		return false, fmt.Errorf("failed to create public witness: %w", err)
	}

	// Verify the proof using PLONK
	err = plonk.Verify(proofData.Proof, proofData.VerifyingKey, publicWitness)
	if err != nil {
		return false, fmt.Errorf("proof verification failed: %w", err)
	}

	return true, nil
}

// VerifyProofWithCommitment verifies a proof and checks the commitment matches expected value
func VerifyProofWithCommitment(proofData *ProofData, expectedCommitment []byte) (bool, error) {
	// First verify the commitment matches
	if len(expectedCommitment) != 32 {
		return false, fmt.Errorf("expected commitment must be 32 bytes, got %d", len(expectedCommitment))
	}

	// Check if commitments match
	for i := 0; i < 32; i++ {
		if proofData.Commitment[i] != expectedCommitment[i] {
			return false, fmt.Errorf("commitment mismatch at byte %d", i)
		}
	}

	// Then verify the proof
	return VerifyProof(proofData)
}
