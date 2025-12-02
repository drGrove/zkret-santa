//go:build zksnark

package zkproof

import (
	"bytes"
	"encoding/hex"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/plonk"
	plonk_bn254 "github.com/consensys/gnark/backend/plonk/bn254"
)

// ProofData contains the zkSNARK proof and verification key
type ProofData struct {
	// Proof is the PLONK proof
	Proof plonk.Proof
	// VerifyingKey is the verification key
	VerifyingKey plonk.VerifyingKey
	// Commitment is the public SHA-256 hash of the seed
	Commitment []byte
	// Curve is the elliptic curve used (BN254)
	Curve ecc.ID
}

// SerializableProof is a JSON-serializable representation of the proof
type SerializableProof struct {
	Proof        string `json:"proof"`
	VerifyingKey string `json:"verifying_key"`
	Commitment   string `json:"commitment"`
	Curve        string `json:"curve"`
}

// ToSerializable converts ProofData to SerializableProof for JSON export
func (pd *ProofData) ToSerializable() (*SerializableProof, error) {
	// Serialize proof
	var proofBuf bytes.Buffer
	_, err := pd.Proof.WriteTo(&proofBuf)
	if err != nil {
		return nil, err
	}

	// Serialize verification key
	var vkBuf bytes.Buffer
	_, err = pd.VerifyingKey.WriteTo(&vkBuf)
	if err != nil {
		return nil, err
	}

	return &SerializableProof{
		Proof:        hex.EncodeToString(proofBuf.Bytes()),
		VerifyingKey: hex.EncodeToString(vkBuf.Bytes()),
		Commitment:   hex.EncodeToString(pd.Commitment),
		Curve:        pd.Curve.String(),
	}, nil
}

// FromSerializable converts SerializableProof back to ProofData
func FromSerializable(sp *SerializableProof) (*ProofData, error) {
	// Decode proof
	proofBytes, err := hex.DecodeString(sp.Proof)
	if err != nil {
		return nil, err
	}

	proof := &plonk_bn254.Proof{}
	_, err = proof.ReadFrom(bytes.NewReader(proofBytes))
	if err != nil {
		return nil, err
	}

	// Decode verification key
	vkBytes, err := hex.DecodeString(sp.VerifyingKey)
	if err != nil {
		return nil, err
	}

	vk := &plonk_bn254.VerifyingKey{}
	_, err = vk.ReadFrom(bytes.NewReader(vkBytes))
	if err != nil {
		return nil, err
	}

	// Decode commitment
	commitment, err := hex.DecodeString(sp.Commitment)
	if err != nil {
		return nil, err
	}

	return &ProofData{
		Proof:        proof,
		VerifyingKey: vk,
		Commitment:   commitment,
		Curve:        ecc.BN254,
	}, nil
}
