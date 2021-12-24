package data

import (
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type APIDataResponse struct {
	Data []interface{} `json:"data,omitempty"`
}

type APIDetailResponse struct {
	Detail []interface{} `json:"detail,omitempty"`
}

type Subject struct {
	ID                  int
	Address     string
	Name                string
	NPhases             int
	Provider            string
	CertAddress string
	DAOAddress  string
}

type Phase struct {
	ID int
	Name string
	Nonce int
	Authors []string
	SubjectAddress string
	Provider string
	CertURI string
	DAOVoteID int
}

type Cert struct {
	ID int
	Name string
	Address string
	SubjectAddress string
	Provider string
	BaseURI string
}

func GetSubjects() (payload []byte, err error) {
	q := `select * from subjects`

	rows, err := db.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	var subjects []Subject
	for rows.Nex() {
		var subject Subject
		rows.Scan(
			&subject.ID,
			&subject.Address,
			&subject.Name,
			&subject.NPhases,
			&subject.Provider,
			&subject.CertAddress,
			&subject.DAOAddress)
		subjects = append(subjects, subject)
	}

	payload, err = json.Marshal(APIDataResponse{Data: subjects})
	return payload, err
}

func GetSubject(subject_id string) (payload []byte, err error) {
	q := `select * from subjects where address = $1`

	var subject Subject
	row := db.QueryRow(q, subject_id)
	err := row.Scan(
		&subject.ID,
		&subject.Address,
		&subject.Name,
		&subject.NPhases,
		&subject.Provider,
		&subject.CertAddress,
		&subject.DAOAddress)
	if err != nil {
		return err
	}

	payload, err = json.Marshal(APIDetailResponse{Detail: subject})
	return payload, err
}

func GetPhases(subject_address string) (payload []byte, err error) {
	q := `select * from phases where subject_address = $1`

	rows, err := db.Query(q, subject_address)
	if err != nil {
		return err
	}
	defer rows.Close()

	var phases []Subject
	for rows.Nex() {
		var phase Subject
		rows.Scan(
			&phase.ID,
			&phase.Name,
			&phase.Nonce,
			&phase.Authors,
			&phase.Provider,
			&phase.CertURI,
			&phase.DAOVoteID)
		phases = append(phases, phase)
	}

	payload, err = json.Marshal(APIDataResponse{Data: phases})
	return payload, err
}

func GetPhase(phase_vote_id int) (payload []byte, err error) {
	q := `select * from phases where dao_vote_id = $1`

	var phase Phase
	rows = db.Query(q, phase_vote_id)
	err := rows.Scan(
			&phase.ID,
			&phase.Name,
			&phase.Nonce,
			&phase.Authors,
			&phase.SubjectAddress,
			&phase.Provider,
			&phase.CertURI,
			&phase.DAOVoteID)

	payload, err = json.Marshal(APIDetailResponse{Detail: phase})
	return payload, err
}

func GetCert(subject_address string) (payload []byte, err error) {
	q := `select * from certs where subject_address = $1`

	var cert Cert
	rows = db.Query(q, subject_address)
	err := rows.Scan(
			&cert.ID,
			&cert.Name,
			&cert.Address,
			&cert.SubjectAddress,
			&cert.Provider,
			&cert.BaseURI)

	payload, err = json.Marshal(APIDetailResponse{Detail: cert})
	return payload, err
}