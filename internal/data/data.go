package data

import (
	"encoding/json"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type APISubjectDataResponse struct {
	Data []Subject `json:"data,omitempty"`
}

type APISubjectDetailResponse struct {
	Detail Subject `json:"detail,omitempty"`
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

type Subject struct {
	ID          int    `json:"id,omitempty"`
	Address     string `json:"assetAddress,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	NPhases     int    `json:"nPhases,omitempty"`
	Provider    string `json:"provier,omitempty"`
	CertAddress string `json:"certAddress,omitempty"`
	DAOAddress  string `json:"daoAddress,omitempty"`
	DAOScheme   string `json:"daoScheme,omitempty"`
}

type Phase struct {
	ID             int
	Name           string
	Description    string
	Nonce          int
	Authors        []string
	SubjectAddress string
	Provider       string
	CertURI        string
	DAOVoteID      int
}

type Cert struct {
	ID             int
	Name           string
	Address        string
	SubjectAddress string
	Provider       string
	BaseURI        string
}

func GetSubjects() (payload []byte, err error) {
	q := `select * from subjects`

	rows, err := db.Query(q)
	if err != nil {
		return payload, err
	}
	defer rows.Close()

	var subjects []Subject
	for rows.Next() {
		var subject Subject
		rows.Scan(
			&subject.ID,
			&subject.Address,
			&subject.Name,
			&subject.Description,
			&subject.NPhases,
			&subject.Provider,
			&subject.CertAddress,
			&subject.DAOAddress)
		subjects = append(subjects, subject)
	}

	payload, err = json.Marshal(APISubjectDataResponse{Data: subjects})
	return payload, err
}

func CreateSubject(s *TypedSubject) (err error) {
	q := `insert into subjects(
		address, name, n_phases, provider, dao_address
		) values($1, $2, $3, $4, $5)`
	_, err = db.Exec(q, s.Message.AssetAddress, s.Message.Name, s.Message.NPhases, s.Message.Provider, s.Message.DaoScheme)
	if err != nil {
		return err
	}

	return nil
}

func GetSubject(subject_id string) (payload []byte, err error) {
	q := `select * from subjects where address = $1`

	var subject Subject
	row := db.QueryRow(q, subject_id)
	err = row.Scan(
		&subject.ID,
		&subject.Address,
		&subject.Name,
		&subject.Description,
		&subject.NPhases,
		&subject.Provider,
		&subject.CertAddress,
		&subject.DAOAddress)
	if err != nil {
		return payload, err
	}

	payload, err = json.Marshal(APISubjectDetailResponse{Detail: subject})
	return payload, err
}

func GetPhases(subject_address string) (payload []byte, err error) {
	q := `select * from phases where subject_address = $1`

	rows, err := db.Query(q, subject_address)
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
		&phase.SubjectAddress,
		&phase.Provider,
		&phase.CertURI,
		&phase.DAOVoteID)
	if err != nil {
		return payload, err
	}

	payload, err = json.Marshal(APIPhaseDetailResponse{Detail: phase})
	return payload, err
}

func GetCert(subject_address string) (payload []byte, err error) {
	q := `select * from certs where subject_address = $1`

	var cert Cert
	rows := db.QueryRow(q, subject_address)
	err = rows.Scan(
		&cert.ID,
		&cert.Name,
		&cert.Address,
		&cert.SubjectAddress,
		&cert.Provider,
		&cert.BaseURI)
	if err != nil {
		return payload, err
	}

	payload, err = json.Marshal(APICertDetailResponse{Detail: cert})
	return payload, err
}
