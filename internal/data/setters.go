package data

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/lib/pq"
)

const CREATED = "CREATED"

func CreateElement(s *Element, hash string) (payload []byte, err error) {
	q := `insert into elements(
		address, name, n_protons, provider, description,
		tx_code, transaction_hash, content_hash, dob
		) values($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err = db.Exec(
		q,
		strings.TrimLeft(s.Address, "0x"),
		s.Name,
		s.NProtons,
		s.Provider,
		s.Description,
		s.TxCode,
		strings.TrimLeft(s.Transaction.TransactionHash, "0x"),
		strings.TrimLeft(hash, "0x"),
		time.Now().Unix(),
	)
	if err != nil {
		return payload, err
	}
	cr := CreatedResponse{
		Status: CREATED,
		Error:  "",
		TxCode: s.TxCode,
	}
	payload, err = json.Marshal(cr)

	return payload, err
}

func (cert Cert) Create(element_id string) (payload []byte, err error) {
	q := `insert into element_cert_keys(cert_uri, cert_keys, fk_element)
	 values($1, $2, $3)`

	_, err = db.Exec(
		q,
		cert.CertUri,
		pq.Array(strings.Split(cert.CertKey, ",")),
		element_id,
	)
	if err != nil {
		return payload, err
	}
	cr := CreatedResponse{
		Status: CREATED,
		Error:  "",
	}

	payload, err = json.Marshal(cr)

	return payload, err
}
