package data

import (
	"fmt"
	"strings"
	"os/exec"
	"encoding/json"

	_ "github.com/lib/pq"
)

func GetElements() (payload []byte, err error) {
	q := `select element_id, address, name, provider, n_protons, description, 
	tx_code, up_votes, down_votes, num_fail, num_pass, transaction_hash
	from elements`

	rows, err := db.Query(q)
	if err != nil {
		return payload, err
	}
	defer rows.Close()

	var elements []Element
	for rows.Next() {
		var element Element
		rows.Scan(
			&element.ElementId,
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
			&element.TransactionHash,
		)

		elements = append(elements, element)
	}

	payload, err = json.Marshal(APIElementDataResponse{Data: elements})
	return payload, err
}

func GetElement(element_id string) (payload []byte, err error) {
	q := `select element_id, address, name, provider, n_protons, description, 
	tx_code, up_votes, down_votes, num_fail, num_pass from elements where element_id = $1`

	var element Element
	row := db.QueryRow(q, element_id)
	row.Scan(
		&element.ElementId,
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
	)

	q = `select transaction_hash from elements where element_id = $1`

	var tx string
	_ = db.QueryRow(q, element_id).Scan(&tx)

	if element.TxCode != "ACCEPTED_ON_L1" {
		stat, _ := GetTransactionStatus(fmt.Sprintf("0x%s", tx), false)
		fmt.Println("STAT: ", stat)
		if stat.TxStatus != element.TxCode {
			q = `update elements set tx_code = $1 where element_id = $2`
			_, err = db.Exec(
				q,
				stat.TxStatus,
				element_id,
			)
			if err != nil {
				fmt.Println("DB ERR: ", err)
			} else {
				element.TxCode = stat.TxStatus
			}
		}
	}

	q = `select cert_uri, rubric_uri from element_cert_keys where fk_element = $1`
	var certUri, rubricUri string
	_ = db.QueryRow(q, element_id).Scan(&certUri, &rubricUri)

	element.TransactionHash = tx
	element.CertUri = certUri
	element.RubricUri = rubricUri

	payload, err = json.Marshal(APIElementDetailResponse{Detail: element})
	return payload, err
}

func GetProtons(element_id string) (payload []byte, err error) {
	q := `select proton_id, name, base_uri, description from protons where fk_element = $1`

	rows, err := db.Query(q, element_id)
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
			&proton.BaseUri,
			&proton.Description)
		protons = append(protons, proton)
	}

	payload, err = json.Marshal(APIProtonDataResponse{Data: protons})
	return payload, err
}

func GetProton(proton_id string) (payload []byte, err error) {
	q := `select proton_id, name, base_uri, description from protons where proton_id = $1`

	var proton Proton
	rows := db.QueryRow(q, proton_id)
	err = rows.Scan(
		&proton.ProtonId,
		&proton.Name,
		&proton.BaseUri,
		&proton.Description)
	if err != nil {
		return payload, err
	}

	payload, err = json.Marshal(APIProtonDetailResponse{Detail: proton})
	return payload, err
}

func GetCustomCert(element_id string) (payload []byte, err error) {
	q := `select answers from custom_exams where fk_element = $1`

	attrs := new(Attrs)
	err = db.QueryRow(q, element_id).Scan(&attrs)
	if err != nil {
		return payload, err
	}
	payload, err = json.Marshal(attrs)
	return payload, err
}

func GetUser(pubKey string) (payload []byte, err error) {
	q := `select username, accumen, location, description, twitter_uri, discord_uri, github_uri from users where address = $1`

	var user User
	row := db.QueryRow(q, pubKey)
	err = row.Scan(
		&user.Username,
		&user.Accumen,
		&user.Location,
		&user.Description,
		&user.TwitterUri,
		&user.DiscordUri,
		&user.GithubUri,
	)
	if err != nil {
		return payload, err
	}
	attempts, err := GetUserAttempts(pubKey)
	if err != nil {
		return payload, err
	}
	payload, err = json.Marshal(UserAttempts{User: user, Attempts: attempts})
	return payload, err
}

func GetUserAttempts(pubKey string) (attempts []ElementAttempts, err error) {
	q := `select passed, elements.element_id, fact, fact_job_id, elements.name, status
		from element_attempts
		join elements on elements.element_id = element_attempts.fk_element
		where public_key = $1`
	rows, err := db.Query(q, pubKey)
	if err != nil {
		return attempts, err
	}
	defer rows.Close()

	for rows.Next() {
		var attempt ElementAttempts
		rows.Scan(
			&attempt.ElementName,
			&attempt.Passed,
			&attempt.ElementId,
			&attempt.Fact,
			&attempt.FactJobId,
			&attempt.Status)

		if attempt.Passed && attempt.Status == "SUBMITTED" {
			cmd := exec.Command("/usr/local/bin/cairo-sharp", "status", attempt.FactJobId)
			stdout, err := cmd.Output()
			if err != nil {
				fmt.Println("unable to get submitted fact status: ", string(stdout))
			} else {
				if strings.Contains(string(stdout), "PROCESSED") {
					attempt.Status = PROCESSED
					q = `update element_attempts set status = 'PROCESSED' where element_id = $1 and public_key = $2`
					_, err = db.Exec(q, attempt.ElementId, pubKey)
					if err != nil {
						fmt.Println("unable to update processed fact: ", string(stdout))
					}
				}
			}
		}
		
		attempts = append(attempts, attempt)
	}

	return attempts, err
}