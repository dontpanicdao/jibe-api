package data

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/dontpanicdao/caigo"
)

func (prot Proton) FmtDefinitionEncoding(field string) (fmtEnc []*big.Int) {
	switch field {
	case "name":
		fmtEnc = append(fmtEnc, caigo.UTF8StrToBig(prot.Name))
	case "baseUri":
		fmtEnc = append(fmtEnc, caigo.UTF8StrToBig(prot.BaseUri))
	case "complete":
		if prot.Complete {
			fmtEnc = append(fmtEnc, big.NewInt(1))
		} else {
			fmtEnc = append(fmtEnc, big.NewInt(0))
		}
	}
	return fmtEnc
}

func (prot Proton) Verify(pubKey, sigKey, r, s, element_id string) (is_valid bool) {
	hash, err := TypedProton.GetMessageHash(caigo.HexToBN(pubKey), prot, StarkCurve)
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
