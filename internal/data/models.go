package data

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	CREATED = "CREATED"
	UPDATED = "UPDATED"
	PASSED = "PASSED"
	FAILED = "FAILED"
	PENDING = "PENDING"
	RECEIVED = "RECEIVED"
	PROCESSED = "PROCESSED"
	SUBMITTED = "SUBMITTED"
	ATTESTED = "ATTESTED"
	CLAIMED = "CLAIMED"
	ACCEPTED_ON_L1 = "ACCEPTED_ON_L1"
	ACCEPTED_ON_L2 = "ACCEPTED_ON_L2"
	JIBE_ADDRESS = "0x0077b19d49e6069372d53e535fc9f3230a99b85ad46cc0934491bb6fb59a5a29"
	JIBE_ID = "alpha.jibe.buzz"
)

type APIElementDataResponse struct {
	Data []Element `json:"data,omitempty"`
}

type APIElementDetailResponse struct {
	Detail Element `json:"detail,omitempty"`
}

type APIProtonDataResponse struct {
	Data []Proton `json:"data,omitempty"`
}

type APIProtonDetailResponse struct {
	Detail Proton `json:"detail,omitempty"`
}

type APIResponse struct {
	Status  string `json:"status,omitempty"`
	TxCode  string `json:"txCode,omitempty"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type Element struct {
	ElementId          int           `json:"elementId"`
	ElementContractId  int           `json:"elementContractId"`
	Address            string        `json:"address"`
	Name               string        `json:"name"`
	Description        string        `json:"description"`
	NProtons           int           `json:"nProtons"`
	Provider           string        `json:"provider"`
	MoleculeAddress    string        `json:"moleculeAddress"`
	RewardErc20Address string        `json:"rewardErc20Address"`
	RewardAmountLow    string        `json:"rewardAmountLow"`
	RewardAmountHigh   string        `json:"rewardAmountHigh"`
	RewardSymbol       string        `json:"rewardSymbol"`
	UpVotes            int           `json:"upVotes,omitempty"`
	DownVotes          int           `json:"downVotes,omitempty"`
	NumFail            int           `json:"numFail,omitempty"`
	NumPass            int           `json:"numPass,omitempty"`
	CertUri            string        `json:"certUri,omitempty"`
	RubricUri          string        `json:"rubricUri,omitempty"`
	RubricHashLow      string        `json:"rubricHashLow"`
	RubricHashHigh     string        `json:"rubricHashHigh"`
	TxCode             string        `json:"txCode"`
	TransactionHash    string        `json:"transactionHash"`
	Transaction        JSTransaction `json:"transaction"`
}

type Proton struct {
	ProtonId    int    `json:"protonId"`
	Name        string `json:"name"`
	BaseUri     string `json:"baseUri"`
	Description string `json:"description,omitempty"`
	Complete    bool   `json:"complete,omitempty"`
	FkElement   int    `json:"fkElement"`
}

type User struct {
	UserId          int    `json:"userId"`
	Address         string `json:"address"`
	Username        string `json:"username"`
	Accumen         int    `json:"accumen"`
	Location        string `json:"location"`
	Description     string `json:"description"`
	PrimaryMolecule string `json:"primaryMolecule,omitempty"`
	PfpUri          string `json:"pfpUri,omitempty"`
	TwitterUri      string `json:"twitterUri,omitempty"`
	DiscordUri      string `json:"discordUri,omitempty"`
	GithubUri       string `json:"githubUri,omitempty"`
	IsStudent       string `json:"isStudent"`
	IsTeacer        string `json:"isTeacher"`
}

type ElementCertKeys struct {
	CertKeys  []string `json:"certKeys"`
	CertUri   string   `json:"certUri"`
	RubricUri string   `json:"rubricUri"`
	FkElement int      `json:"fkElement"`
}

type ElementAttempts struct {
	ElementName       string `json:"elementName"`
	Passed            bool   `json:"passed"`
	Score             int    `json:"score"`
	Fact              string `json:"fact"`
	FactLow           string `json:"factLow"`
	FactHigh          string `json:"factHigh"`
	FactJobId         string `json:"factJobId"`
	Status            string `json:"status"`
	L1Tx              string `json:"l1Tx"`
	PublicKey         string `json:"publicKey"`
	ElementContractId int    `json:"elementContractId"`
	FkUser            int    `json:"fkUser"`
}

type UserAttempts struct {
	User     User              `json:"user"`
	Attempts []ElementAttempts `json:"attempts"`
}

type ProtonCompletions struct {
	Passed      bool   `json:"passed"`
	Score       int    `json:"score"`
	ResponseUri string `json:"responseUri,omitempty"`
	fkProton    int    `json:"fkProton"`
	fkUser      int    `json:"fkUser"`
}

// struct to catch starknet.js transaction payloads
type JSTransaction struct {
	Calldata           []string `json:"calldata"`
	ContractAddress    string   `json:"contract_address"`
	EntryPointSelector string   `json:"entry_point_selector"`
	EntryPointType     string   `json:"entry_point_type"`
	JSSignature        []string `json:"signature"`
	TransactionHash    string   `json:"transaction_hash"`
	Type               string   `json:"type"`
	Nonce              string   `json:"nonce"`
}

type FmtCredential struct {
	AAGUID string `json:"aaguid"`
	CredentialID string `json:"credentialId"`
	PublicKeyX string `json:"publicKeyX"`
	PublicKeyY string `json:"publicKeyY"`
	Counter uint32 `json:"counter"`
	DisplayName string `json:"displayName"`
	StarkKey string `json:"starkKey"`
}

type CredentialAssertion struct {
	Response PublicKeyCredentialRequestOptions `json:"publicKey"`
}

type PublicKeyCredentialRequestOptions struct {
	Challenge          []byte                   `json:"challenge"`
	Timeout            int                         `json:"timeout,omitempty"`
	RelyingPartyID     string                      `json:"rpId,omitempty"`
	AllowedCredentials []CredentialDescriptor      `json:"allowCredentials"`
	UserVerification   string `json:"userVerification,omitempty"` // Default is "preferred"
	Extensions         string    `json:"extensions,omitempty"`
}

type CredentialDescriptor struct {
	Type string `json:"type"`
	CredentialID []byte `json:"id"`
	Transport []string `json:"transports,omitempty"`
}

type Attrs struct {
	Questions []struct {
		Question string `json:"question,omitempty"`
		Answers  []struct {
			Answer string `json:"answer,omitempty`
		} `json:"answers,omitempty"`
	} `json:"questions"`
}

func (a Attrs) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Attrs) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
