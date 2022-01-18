package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var ALPHA_BASE string = "https://alpha4.starknet.io/feeder_gateway"
var ALPHA_MAINNET_BASE string = "https://alpha-mainnet.starknet.io/feeder_gateway"

type StarkNetRequest struct {
	ContractAddress    string   `json:"contract_address"`
	EntryPointSelector string   `json:"entry_point_selector"`
	Calldata           []string `json:"calldata"`
	Signature          []string `json:"signature"`
}

type StarkResp struct {
	Result []string `json:"result"`
}

type TransactionStatus struct {
	TxStatus  string `json:"tx_status"`
	BlockHash string `json:"block_hash"`
}

type Transaction struct {
	TransactionIndex int `json:"transaction_index"`
	BlockNumber      int `json:"block_number"`
	Transaction      struct {
		Signature          []string `json:"signature"`
		EntryPointType     string   `json:"entry_point_type"`
		TransactionHash    string   `json:"transaction_hash"`
		Calldata           []string `json:"calldata"`
		EntryPointSelector string   `json:"entry_point_selector"`
		ContractAddress    string   `json:"contract_address"`
		Type               string   `json:"type"`
	} `json:"transaction"`
	BlockHash string `json:"block_hash"`
	Status    string `json:"status"`
}

func (sn StarkNetRequest) Call(prod bool) (resp []string, err error) {
	var url string
	if prod {
		url = fmt.Sprintf("%s/call_contract", ALPHA_MAINNET_BASE)
	} else {
		url = fmt.Sprintf("%s/call_contract", ALPHA_BASE)
	}
	method := "POST"

	pay, err := json.Marshal(sn)
	if err != nil {
		return resp, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(pay))
	if err != nil {
		return resp, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	defer res.Body.Close()

	var sRes StarkResp
	json.NewDecoder(res.Body).Decode(&sRes)
	return sRes.Result, nil
}

func GetTransactionStatus(txHash string, prod bool) (status TransactionStatus, err error) {
	var url string
	if prod {
		url = fmt.Sprintf("%s/get_transaction_status?transactionHash=%s", ALPHA_MAINNET_BASE, txHash)
	} else {
		url = fmt.Sprintf("%s/get_transaction_status?transactionHash=%s", ALPHA_BASE, txHash)
	}
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return status, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return status, err
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&status)
	return status, nil
}

func GetTransaction(txHash string, prod bool) (tx Transaction, err error) {
	var url string
	if prod {
		url = fmt.Sprintf("%s/get_transaction?transactionHash=%s", ALPHA_MAINNET_BASE, txHash)
	} else {
		url = fmt.Sprintf("%s/get_transaction?transactionHash=%s", ALPHA_BASE, txHash)
	}
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return tx, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return tx, err
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&tx)
	return tx, nil
}
