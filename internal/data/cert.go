package data

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/dontpanicdao/caigo"
	"github.com/lib/pq"
)

var (
	UNSUBMITTED = "UNSUBMITTED"
)

type Cert struct {
	CertUri     string
	CertKey     string
	CertAttempt string
	RubricUri   string
}

type ExamAttempt struct {
	Exam    *big.Int `json:"exam"`
	Key     *big.Int `json:"key"`
	Address *big.Int `json:"address"`
}

func (cert Cert) FmtDefinitionEncoding(field string) (fmtEnc []*big.Int) {
	switch field {
	case "certUri":
		if cert.CertUri == "" {
			fmtEnc = append(fmtEnc, big.NewInt(0))
		} else {
			fmtEnc = append(fmtEnc, caigo.UTF8StrToBig(cert.CertUri))
		}
	case "certKey":
		if cert.CertKey == "" {
			fmtEnc = append(fmtEnc, big.NewInt(0))
		} else {
			fmtEnc = append(fmtEnc, caigo.UTF8StrToBig(cert.CertKey))
		}
	case "certAttempt":
		if cert.CertAttempt == "" {
			fmtEnc = append(fmtEnc, big.NewInt(0))
		} else {
			fmtEnc = append(fmtEnc, caigo.UTF8StrToBig(cert.CertAttempt))
		}
	case "rubricUri":
		if cert.RubricUri == "" {
			fmtEnc = append(fmtEnc, big.NewInt(0))
		} else {
			fmtEnc = append(fmtEnc, caigo.UTF8StrToBig(cert.RubricUri))
		}
	}
	return fmtEnc
}

func (cert Cert) Verify(pubKey, sigKey, r, s, element_id string) (is_valid bool) {
	keys := strings.Split(cert.CertKey, ",")
	if len(keys) == 0 {
		fmt.Println("length is bad: ", len(keys))
		return false
	}

	fmt.Println("PUBKEY CERT THANG: ", pubKey, cert)
	hash, err := TypedCert.GetMessageHash(caigo.HexToBN(pubKey), cert, StarkCurve)
	if err != nil {
		fmt.Println("hash err: ", hash, err)
		return false
	}
	q := `select address from elements where element_contract_id = $1`

	var addr string
	row := db.QueryRow(q, element_id)
	row.Scan(&addr)

	if addr != strings.TrimLeft(pubKey, "0x") {
		fmt.Println("not the owner: ", addr, strings.TrimLeft(pubKey, "0x"))
		return false
	}

	x := caigo.HexToBN(sigKey)
	y := StarkCurve.GetYCoordinate(x)

	is_valid = StarkCurve.Verify(hash, caigo.StrToBig(r), caigo.StrToBig(s), x, y)
	return is_valid
}

func (cert Cert) VerifyAttempt(pubKey, sigKey, r, s, element_id string) (is_valid bool) {
	keys := strings.Split(cert.CertKey, ",")
	if len(keys) == 0 {
		fmt.Println("length is bad: ", len(keys))
		return false
	}

	fmt.Println("PUBKEY CERT THANG: ", pubKey, cert)
	hash, err := TypedCert.GetMessageHash(caigo.HexToBN(pubKey), cert, StarkCurve)
	if err != nil {
		fmt.Println("hash err: ", hash, err)
		return false
	}

	x := caigo.HexToBN(sigKey)
	y := StarkCurve.GetYCoordinate(x)

	is_valid = StarkCurve.Verify(hash, caigo.StrToBig(r), caigo.StrToBig(s), x, y)
	return is_valid
}

func (cert Cert) AddRubric(element_id string) (payload []byte, err error) {
	q := `update element_cert_keys set rubric_uri = $1 where fk_element = $2`

	_, err = db.Exec(q, cert.RubricUri, element_id)
	if err != nil {
		return payload, err
	}

	cr := APIResponse{
		Status:  UPDATED,
		Message: cert.RubricUri,
		Error:   "",
	}

	payload, err = json.Marshal(cr)

	return payload, err
}

