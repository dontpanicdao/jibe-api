package data

import (
	"encoding/json"

	_ "github.com/lib/pq"
)
const CREATED = "CREATED"

type APIElementDataResponse struct {
	Data []Element `json:"data,omitempty"`
}

type APIElementDetailResponse struct {
	Detail Element `json:"detail,omitempty"`
}

type APIPhaseDataResponse struct {
	Data []Phase `json:"data,omitempty"`
}

type APIPhaseDetailResponse struct {
	Detail Phase `json:"detail,omitempty"`
}

type APICertDataResponse struct {
	Data []Cert `json:"data,omitempty"`
}

type APICertDetailResponse struct {
	Detail Cert `json:"detail,omitempty"`
}

type CreatedResponse struct {
	Status string `json:"status,omitempty"`
	TxCode string `json:"txCode"`
	Error string `json:"error"`
}

type Element struct {
	ID          int    `json:"id,omitempty"`
	Address     string `json:"address,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	NPhases     int    `json:"nPhases,omitempty"`
	Provider    string `json:"provider,omitempty"`
	ElementID   int `json:"elementID,omitempty"`
	TxCode   string `json:"txCode,omitempty"`
	Transaction JSTransaction `json:"transaction"`
}

type Phase struct {
	ID             int
	Name           string
	Description    string
	Nonce          int
	Authors        []string
	ElementAddress string
	Provider       string
	CertURI        string
	DAOVoteID      int
}

type Cert struct {
	ID             int
	Name           string
	Address        string
	ElementAddress string
	Provider       string
	BaseURI        string
}

func GetElements() (payload []byte, err error) {
	q := `select * from elements`

	rows, err := db.Query(q)
	if err != nil {
		return payload, err
	}
	defer rows.Close()

	var elements []Element
	for rows.Next() {
		var element Element
		rows.Scan(
			&element.ID,
			&element.Address,
			&element.Name,
			&element.NPhases,
			&element.Provider,
			&element.Description,
			&element.TxCode,
			&element.Transaction.TransactionHash,
		)
		elements = append(elements, element)
	}

	payload, err = json.Marshal(APIElementDataResponse{Data: elements})
	return payload, err
}

func CreateElement(s *Element) (payload []byte, err error) {
	q := `insert into elements(
		address, name, n_phases, provider, dao_address, tx_code, tx_hash, description 
		) values($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = db.Exec(
		q, 
		s.Address, 
		s.Name, 
		s.NPhases, 
		s.Provider, 
		s.Description,
		s.TxCode,
		s.Transaction.TransactionHash,
	)
	if err != nil {
		return payload, err
	}
	cr := CreatedResponse{
		Status: CREATED,
		Error: "",
		TxCode: s.TxCode,
	}
	payload, err = json.Marshal(cr)

	return payload, err
}

func GetElement(element_id string) (payload []byte, err error) {
	q := `select * from elements where address = $1`

	var element Element
	row := db.QueryRow(q, element_id)
	err = row.Scan(
		&element.ID,
		&element.Address,
		&element.Name,
		&element.Description,
		&element.NPhases,
		&element.Provider,
	)
	if err != nil {
		return payload, err
	}

	payload, err = json.Marshal(APIElementDetailResponse{Detail: element})
	return payload, err
}

func GetPhases(element_address string) (payload []byte, err error) {
	q := `select * from phases where element_address = $1`

	rows, err := db.Query(q, element_address)
	if err != nil {
		return payload, err
	}
	defer rows.Close()

	var phases []Phase
	for rows.Next() {
		var phase Phase
		rows.Scan(
			&phase.ID,
			&phase.Name,
			&phase.Description,
			&phase.Nonce,
			&phase.Authors,
			&phase.Provider,
			&phase.CertURI,
			&phase.DAOVoteID)
		phases = append(phases, phase)
	}

	payload, err = json.Marshal(APIPhaseDataResponse{Data: phases})
	return payload, err
}

func GetPhase(phase_vote_id string) (payload []byte, err error) {
	q := `select * from phases where dao_vote_id = $1`

	var phase Phase
	rows := db.QueryRow(q, phase_vote_id)
	err = rows.Scan(
		&phase.ID,
		&phase.Name,
		&phase.Description,
		&phase.Nonce,
		&phase.Authors,
		&phase.ElementAddress,
		&phase.Provider,
		&phase.CertURI,
		&phase.DAOVoteID)
	if err != nil {
		return payload, err
	}

	payload, err = json.Marshal(APIPhaseDetailResponse{Detail: phase})
	return payload, err
}

func GetCert(element_address string) (payload []byte, err error) {
	q := `select * from certs where element_address = $1`

	var cert Cert
	rows := db.QueryRow(q, element_address)
	err = rows.Scan(
		&cert.ID,
		&cert.Name,
		&cert.Address,
		&cert.ElementAddress,
		&cert.Provider,
		&cert.BaseURI)
	if err != nil {
		return payload, err
	}

	payload, err = json.Marshal(APICertDetailResponse{Detail: cert})
	return payload, err
}
