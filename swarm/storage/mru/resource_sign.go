package mru

import (
	"crypto/ecdsa"

	"github.com/TreasureChain/go-tstchain/common"
	"github.com/TreasureChain/go-tstchain/crypto"
)

const signatureLength = 65

// Signature is an alias for a static byte array with the size of a signature
type Signature [signatureLength]byte

// Signer signs Mutable Resource update payloads
type Signer interface {
	Sign(common.Hash) (Signature, error)
	Address() common.Address
}

// GenericSigner implements the Signer interface
// It is the vanilla signer that probably should be used in most cases
type GenericSigner struct {
	PrivKey *ecdsa.PrivateKey
	address common.Address
}

// NewGenericSigner builds a signer that will sign everything with the provided private key
func NewGenericSigner(privKey *ecdsa.PrivateKey) *GenericSigner {
	return &GenericSigner{
		PrivKey: privKey,
		address: crypto.PubkeyToAddress(privKey.PublicKey),
	}
}

// Sign signs the supplied data
// It wraps the tstchain crypto.Sign() method
func (s *GenericSigner) Sign(data common.Hash) (signature Signature, err error) {
	signaturebytes, err := crypto.Sign(data.Bytes(), s.PrivKey)
	if err != nil {
		return
	}
	copy(signature[:], signaturebytes)
	return
}

// PublicKey returns the public key of the signer's private key
func (s *GenericSigner) Address() common.Address {
	return s.address
}
