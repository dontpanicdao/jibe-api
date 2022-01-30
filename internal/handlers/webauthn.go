package handlers 

import (
	"fmt"
	"bytes"
	"reflect"
	"strings"
	"net/http"
	"math/big"
	"crypto/rand"
	"crypto/sha256"
	"crypto/elliptic"
	"crypto/ecdsa"
	"encoding/asn1"
	"encoding/json"
	"encoding/base64"

	// "github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/fxamacker/cbor/v2"
	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/jibe-api/internal/data"
)

// ChallengeLength - Length of bytes to generate for a challenge
const (
	ChallengeLength = 32
	JIBE_ID = "alpha.jibe.buzz"
)

type URLEncodedBase64 []byte

// Challenge that should be signed and returned by the authenticator
type Challenge URLEncodedBase64

type CredentialType string

type AuthenticationExtensions map[string]interface{}

type AuthenticatorTransport string

type AuthenticatorAttachment string

type UserVerificationRequirement string

type ConveyancePreference string

// https://github.com/duo-labs/webauthn/blob/master/protocol/options.go
type CredentialCreation struct {
	Response PublicKeyCredentialCreationOptions `json:"publicKey"`
}

// https://github.com/duo-labs/webauthn/blob/master/protocol/options.go
type CredentialAssertion struct {
	Response PublicKeyCredentialRequestOptions `json:"publicKey"`
}

// https://github.com/duo-labs/webauthn/blob/master/protocol/entities.go
type CredentialEntity struct {
	Name string `json:"name"`
	Icon string `json:"icon,omitempty"`
}

// https://github.com/duo-labs/webauthn/blob/master/protocol/assertion.go
type CredentialAssertionResponse struct {
	Id string `json:"id"`
	RawId URLEncodedBase64 `json:"rawId"`
	Type string `json:"type"`
	AssertionResponse AuthenticatorAssertionResponse `json:"response"`
}

// https://github.com/duo-labs/webauthn/blob/master/protocol/entities.go
type RelyingPartyEntity struct {
	CredentialEntity
	ID string `json:"id"`
}

// https://github.com/duo-labs/webauthn/blob/master/protocol/options.go
type CredentialParameter struct {
	Type      string                       `json:"type"` //"public-key"
	Algorithm  int `json:"alg"` //webauthncose.COSEAlgorithmIdentifier(-7)
}

// https://github.com/duo-labs/webauthn/blob/master/protocol/assertion.go
type AuthenticatorAssertionResponse struct {
	ClientDataJSON URLEncodedBase64 `json:"clientDataJSON"`
	AuthenticatorData URLEncodedBase64 `json:"authenticatorData"`
	Signature         URLEncodedBase64 `json:"signature"`
	UserHandle        URLEncodedBase64 `json:"userHandle,omitempty"`
}

// https://github.com/duo-labs/webauthn/blob/master/protocol/entities.go
type UserEntity struct {
	CredentialEntity
	DisplayName string `json:"displayName,omitempty"`
	ID []byte `json:"id"`
}

type CredentialDescriptor struct {
	Type CredentialType `json:"type"`
	CredentialID []byte `json:"id"`
	Transport []AuthenticatorTransport `json:"transports,omitempty"`
}

// https://github.com/duo-labs/webauthn/blob/master/protocol/entities.go
type WebAuthnCredential struct {
	PublicKey string `json:"publicKey"`
	DisplayName string `json:"displayName"`
	AttestationType string `json:"attestationType"`
	AuthType string `json:"authType"`
	UserVerification string `json:"userVerification"`
	ResidentKeyRequirement string `json:"residentKeyRequirement"`
	TxAuthExtension string `json:"txAuthExtension"`
}

type AuthenticatorSelection struct {
	AuthenticatorAttachment AuthenticatorAttachment `json:"authenticatorAttachment,omitempty"`
	RequireResidentKey *bool `json:"requireResidentKey,omitempty"`
	UserVerification UserVerificationRequirement `json:"userVerification,omitempty"`
}

// https://github.com/duo-labs/webauthn/blob/master/protocol/options.go
type PublicKeyCredentialCreationOptions struct {
	Challenge              Challenge                `json:"challenge"`
	RelyingParty           RelyingPartyEntity       `json:"rp"`
	User                   UserEntity               `json:"user"`
	Parameters             []CredentialParameter    `json:"pubKeyCredParams,omitempty"`
	AuthenticatorSelection AuthenticatorSelection   `json:"authenticatorSelection,omitempty"`
	Timeout                int                      `json:"timeout,omitempty"`
	CredentialExcludeList  []CredentialDescriptor   `json:"excludeCredentials,omitempty"`
	Extensions             AuthenticationExtensions `json:"extensions,omitempty"`
	Attestation            ConveyancePreference     `json:"attestation,omitempty"`
}

