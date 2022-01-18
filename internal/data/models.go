package data

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"
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
	ElementId         int           `json:"elementId"`
	ElementContractId int           `json:"elementContractId"`
	Address           string        `json:"address"`
	Name              string        `json:"name"`
	Description       string        `json:"description"`
	NProtons          int           `json:"nProtons"`
	Provider          string        `json:"provider"`
	UpVotes           int           `json:"upVotes,omitempty"`
	DownVotes         int           `json:"downVotes,omitempty"`
	NumFail           int           `json:"numFail,omitempty"`
	NumPass           int           `json:"numPass,omitempty"`
	CertUri           string        `json:"certUri,omitempty"`
	RubricUri         string        `json:"rubricUri,omitempty"`
	TxCode            string        `json:"txCode"`
	TransactionHash   string        `json:"transactionHash"`
	Transaction       JSTransaction `json:"transaction"`
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
	UserId      int    `json:"userId"`
	Address     string `json:"address"`
	Username    string `json:"username"`
	Accumen    int `json:"accumen"`
	Location 		string `json:"location"`
	Description string `json:"description"`
	PrimaryMolecule string `json:"primaryMolecule,omitempty"`
	PfpUri      string `json:"pfpUri,omitempty"`
	TwitterUri  string `json:"twitterUri,omitempty"`
	DiscordUri  string `json:"discordUri,omitempty"`
	GithubUri   string `json:"githubUri,omitempty"`
	IsStudent   string `json:"isStudent"`
	IsTeacer    string `json:"isTeacher"`
}

type ElementCertKeys struct {
	CertKeys  []string `json:"certKeys"`
	CertUri   string   `json:"certUri"`
	RubricUri string   `json:"rubricUri"`
	FkElement int      `json:"fkElement"`
}

type ElementAttempts struct {
	ElementName string `json:"elementName"`
	Passed    bool   `json:"passed"`
	Score     int    `json:"score"`
	Fact      string `json:"fact"`
	FactJobId string `json:"factJobId"`
	Status string `json:"status"`
	PublicKey string `json:"publicKey"`
	ElementId int    `json:"elementId"`
	FkUser    int    `json:"fkUser"`
}

type UserAttempts struct {
	User User `json:"user"`
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
