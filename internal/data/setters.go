package data

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

const CREATED = "CREATED"
const UPDATED = "UPDATED"
const PASSED = "PASSED"
const FAILED = "FAILED"
const PROCESSED = "PROCESSED"

func CreateElement(s *Element, hash string) (payload []byte, err error) {
	q := `insert into elements(
		num_pass, num_fail, address, name, n_protons, provider, description,
		tx_code, transaction_hash, content_hash, dob
		) values(0, 0, $1, $2, $3, $4, $5, $6, $7, $8, $9)`
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
	cr := APIResponse{
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
	cr := APIResponse{
		Status: CREATED,
		Error:  "",
	}

	payload, err = json.Marshal(cr)

	return payload, err
}

func (prot Proton) Create(element_id string) (payload []byte, err error) {
	q := `insert into protons(name, base_uri, fk_element)
	 values($1, $2, $3)`

	_, err = db.Exec(
		q,
		prot.Name,
		prot.BaseUri,
		element_id,
	)
	q = `update elements set n_protons = n_protons + 1 where element_id = $1`
	_, err = db.Exec(q, element_id)

	if err != nil {
		return payload, err
	}

	cr := APIResponse{
		Status: CREATED,
		Error:  "",
	}

	payload, err = json.Marshal(cr)

	return payload, err
}

func CreateQuestions(attrs Attrs, element_id string) (payload []byte, err error) {
	fmt.Println("ATTRS: ", attrs)
	var id_check string
	q := `select element_id from elements where element_id = $1`
	err = db.QueryRow(q, element_id).Scan(&id_check)
	if err != nil || id_check != element_id {
		return payload, fmt.Errorf("could not find corresponding exam key: %v %v\n", id_check, err)
	}

	q = `insert into custom_exams(answers, fk_element)
		values($1, $2)`

	_, err = db.Exec(
		q,
		attrs,
		element_id,
	)
	if err != nil {
		return payload, err
	}

	cr := APIResponse{
		Status: CREATED,
		Error:  "",
	}

	payload, err = json.Marshal(cr)

	return payload, err
}