// https://github.com/duo-labs/webauthn/blob/master/protocol/options.go
type PublicKeyCredentialRequestOptions struct {
	Challenge          Challenge                   `json:"challenge"`
	Timeout            int                         `json:"timeout,omitempty"`
	RelyingPartyID     string                      `json:"rpId,omitempty"`
	AllowedCredentials []CredentialDescriptor      `json:"allowCredentials,omitempty"`
	UserVerification   UserVerificationRequirement `json:"userVerification,omitempty"` // Default is "preferred"
	Extensions         AuthenticationExtensions    `json:"extensions,omitempty"`
}

// https://github.com/duo-labs/webauthn/blob/master/protocol/session.go
type SessionData struct {
	Challenge            string                               `json:"challenge"`
	UserID               []byte                               `json:"user_id"`
	AllowedCredentialIDs [][]byte                             `json:"allowed_credentials,omitempty"`
	UserVerification     UserVerificationRequirement `json:"userVerification"`
	Extensions           AuthenticationExtensions    `json:"extensions,omitempty"`
}

// https://github.com/duo-labs/webauthn/blob/master/protocol/client.go
type CollectedClientData struct {
	Type         string  `json:"type"` // webauthn.create || webauthn.get
	Challenge    string        `json:"challenge"`
	Origin       string        `json:"origin"`
	TokenBinding *TokenBinding `json:"tokenBinding,omitempty"`
	Hint string `json:"new_keys_may_be_added_here,omitempty"`
}

type TokenBinding struct {
	Status string `json:"status"` // "present" || "not-present"
	ID     string             `json:"id,omitempty"`
}

type AttestationObject struct {
	AuthnData []byte          `cbor:"authData"`
	Fmt       string          `cbor:"fmt"`
	AttStmt   cbor.RawMessage `cbor:"attStmt"`
}

type SessionResponse struct {
	Id string `json:"id"`
	RawId string `json:"rawId"`
	Type string `json:"type"`
	Response struct {
		AttestationObject string `json:"attestationObject"`
		ClientDataJSON string `json:"clientDataJSON"`
	} `json:"response"`
}

type ECDSASignature struct {
	R, S *big.Int
}

func GetAllCredentials(w http.ResponseWriter, r *http.Request) {
	challenge, err := CreateChallenge()
	if err != nil {
		httpError(err, "could not create challenge", http.StatusInternalServerError, w)
		return
	}

	pay, err := data.GetCredentials(challenge[:])
	if err != nil {
		httpError(err, "could not fetch credentials", http.StatusInternalServerError, w)
		return
	}

	writeGoodJSON(pay, http.StatusOK, w)
	return
}

