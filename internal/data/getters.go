package data

import (
	"fmt"
	"strings"
	"strconv"
	"os/exec"
	"database/sql"
	"encoding/json"

	"github.com/dontpanicdao/caigo"
	_ "github.com/lib/pq"
)

func GetElements() (payload []byte, err error) {
	q := `select element_id, element_contract_id, address, name, reward_amount_low, reward_amount_high,
	n_protons, tx_code, transaction_hash, up_votes, down_votes, num_fail, num_pass from elements`

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
			&element.ElementContractId,
			&element.Address,
			&element.Name,
			&element.RewardAmountLow,
			&element.RewardAmountHigh,
			&element.NProtons,
			&element.TxCode,
			&element.TransactionHash,
			&element.UpVotes,
			&element.DownVotes,
			&element.NumFail,
			&element.NumPass,
		)

		if element.TxCode != ACCEPTED_ON_L1 {
			stat, err := GetTransactionStatus(fmt.Sprintf("0x%s", element.TransactionHash), false)
			if err != nil {
				fmt.Println("starknet tx err: ", err)
			}
			fmt.Println("STAT: ", stat)
			if stat.TxStatus != element.TxCode {
				q = `update elements set tx_code = $1 where element_contract_id = $2`
				_, err = db.Exec(
					q,
					stat.TxStatus,
					element.ElementContractId,
				)
				if err != nil {
					fmt.Println("DB ERR: ", err)
				} else {
					element.TxCode = stat.TxStatus
				}
			}
		}

		elements = append(elements, element)
	}

	payload, err = json.Marshal(APIElementDataResponse{Data: elements})
	return payload, err
}