func (cert Cert) Grade(element_id, pubKey string) (payload []byte, err error) {
	// resp is an array of hex strings
	sn := StarkNetRequest{
		ContractAddress:    JIBE_ADDRESS,
		EntryPointSelector: caigo.BigToHex(caigo.GetSelectorFromName("get_element")),
		Calldata:           []string{element_id},
		Signature:          []string{},
	}

	snResp, err := sn.Call(false)
	if err != nil || len(snResp) < 11 {
		return payload, fmt.Errorf("could not get element details from starknet: %v %v", err, len(snResp))
	}

	rubricUri := caigo.HexToShortStr(snResp[8])

	fmt.Println("RUBRIC URI FROM CHAIN: ", rubricUri)

	q := `select cert_keys from element_cert_keys where fk_element = $1`

	var certKeys []string
	err = db.QueryRow(q, element_id).Scan(pq.Array(&certKeys))
	if err != nil {
		return payload, err
	}
	if len(certKeys) < 1 {
		return payload, fmt.Errorf("could not get cert keys for element: %v", element_id)
	}

	resp, err := http.Get(rubricUri)
	if err != nil {
		return payload, err
	}

	rubricPath := fmt.Sprintf("/opt/jibe-api/rubrics/%v.cairo", element_id)
	grader, err := os.Create(rubricPath)
	if err != nil {
		return payload, err
	}

	_, err = io.Copy(grader, resp.Body)
	if err != nil {
		return payload, err
	}

	resp.Body.Close()
	grader.Close()

	attempt := ExamAttempt{
		Exam:    caigo.UTF8StrToBig(strings.ReplaceAll(cert.CertAttempt, " ", "")),
		Key:     caigo.UTF8StrToBig(strings.ReplaceAll(strings.Join(certKeys, ","), " ", "")),
		Address: caigo.HexToBN(pubKey),
	}

	ja, err := json.Marshal(attempt)
	if err != nil {
		return payload, err
	}
	attemptPath := fmt.Sprintf("/opt/jibe-api/attempts/%v_%v.json", element_id, pubKey)
	err = os.WriteFile(attemptPath, ja, 0644)
	if err != nil {
		return payload, err
	}

	cmd := exec.Command("/usr/local/bin/cairo-sharp", "submit", "--source", rubricPath, "--program_input", attemptPath)
	stdout, err := cmd.Output()
	if err != nil {
		return payload, err
	}

	factJobId := factJobReg.FindString(string(stdout))
	fact := factReg.FindString(string(stdout))

	if fact == "" {
		q = `update elements set num_fail = num_fail + 1 where element_contract_id = $1`
		_, err = db.Exec(q, element_id)
		if err != nil {
			return payload, fmt.Errorf("invalid attempt: %v %v", cert.CertAttempt, err)
		}
		q = `insert into element_attempts(passed, public_key, fk_element) values(false, $1, $2)`
		_, err = db.Exec(q, strings.TrimLeft(pubKey, "0x"), element_id)
		if err != nil {
			return payload, fmt.Errorf("invalid attempt: %v %v", cert.CertAttempt, err)
		}

		cr := APIResponse{
			Status:  FAILED,
			Message: cert.CertUri,
			Error:   "",
		}
		payload, err = json.Marshal(cr)
		return payload, err
	} else {
		q = `update elements set num_pass = num_pass + 1 where element_contract_id = $1`
		_, err = db.Exec(q, element_id)
		if err != nil {
			return payload, fmt.Errorf("could not write good fact to db: %v %v", fact, err)
		}

		factLow, factHigh := caigo.SplitFactStr(fact)
	
		q = `insert into element_attempts(passed, status, public_key, fact, fact_job_id, fact_low, fact_high, fk_element) values(true, $1, $2, $3, $4, $5, $6, $7)`
		_, err = db.Exec(q, UNSUBMITTED, strings.TrimLeft(pubKey, "0x"), strings.TrimLeft(fact, "0x"), factJobId, strings.TrimLeft(factLow, "0x"), strings.TrimLeft(factHigh, "0x"), element_id)
		if err != nil {
			return payload, fmt.Errorf("could not write good fact to db: %v %v", fact, err)
		}
	
		cr := APIResponse{
			Status:  PASSED,
			Message: fact,
			Error:   "",
		}
		payload, err = json.Marshal(cr)
	
		return payload, err
	}
}
