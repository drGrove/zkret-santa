//go:build zksnark

package zkproof

import (
	"crypto/sha256"
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/frontend"
)

// VerifyProof verifies a zkSNARK proof against a commitment and participants hash
// Returns true if the proof is valid, false otherwise
func VerifyProof(proofData *ProofData) (bool, error) {
	// Validate inputs
	if proofData == nil {
		return false, fmt.Errorf("proof data is nil")
	}
	if len(proofData.Commitment) != 32 {
		return false, fmt.Errorf("commitment must be 32 bytes, got %d", len(proofData.Commitment))
	}
	if len(proofData.ParticipantsHash) != 32 {
		return false, fmt.Errorf("participants hash must be 32 bytes, got %d", len(proofData.ParticipantsHash))
	}

	// Create the public witness (participants hash and commitment, not the seed)
	publicAssignment := &SeedCommitmentCircuit{}

	// Assign the public input (participants hash)
	for i := 0; i < 32; i++ {
		publicAssignment.ParticipantsHash[i] = proofData.ParticipantsHash[i]
	}

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

// VerifyProofWithCommitment verifies a proof and checks the participants hash matches the provided participants content
func VerifyProofWithCommitment(proofData *ProofData, participantsContent []byte) (bool, error) {
	// Compute the expected participants hash from the provided content
	expectedParticipantsHash := sha256.Sum256(participantsContent)

	// Verify the participants hash matches
	if len(proofData.ParticipantsHash) != 32 {
		return false, fmt.Errorf("proof participants hash must be 32 bytes, got %d", len(proofData.ParticipantsHash))
	}

	// Check if participants hashes match
	for i := 0; i < 32; i++ {
		if proofData.ParticipantsHash[i] != expectedParticipantsHash[i] {
			return false, fmt.Errorf("participants hash mismatch at byte %d: proof has %x, expected %x",
				i, proofData.ParticipantsHash[i], expectedParticipantsHash[i])
		}
	}

	// Then verify the proof
	return VerifyProof(proofData)
}
