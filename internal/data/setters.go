package data

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dontpanicdao/caigo"
	"github.com/lib/pq"
)


func CreateElement(s *Element, hash string) (payload []byte, err error) {
	sn := StarkNetRequest{
		ContractAddress:    JIBE_ADDRESS,
		EntryPointSelector: caigo.BigToHex(caigo.GetSelectorFromName("get_count")),
		Calldata:           []string{},
		Signature:          []string{},
	}
	element_contract_id := 0
	resp, err := sn.Call(false)
	if err == nil {
		val, err := strconv.ParseInt(strings.Replace(resp[0], "0x", "", -1), 16, 32)
		if err == nil {
			element_contract_id = int(val)
		}
	}

	q := `insert into elements(
		num_pass, num_fail, element_contract_id, address, name, n_protons, provider, description, 
		molecule_address, reward_erc20_address, reward_amount_low, reward_amount_high, 
		cert_uri, rubric_uri, rubric_hash_low, rubric_hash_high,
		tx_code, transaction_hash, content_hash, dob, reward_symbol
		) values(0, 0, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`
	_, err = db.Exec(
		q,
		element_contract_id,
		strings.TrimLeft(s.Address, "0x"),
		s.Name,
		s.NProtons,
		s.Provider,
		s.Description,
		strings.TrimLeft(s.MoleculeAddress, "0x"),
		strings.TrimLeft(s.RewardErc20Address, "0x"),
		s.RewardAmountLow,
		s.RewardAmountHigh,
		s.CertUri,
		s.RubricUri,
		s.RubricHashLow,
		s.RubricHashHigh,
		s.TxCode,
		strings.TrimLeft(s.Transaction.TransactionHash, "0x"),
		strings.TrimLeft(hash, "0x"),
		time.Now().Unix(),
		s.RewardSymbol,
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
	q := `insert into element_cert_keys(cert_keys, fk_element)
	 values($1, $2)`

	_, err = db.Exec(
		q,
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
	q = `update elements set n_protons = n_protons + 1 where element_contract_id = $1`
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
	q := `select element_contract_id from elements where element_contract_id = $1`
	err = db.QueryRow(q, element_id).Scan(&id_check)
	if err != nil || id_check != element_id {
		return payload, fmt.Errorf("could not find corresponding exam key: %v %v", id_check, err)
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

func (fmtCred *FmtCredential) Create(pubKey string) (err error) {
	q := `insert into credentials(aaguid, credential_id, public_x, public_y, stark_key, counter) values($1, $2, $3, $4, $5, $6)`
	_, err = db.Exec(q, fmtCred.AAGUID, fmtCred.CredentialID, fmtCred.PublicKeyX, fmtCred.PublicKeyY, pubKey, fmtCred.Counter)
	return err
}

func SaveSession(challenge, user, ver, pub string) (err error) {
	q := `insert into webauthn_sessions(challenge, display_name, user_verification, public_key) values($1, $2, $3, $4)`
	_, err = db.Exec(q, challenge, user, ver, pub)
	return err
}

func GetSession(challenge string) (public_key string, err error) {
	q := `select public_key from webauthn_sessions where challenge = $1`
	err = db.QueryRow(q, challenge).Scan(&public_key)
	return public_key, err
}