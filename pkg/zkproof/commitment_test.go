//go:build zksnark

package zkproof

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestComputeCommitment(t *testing.T) {
	// Test case 1: Valid seed and participants content
	seed := make([]byte, 32)
	for i := range seed {
		seed[i] = byte(i)
	}
	participantsContent := []byte("Alice\nBob\nCarol\nDave\n")

	commitment, participantsHash, err := ComputeCommitment(seed, participantsContent)
	if err != nil {
		t.Fatalf("ComputeCommitment failed: %v", err)
	}

	// Verify commitment is 32 bytes
	if len(commitment) != 32 {
		t.Errorf("Expected commitment to be 32 bytes, got %d", len(commitment))
	}

	// Verify participants hash is 32 bytes
	if len(participantsHash) != 32 {
		t.Errorf("Expected participants hash to be 32 bytes, got %d", len(participantsHash))
	}

	// Verify participants hash is correct
	expectedParticipantsHash := sha256.Sum256(participantsContent)
	if hex.EncodeToString(participantsHash) != hex.EncodeToString(expectedParticipantsHash[:]) {
		t.Errorf("Participants hash mismatch.\nExpected: %s\nGot: %s",
			hex.EncodeToString(expectedParticipantsHash[:]),
			hex.EncodeToString(participantsHash))
	}

	// Verify commitment is computed correctly
	// commitment = SHA256(hex(seed) + hex(participantsHash))
	seedHex := hex.EncodeToString(seed)
	participantsHashHex := hex.EncodeToString(participantsHash)
	combined := seedHex + participantsHashHex
	expectedCommitment := sha256.Sum256([]byte(combined))

	if hex.EncodeToString(commitment) != hex.EncodeToString(expectedCommitment[:]) {
		t.Errorf("Commitment mismatch.\nExpected: %s\nGot: %s",
			hex.EncodeToString(expectedCommitment[:]),
			hex.EncodeToString(commitment))
	}
}

func TestComputeCommitmentDeterministic(t *testing.T) {
	// Test that the same inputs produce the same outputs
	seed := make([]byte, 32)
	for i := range seed {
		seed[i] = byte(i * 2)
	}
	participantsContent := []byte("Alice\nBob\n")

	commitment1, participantsHash1, err := ComputeCommitment(seed, participantsContent)
	if err != nil {
		t.Fatalf("First ComputeCommitment failed: %v", err)
	}

	commitment2, participantsHash2, err := ComputeCommitment(seed, participantsContent)
	if err != nil {
		t.Fatalf("Second ComputeCommitment failed: %v", err)
	}

	if hex.EncodeToString(commitment1) != hex.EncodeToString(commitment2) {
		t.Errorf("Commitments are not deterministic")
	}

	if hex.EncodeToString(participantsHash1) != hex.EncodeToString(participantsHash2) {
		t.Errorf("Participants hashes are not deterministic")
	}
}

func TestComputeCommitmentInvalidSeed(t *testing.T) {
	// Test with invalid seed length
	testCases := []struct {
		name     string
		seedLen  int
		shouldError bool
	}{
		{"valid 32 bytes", 32, false},
		{"invalid 16 bytes", 16, true},
		{"invalid 64 bytes", 64, true},
		{"invalid 0 bytes", 0, true},
	}

	participantsContent := []byte("Alice\nBob\n")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			seed := make([]byte, tc.seedLen)
			_, _, err := ComputeCommitment(seed, participantsContent)

			if tc.shouldError && err == nil {
				t.Errorf("Expected error for seed length %d, but got none", tc.seedLen)
			}
			if !tc.shouldError && err != nil {
				t.Errorf("Expected no error for seed length %d, but got: %v", tc.seedLen, err)
			}
		})
	}
}

func TestComputeCommitmentDifferentParticipants(t *testing.T) {
	// Test that different participants produce different commitments
	seed := make([]byte, 32)
	for i := range seed {
		seed[i] = byte(i)
	}

	participants1 := []byte("Alice\nBob\nCarol\n")
	participants2 := []byte("Alice\nBob\nDave\n")

	commitment1, _, err := ComputeCommitment(seed, participants1)
	if err != nil {
		t.Fatalf("First ComputeCommitment failed: %v", err)
	}

	commitment2, _, err := ComputeCommitment(seed, participants2)
	if err != nil {
		t.Fatalf("Second ComputeCommitment failed: %v", err)
	}

	if hex.EncodeToString(commitment1) == hex.EncodeToString(commitment2) {
		t.Errorf("Different participants should produce different commitments")
	}
}

func TestComputeCommitmentDifferentSeeds(t *testing.T) {
	// Test that different seeds produce different commitments
	seed1 := make([]byte, 32)
	seed2 := make([]byte, 32)
	for i := range seed1 {
		seed1[i] = byte(i)
		seed2[i] = byte(i + 1)
	}

	participants := []byte("Alice\nBob\n")

	commitment1, _, err := ComputeCommitment(seed1, participants)
	if err != nil {
		t.Fatalf("First ComputeCommitment failed: %v", err)
	}

	commitment2, _, err := ComputeCommitment(seed2, participants)
	if err != nil {
		t.Fatalf("Second ComputeCommitment failed: %v", err)
	}

	if hex.EncodeToString(commitment1) == hex.EncodeToString(commitment2) {
		t.Errorf("Different seeds should produce different commitments")
	}
}
