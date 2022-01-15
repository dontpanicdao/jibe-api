package data

import (
	"fmt"
	"math/big"
	"strings"
	"os/exec"

	"github.com/dontpanicdao/caigo"
)

func (cert Cert) FmtDefinitionEncoding(field string) (fmtEnc []*big.Int) {
	switch field {
	case "certUri":
		fmtEnc = append(fmtEnc, caigo.UTF8StrToBig(cert.CertUri))
	case "certKey":
		if cert.CertKey == "" {
			fmtEnc = append(fmtEnc, big.NewInt(0))
		} else {
			fmtEnc = append(fmtEnc, caigo.UTF8StrToBig(cert.CertKey))
		}
	case "certAttempt":
		fmtEnc = append(fmtEnc, caigo.UTF8StrToBig(cert.CertAttempt))
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
	q := `select address from elements where element_id = $1`

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


func (cert Cert) Grade(element_id string) (payload []byte, err error) {
	path2Grader := "/opt/jibe-api/grader.cairo"

	exec.Command("cairo-sharp", "submit", "--source", path2Grader, "--program_input")

	return payload, err	
}