func BeginLogin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	challenge, err := CreateChallenge()
	if err != nil {
		httpError(err, "could not create challenge", http.StatusInternalServerError, w)
		return
	}

	cred, err := data.GetCredential(vars["public_key"])
	if err != nil {
		pay, err := json.Marshal(data.APIResponse{Status: "no-user", Message: "no-user"})
		if err != nil {
			httpError(err, "could not marshal response", http.StatusInternalServerError, w)
			return
		}
	
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(pay)
		return
	}

	err = data.SaveSession(base64.RawURLEncoding.EncodeToString(challenge), "", "required", vars["public_key"])
	if err != nil {
		httpError(err, "could not create webauthn session", http.StatusBadRequest, w)
		return
	}
	
	credId, err := base64.URLEncoding.DecodeString(cred.CredentialID)
	if err != nil {
		httpError(err, "could not decode credential", http.StatusInternalServerError, w)
		return
	}

	reqOptions := PublicKeyCredentialRequestOptions{
		Challenge: challenge,
		Timeout: 15000,
		RelyingPartyID: JIBE_ID,
		AllowedCredentials: []CredentialDescriptor{CredentialDescriptor{Type: "public-key", CredentialID: credId}},
		UserVerification: "required",
	}

	pay, err := json.Marshal(CredentialAssertion{Response: reqOptions})
	if err != nil {
		httpError(err, "could not create webauthn session", http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(pay)
	return
}

func FinishLogin(w http.ResponseWriter, r *http.Request) {
	var car CredentialAssertionResponse
	err := json.NewDecoder(r.Body).Decode(&car)
	if err != nil {
		httpError(err, "could not parse credential assertion", http.StatusInternalServerError, w)
		return
	}
	_, err = base64.RawURLEncoding.DecodeString(car.Id)
	if err != nil || car.Id == "" || car.Type != "public-key" {
		httpError(err, "assertion incorrect format", http.StatusBadRequest, w)
		return
	}

	var cCollected CollectedClientData
	err = json.Unmarshal(car.AssertionResponse.ClientDataJSON, &cCollected)
	if err != nil || car.Id == "" || car.Type != "public-key" {
		httpError(err, "assertion incorrect format", http.StatusBadRequest, w)
		return
	}

	fmt.Println("USER: ", string(car.AssertionResponse.UserHandle))

	cred, err := data.GetCredential(string(car.AssertionResponse.UserHandle))
	if err != nil {
		httpError(err, "could not get user", http.StatusInternalServerError, w)
		return
	}
	fmt.Println("CRED: ", cred.DisplayName, cred.StarkKey)
	
	storedId, err := base64.URLEncoding.DecodeString(cred.CredentialID)
	if !bytes.Equal(car.AssertionResponse.UserHandle, []byte(cred.StarkKey)) || !bytes.Equal(storedId, car.RawId) || err != nil {
		httpError(err, "incorrect user handler", http.StatusBadRequest, w)
		return
	}
	
	if !strings.Contains(cCollected.Origin, JIBE_ID) {
		httpError(err, "incorrect originr", http.StatusInternalServerError, w)
		return	
	}

	clientDataHash := sha256.Sum256(car.AssertionResponse.ClientDataJSON)
	sigData := append(car.AssertionResponse.AuthenticatorData, clientDataHash[:]...)

	pub := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X: caigo.HexToBN(cred.PublicKeyX),
		Y: caigo.HexToBN(cred.PublicKeyY),
	}

	e := &ECDSASignature{}
	_, err = asn1.Unmarshal(car.AssertionResponse.Signature, e)
	hash := sha256.Sum256(sigData)


	if !ecdsa.Verify(pub, hash[:], e.R, e.S) {
		httpError(fmt.Errorf("invalid signature"), "invalid signature", http.StatusBadRequest, w)
		return
	}

	pay, err := json.Marshal(data.APIResponse{Status: "login-success", Message: cred.DisplayName})
	if err != nil {
		httpError(err, "could not marshal response", http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(pay)
	return
}

func BeginRegistration(w http.ResponseWriter, r *http.Request) {
	wac := &WebAuthnCredential{}
	err := json.NewDecoder(r.Body).Decode(&wac)
	if err != nil {
		httpError(err, "could not parse json", http.StatusBadRequest, w)
		return
	}

	challenge, err := CreateChallenge()
	if err != nil {
		httpError(err, "could not create challenge", http.StatusInternalServerError, w)
		return
	}

	// cookie and local storage and webauthn and 
	wau := UserEntity{
		ID: []byte(wac.PublicKey),
		DisplayName: wac.DisplayName,
		CredentialEntity: CredentialEntity{
			Name: wac.DisplayName,
			Icon: "",
		},
	}

	relyingParty := RelyingPartyEntity{
		ID: JIBE_ID,
		CredentialEntity: CredentialEntity{
			Name: "Jibe(alpha)",
			Icon: "https://alpha.jibe.buzz/img/brand/jibeHex.png",
		},
	}

	reqResKey := false
	creationOptions := PublicKeyCredentialCreationOptions{
		Challenge: challenge,
		RelyingParty: relyingParty,
		User: wau,
		Parameters: []CredentialParameter{CredentialParameter{Type: "public-key", Algorithm: -7}},
		AuthenticatorSelection: AuthenticatorSelection{
			AuthenticatorAttachment: "platform",
			RequireResidentKey: &reqResKey,
			UserVerification: "required",
		},
		Timeout: 15000,
		Attestation: "direct",
	}

	resp := CredentialCreation{Response: creationOptions}
	sess := SessionData{
		Challenge: base64.RawURLEncoding.EncodeToString(challenge),
		UserID: wau.ID,
		UserVerification: creationOptions.AuthenticatorSelection.UserVerification,
	}

	err = data.SaveSession(sess.Challenge, wau.DisplayName, string(sess.UserVerification), string(wau.ID))
	if err != nil {
		httpError(err, "could not create webauthn session", http.StatusInternalServerError, w)
		return
	}

	pay, err := json.Marshal(resp)
	if err != nil {
		httpError(err, "could not create webauthn session", http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(pay)
	return
}

func FinishRegistration(w http.ResponseWriter, r *http.Request) {
	sess := &SessionResponse{}
	if err := json.NewDecoder(r.Body).Decode(&sess); err != nil {
		httpError(err, "could not parse json", http.StatusBadRequest, w)
		return
	}

	_, err := base64.RawURLEncoding.DecodeString(sess.Id)
	if err != nil {
		httpError(err, "could not decode session data", http.StatusBadRequest, w)
		return			
	}

	cJson, err := base64.RawURLEncoding.DecodeString(sess.Response.ClientDataJSON)
	if err != nil {
		httpError(err, "could not decode JSON data", http.StatusBadRequest, w)
		return		
	}

	cCollected := &CollectedClientData{}
	if err := json.NewDecoder(bytes.NewReader(cJson)).Decode(cCollected); err != nil {
		httpError(err, "could not decode JSON data", http.StatusBadRequest, w)
		return
	}

	atObj, err := base64.RawURLEncoding.DecodeString(sess.Response.AttestationObject)
	if err != nil {
		httpError(err, "could not decode attestation", http.StatusBadRequest, w)
		return
	}

	atCollected := &AttestationObject{}
	if err := cbor.Unmarshal(atObj, &atCollected); err != nil {
		httpError(err, "could not parse cbor", http.StatusBadRequest, w)
		return
	}

	authParsed, _, err := parseAuthData(atCollected.AuthnData)
	if err != nil {
		httpError(err, "could not parse auth data", http.StatusBadRequest, w)
		return
	}

	pubKey, displayName, err := data.GetSession(cCollected.Challenge)
	if err != nil {
		httpError(err, "no session for this challenge", http.StatusBadRequest, w)
		return
	}

	if !strings.Contains(cCollected.Origin, JIBE_ID) {
		httpError(fmt.Errorf("%v %v", cCollected.Origin, JIBE_ID), "incorrect origin", http.StatusBadRequest, w)
		return
	}

	rpIDHash := sha256.Sum256([]byte(JIBE_ID))
	if !bytes.Equal(authParsed.RPIDHash, rpIDHash[:]) {
		httpError(fmt.Errorf("%v", authParsed.RPIDHash), "incorrect rpId hash", http.StatusBadRequest, w)
		return
	}

	fmtCred := &data.FmtCredential{
		AAGUID: base64.URLEncoding.EncodeToString(authParsed.AAGUID),
		CredentialID: base64.URLEncoding.EncodeToString(authParsed.CredentialID),
		PublicKeyX: caigo.BigToHex(authParsed.Credential.PublicKey.X),
		PublicKeyY: caigo.BigToHex(authParsed.Credential.PublicKey.Y),
		Counter: authParsed.Counter,
		DisplayName: displayName,
	}

	err = fmtCred.Create(pubKey)
	if err != nil {
		httpError(err, "could not save credential", http.StatusBadRequest, w)
		return
	}

	pay, err := json.Marshal(fmtCred)
	if err != nil {
		httpError(err, "could not marshal credential", http.StatusBadRequest, w)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(pay)
	return
}

// Create a new challenge to be sent to the authenticator. The spec recommends using
// at least 16 bytes with 100 bits of entropy. We use 32 bytes.
func CreateChallenge() (Challenge, error) {
	challenge := make([]byte, ChallengeLength)
	_, err := rand.Read(challenge)
	if err != nil {
		return nil, err
	}
	return challenge, nil
}

func (c Challenge) String() string {
	return base64.RawURLEncoding.EncodeToString(c)
}

// UnmarshalJSON base64 decodes a URL-encoded value, storing the result in the
// provided byte slice.
func (dest *URLEncodedBase64) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		return nil
	}

	// Trim the leading spaces
	data = bytes.Trim(data, "\"")
	out := make([]byte, base64.RawURLEncoding.DecodedLen(len(data)))
	n, err := base64.RawURLEncoding.Decode(out, data)
	if err != nil {
		return err
	}

	v := reflect.ValueOf(dest).Elem()
	v.SetBytes(out[:n])
	return nil
}

// MarshalJSON base64 encodes a non URL-encoded value, storing the result in the
// provided byte slice.
func (data URLEncodedBase64) MarshalJSON() ([]byte, error) {
	if data == nil {
		return []byte("null"), nil
	}
	return []byte(`"` + base64.RawURLEncoding.EncodeToString(data) + `"`), nil
}