func GetElement(element_id string) (payload []byte, err error) {
	q := `select element_id, element_contract_id, address, name, provider, molecule_address,
	reward_erc20_address, reward_amount_low, reward_amount_high, n_protons, cert_uri, rubric_uri, description,
	tx_code, transaction_hash, up_votes, down_votes, num_fail, num_pass, reward_symbol from elements where element_contract_id = $1`

	var element Element
	row := db.QueryRow(q, element_id)
	row.Scan(
		&element.ElementId,
		&element.ElementContractId,
		&element.Address,
		&element.Name,
		&element.Provider,
		&element.MoleculeAddress,
		&element.RewardErc20Address,
		&element.RewardAmountLow,
		&element.RewardAmountHigh,
		&element.NProtons,
		&element.CertUri,
		&element.RubricUri,
		&element.Description,
		&element.TxCode,
		&element.TransactionHash,
		&element.UpVotes,
		&element.DownVotes,
		&element.NumFail,
		&element.NumPass,
		&element.RewardSymbol,
	)

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

func GetCredential(pubKey string) (cred FmtCredential, err error) {
	q := `select credential_id, public_x, public_y from credentials where stark_key = $1`
	err = db.QueryRow(q, pubKey).Scan(&cred.CredentialID, &cred.PublicKeyX, &cred.PublicKeyY)
	return cred, err
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
	if err != nil && err != sql.ErrNoRows {
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
	q := `select passed, elements.element_contract_id, elements.name, fact, fact_low, fact_high,
	fact_job_id, status from element_attempts join elements on fk_element = element_contract_id where public_key = $1`
	rows, err := db.Query(q, pubKey)
	if err != nil && err != sql.ErrNoRows {
		return attempts, err
	}
	defer rows.Close()

	for rows.Next() {
		var attempt ElementAttempts
		rows.Scan(
			&attempt.Passed,
			&attempt.ElementContractId,
			&attempt.ElementName,
			&attempt.Fact,
			&attempt.FactLow,
			&attempt.FactHigh,
			&attempt.FactJobId,
			&attempt.Status)

		attempts = append(attempts, attempt)
	}

	return attempts, err
}

func CheckFact(fact string) (payload []byte, err error) {
	factLow, factHigh := caigo.SplitFactStr(fact)
	fmt.Println("FACT: ", factLow, factHigh)

	low := caigo.HexToBN(factLow)
	high := caigo.HexToBN(factHigh)
	sn := StarkNetRequest{
		ContractAddress:    JIBE_ADDRESS,
		EntryPointSelector: caigo.BigToHex(caigo.GetSelectorFromName("get_fact_check")),
		Calldata:           []string{low.Text(10), high.Text(10)},
		Signature:          []string{},
	}
	fmt.Println("SELECTOR: ", caigo.BigToHex(caigo.GetSelectorFromName("get_fact_check")))

	snResp, err := sn.Call(false)
	fmt.Println("fact check: ", snResp)
	if err != nil || len(snResp) < 5 || snResp[2] == "0x0" || snResp[2] == "0" {
		return payload, fmt.Errorf("could not get element details from starknet: %v %v", err, len(snResp))
	}
	
	status := SUBMITTED

	if snResp[4] == "0x1" || snResp[4] == "1" {
		if snResp[1] == "0x1" || snResp[1] == "1" {
			status = CLAIMED
		} else {
			status = ATTESTED
		}
	} else {
		elemId, err := strconv.Atoi(strings.TrimLeft(snResp[0], "0x"))
		if err == nil {
			q := `select status, l1_tx from element_attempts where public_key = $1 and fk_element = $2`
			
			var dbStat, l1Tx string
			err = db.QueryRow(q, strings.TrimLeft(snResp[2], "0x"), elemId).Scan(&dbStat, &l1Tx)
			fmt.Println("THIS THAT THE OTHER: ", dbStat, l1Tx, err, snResp)
			if err == nil && l1Tx != "" {
				status = PENDING
			}
		}
	}

	q := `update element_attempts set status = $1 where fact = $2`
	_, err = db.Exec(
		q,
		status,
		fact,
	)
	if err != nil {
		return payload, err
	}

	cr := APIResponse{
		Status: status,
		Error:  "",
	}

	payload, err = json.Marshal(cr)

	return payload, err
}

func ShipSharpStark(fact, tx_hash string) (payload []byte, err error) {
	q := `update element_attempts set status = $1, l1_tx = $2 where fact = $3`
	_, err = db.Exec(
		q,
		PENDING,
		tx_hash,
		fact,
	)
	if err != nil {
		return payload, err
	}

	cr := APIResponse{
		Status: PENDING,
		Error:  "",
	}

	payload, err = json.Marshal(cr)

	return payload, err
}

func CheckFactJob(factJobId string) (payload []byte, err error) {
	q := `select status, fact_low, fact_high from element_attempts where fact_job_id = $1`
	var elemAttempt ElementAttempts
	err = db.QueryRow(q, factJobId).Scan(&elemAttempt.Status, &elemAttempt.FactLow, &elemAttempt.FactHigh)
	if err != nil {
		return payload, err
	}
	fmt.Println("CHECK FACT: ", elemAttempt)
	if elemAttempt.Status == SUBMITTED {
		cmd := exec.Command("/usr/local/bin/cairo-sharp", "status", factJobId)
		stdout, _ := cmd.Output()

		if strings.Contains(string(stdout), PROCESSED) {
			elemAttempt.Status = PROCESSED
			q := `update element_attempts set status = $1 where fact_job_id = $2`
			_, err = db.Exec(q, PROCESSED, factJobId)
			if err != nil {
				return payload, fmt.Errorf("unable to update processed fact: %v", err)
			}
		} else if strings.Contains(string(stdout), SUBMITTED) {
			elemAttempt.Status = SUBMITTED
			q := `update element_attempts set status = $1 where fact_job_id = $2`
			_, err = db.Exec(q, SUBMITTED, factJobId)
			if err != nil {
				return payload, fmt.Errorf("unable to update processed fact: : %v", err)
			}
		}
	}

	payload, err = json.Marshal(elemAttempt)
	return payload, err
}
