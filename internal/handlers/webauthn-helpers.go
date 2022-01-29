package handlers

import (
	"fmt"
	"bytes"
	"crypto"
	"math/big"
	"crypto/x509"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/binary"

	"github.com/fxamacker/cbor/v2"
)

const COSEAlgES256 = -7

// https://github.com/fxamacker/webauthn/common.go
type AuthenticatorData struct {
	Raw          []byte                 // Complete raw authenticator data content.
	RPIDHash     []byte                 // SHA-256 hash of the RP ID the credential is scoped to.
	UserPresent  bool                   // User is present.
	UserVerified bool                   // User is verified.
	Counter      uint32                 // Signature Counter.
	AAGUID       []byte                 // AAGUID of the authenticator (optional).
	CredentialID []byte                 // Identifier of a public key credential source (optional).
	Credential   *Credential            // Algorithm and public key portion of a Relying Party-specific credential key pair (optional).
	Extensions   map[string]interface{} // Extension-defined authenticator data (optional).
}

// https://github.com/fxamacker/webauthn/credential.go
type Credential struct {
	Raw []byte
	SignatureAlgorithm
	ecdsa.PublicKey
}

// https://github.com/fxamacker/webauthn/signatureAlgorithm.go
type SignatureAlgorithm struct {
	Algorithm          x509.SignatureAlgorithm
	PublicKeyAlgorithm x509.PublicKeyAlgorithm
	Hash               crypto.Hash
	COSEAlgorithm      int
}

type rawCredential struct {
	Kty    int             `cbor:"1,keyasint"`
	Alg    int             `cbor:"3,keyasint"`
	CrvOrN cbor.RawMessage `cbor:"-1,keyasint"`
	XOrE   cbor.RawMessage `cbor:"-2,keyasint"`
	Y      cbor.RawMessage `cbor:"-3,keyasint"`
}

// https://github.com/fxamacker/webauthn/common.go
func parseAuthData(data []byte) (authData *AuthenticatorData, rest []byte, err error) {
	if len(data) < 37 {
		return authData, rest, fmt.Errorf("not enough bytes %v\n", len(data))
	}

	authData = &AuthenticatorData{Raw: data}
	copy(authData.RPIDHash, data)

	flags := data[32]
	authData.UserPresent = (flags & 0x01) > 0
	authData.UserVerified = (flags & 0x04) > 0
	credentialDataIncluded := (flags & 0x40) > 0

	if !authData.UserVerified {
		return authData, rest, fmt.Errorf("must have a verified user")
	}

	authData.Counter = binary.BigEndian.Uint32(data[33:37])
	rest = data[37:]

	if credentialDataIncluded {
		if len(rest) < 18 {
			return authData, rest, fmt.Errorf("not enough bytes in credential data %v\n", len(rest))
		}

		authData.AAGUID = make([]byte, 16)
		copy(authData.AAGUID, rest)

		idLength := binary.BigEndian.Uint16(rest[16:18])
		if len(rest[18:]) < int(idLength) {
			return authData, rest, fmt.Errorf("not enough bytes in id %v\n", len(rest))
		}

		authData.CredentialID = make([]byte, idLength)
		copy(authData.CredentialID, rest[18:])

		if authData.Credential, rest, err = parseCredential(rest[18+idLength:]); err != nil {
			return authData, rest, err
		}
	}
	return authData, rest, nil
}

// https://github.com/fxamacker/webauthn/credential.go
func parseCredential(data []byte) (credData *Credential, rest []byte, err error) {
	var raw rawCredential
	decoder := cbor.NewDecoder(bytes.NewReader(data))
	if err = decoder.Decode(&raw); err != nil {
		return credData, rest, err
	}
	if raw.Alg != COSEAlgES256 || raw.CrvOrN == nil || raw.XOrE == nil || raw.Y == nil {
		return credData, rest, fmt.Errorf("signature alg is not ES256")
	}

	var crvId int
	if err := cbor.Unmarshal(raw.CrvOrN, &crvId); err != nil {
		return credData, rest, fmt.Errorf("invalid ECDSA curve")
	}

	var xb []byte
	if err := cbor.Unmarshal(raw.XOrE, &xb); err != nil {
		return credData, rest, fmt.Errorf("invalid ECDSA x")
	}

	var yb []byte
	if err := cbor.Unmarshal(raw.Y, &yb); err != nil {
		return credData, rest, fmt.Errorf("invalid ECDSA y")
	}

	credData = &Credential{
		Raw: data,
		SignatureAlgorithm: SignatureAlgorithm{
			Algorithm: x509.ECDSAWithSHA256,
			PublicKeyAlgorithm: x509.ECDSA,
			Hash: crypto.SHA256,
			COSEAlgorithm: COSEAlgES256,
		},
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X: new(big.Int).SetBytes(xb),
			Y: new(big.Int).SetBytes(yb),
		},
	}
	return credData, data[decoder.NumBytesRead():], nil
}
