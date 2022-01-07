package data

import (
	"encoding/json"

	_ "github.com/lib/pq"
)

func GetElements() (payload []byte, err error) {
	q := `select address, name, provider, n_protons, description, tx_code, up_votes, down_votes, num_fail, num_pass, transaction_hash from elements`

	rows, err := db.Query(q)
	if err != nil {
		return payload, err
	}
	defer rows.Close()

	var elements []Element
	for rows.Next() {
		var element Element
		var tx string
		rows.Scan(
			&element.Address,
			&element.Name,
			&element.Provider,
			&element.NProtons,
			&element.Description,
			&element.TxCode,
			&element.UpVotes,
			&element.DownVotes,
			&element.NumFail,
			&element.NumPass,
			&tx,
		)

		element.Transaction.TransactionHash = tx
		elements = append(elements, element)
	}

	payload, err = json.Marshal(APIElementDataResponse{Data: elements})
	return payload, err
}

func GetElement(element_id string) (payload []byte, err error) {
	q := `select address, name, provider, n_protons, description, tx_code, up_votes, down_votes, num_fail, num_pass, transaction_hash from elements`

	var element Element
	var tx string
	row := db.QueryRow(q, element_id)
	row.Scan(
		&element.Address,
		&element.Name,
		&element.Provider,
		&element.NProtons,
		&element.Description,
		&element.TxCode,
		&element.UpVotes,
		&element.DownVotes,
		&element.NumFail,
		&element.NumPass,
		&tx,
	)
	element.Transaction.TransactionHash = tx
	if err != nil {
		return payload, err
	}

	payload, err = json.Marshal(APIElementDetailResponse{Detail: element})
	return payload, err
}

func GetProtons(element_address string) (payload []byte, err error) {
	q := `select * from protons where element_address = $1`

	rows, err := db.Query(q, element_address)
	if err != nil {
		return payload, err
	}
	defer rows.Close()

	var protons []Proton
	for rows.Next() {
		var proton Proton
		rows.Scan(
			&proton.ProtonId,
			&proton.Name,
			&proton.Description,
			&proton.BaseUri,
			&proton.FkElement)
		protons = append(protons, proton)
	}

	payload, err = json.Marshal(APIProtonDataResponse{Data: protons})
	return payload, err
}

func GetProton(proton_id string) (payload []byte, err error) {
	q := `select * from protons where proton_id = $1`

	var proton Proton
	rows := db.QueryRow(q, proton_id)
	err = rows.Scan(
		&proton.ProtonId,
		&proton.Name,
		&proton.Description,
		&proton.BaseUri,
		&proton.FkElement)
	if err != nil {
		return payload, err
	}

	payload, err = json.Marshal(APIProtonDetailResponse{Detail: proton})
	return payload, err
}
