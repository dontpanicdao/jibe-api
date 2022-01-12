package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/jibe-api/internal/data"
)

type HTTPError struct {
	Error    string `json:"error"`
	Metadata string `json:"metadata"`
}

func writeGoodJSON(payload []byte, code int, w http.ResponseWriter) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
	return
}

func httpError(err error, metadata string, code int, w http.ResponseWriter) error {
	if code < 400 {
		return fmt.Errorf("%v not an http error code", code)
	}
	erry := HTTPError{
		Error:    err.Error(),
		Metadata: metadata,
	}
	returnErr, err := json.Marshal(erry)
	if err != nil {
		return err
	}
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(returnErr)
	return nil
}

func VerifyTx(jtx data.JSTransaction, r *http.Request) (hash string, is_valid bool) {
	pubKey := r.Header.Get("Public-Key")
	sigKey := r.Header.Get("Signing-Key")
	rSig := r.Header.Get("Signature-R")
	sSig := r.Header.Get("Signature-S")

	tx := jtx.ConvertTx()
	contentHash, err := data.StarkCurve.HashTx(caigo.HexToBN(pubKey), tx)
	if err != nil {
		return caigo.BigToHex(contentHash), false
	}

	pubX, pubY := data.StarkCurve.XToPubKey(sigKey)

	valid := data.StarkCurve.Verify(
		contentHash,
		caigo.StrToBig(rSig),
		caigo.StrToBig(sSig),
		pubX,
		pubY,
	)
	return caigo.BigToHex(contentHash), valid
}
