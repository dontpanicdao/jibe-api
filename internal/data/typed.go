package data

import (
	"strings"
	"math/big"

	"github.com/dontpanicdao/caigo"
)

type TypedSubject struct {
	Types struct {
		StarkNetDomain []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"StarkNetDomain"`
		Exam []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"Exam"`
	} `json:"types"`
	PrimaryType string `json:"primaryType"`
	Domain      struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		ChainID int    `json:"chainId"`
	} `json:"domain"`
	Message struct {
		Name         string `json:"name"`
		AssetAddress string `json:"assetAddress"`
		NPhases      string `json:"nPhases"`
		DaoScheme    string `json:"daoScheme"`
		SignerScheme string `json:"signerScheme"`
		Provider     string `json:"provider"`
	} `json:"message"`
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
		ContractAddress: jsToBN(jtx.ContractAddress),
		EntryPointSelector: jsToBN(jtx.EntryPointSelector),
		EntryPointType: jtx.EntryPointType,
		TransactionHash: jsToBN(jtx.TransactionHash),
		Type: jtx.Type,
		Nonce: jsToBN(jtx.Nonce),
	}
	for _, cd := range jtx.Calldata {
		tx.Calldata = append(tx.Calldata, jsToBN(cd))
	}
	for _, sigElem := range jtx.JSSignature {
		tx.Signature = append(tx.Signature, jsToBN(sigElem))
	}
	return tx
}

func jsToBN(str string) (*big.Int) {
	if strings.Contains(str, "0x") {
		return caigo.HexToBN(str)
	} else {
		return caigo.StrToBig(str)
	}
}