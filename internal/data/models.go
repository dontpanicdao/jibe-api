package data

import (
    "fmt"
    "database/sql/driver"
    "encoding/json"

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

type CreatedResponse struct {
	Status string `json:"status,omitempty"`
	TxCode string `json:"txCode"`
	Error  string `json:"error"`
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
	TxCode            string        `json:"txCode"`
	TransactionHash   string        `json:"transactionHash"`
	Transaction       JSTransaction `json:"transaction"`
}

type Proton struct {
	ProtonId    int    `json:"protonId"`
	Name        string `json:"name"`
	BaseUri     string `json:"baseUri"`
	Description string `json:"description,omitempty"`
	Complete     bool `json:"complete,omitempty"`
	FkElement   int    `json:"fkElement"`
}

type User struct {
	UserId      int    `json:"userId"`
	Address     string `json:"address"`
	Username    string `json:"username"`
	PfpUri      string `json:"pfpUri,omitempty"`
	Description string `json:"description,omitempty"`
	TwitterUri  string `json:"twitterUri,omitempty"`
	GithubUri   string `json:"githubUri,omitempty"`
	IsStudent   string `json:"isStudent"`
	IsTeacer    string `json:"isTeacher"`
}

type Fact struct {
	FactId     int    `json:"factId"`
	Fact       string `json:"fact"`
	FactHash   string `json:"factHash"`
	FactR      string `json:"factR"`
	FactS      string `json:"factS"`
	FactOutput string `json:"factOutput"`
	FactStatus string `json:"factStatus"`
}

type ElementCertKeys struct {
	CertKeys  []string `json:"certKeys"`
	CertUri   string   `json:"certUri"`
	FkElement int      `json:"fkElement"`
}

type ElementAttempts struct {
	Passed    bool `json:"passed"`
	Score     int  `json:"score"`
	FactId    int  `json:"factId"`
	ElementId int  `json:"elementId"`
	FkUser    int  `json:"fkUser"`
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

type Cert struct {
	CertUri string
	CertKey string
	CertAttempt string
}

type Attrs struct {
	Questions []struct{
		Question string `json:"question,omitempty"`
		Answers []struct {
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