package data

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/dontpanicdao/caigo"
)

type Cert struct {
	CertUri string
	CertKey string
}

func (cert Cert) FmtDefinitionEncoding(field string) (fmtEnc []*big.Int) {
	switch field {
	case "certUri":
		fmtEnc = append(fmtEnc, caigo.UTF8StrToBig(cert.CertUri))
	case "certKey":
		fmtEnc = append(fmtEnc, caigo.UTF8StrToBig(cert.CertKey))
	}
	return fmtEnc
}

func (cert Cert) Verify(pubKey, sigKey, r, s, element_id string) (is_valid bool) {
	keys := strings.Split(cert.CertKey, ",")
	if len(keys) == 0 {
		fmt.Println("length is bad: ", len(keys))
		return false
	}

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

// struct to catch starknet.js transaction payloads
type JSTransaction struct {
	Calldata           []string `json:"calldata"`
	ContractAddress    string   `json:"contract_address"`
	EntryPointSelector string   `json:"entry_point_selector"`
	EntryPointType     string   `json:"entry_point_type"`
	JSSignature        []string `json:"signature"`
	TransactionHash    string   `json:"transaction_hash"`
	Type               string   `json:"type"`
	Nonce              string   `json:"nonce"`
}

func (jtx JSTransaction) ConvertTx() (tx caigo.Transaction) {
	tx = caigo.Transaction{
		ContractAddress:    jsToBN(jtx.ContractAddress),
		EntryPointSelector: jsToBN(jtx.EntryPointSelector),
		EntryPointType:     jtx.EntryPointType,
		TransactionHash:    jsToBN(jtx.TransactionHash),
		Type:               jtx.Type,
		Nonce:              jsToBN(jtx.Nonce),
	}
	for _, cd := range jtx.Calldata {
		tx.Calldata = append(tx.Calldata, jsToBN(cd))
	}
	for _, sigElem := range jtx.JSSignature {
		tx.Signature = append(tx.Signature, jsToBN(sigElem))
	}
	return tx
}

func jsToBN(str string) *big.Int {
	if strings.Contains(str, "0x") {
		return caigo.HexToBN(str)
	} else {
		return caigo.StrToBig(str)
	}
}
