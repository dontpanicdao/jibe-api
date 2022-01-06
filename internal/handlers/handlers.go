package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/jibe-api/internal/data"
)

func SubjectsFetch(w http.ResponseWriter, r *http.Request) {
	subjects, err := data.GetSubjects()
	if err != nil {
		httpError(err, "subjects db pull", http.StatusInternalServerError, w)
		return
	}

	writeGoodJSON(subjects, http.StatusOK, w)
	return
}

func SubjectFetch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	subject, err := data.GetSubject(vars["subject_id"])
	if err != nil {
		httpError(err, "subject db pull", http.StatusInternalServerError, w)
		return
	}

	writeGoodJSON(subject, http.StatusOK, w)
	return
}

func PhasesFetch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	phases, err := data.GetPhases(vars["subject_id"])
	if err != nil {
		httpError(err, "phases db pull", http.StatusInternalServerError, w)
		return
	}

	writeGoodJSON(phases, http.StatusOK, w)
	return
}

func PhaseFetch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	phase, err := data.GetPhase(vars["phase_id"])
	if err != nil {
		httpError(err, "phase db pull", http.StatusInternalServerError, w)
		return
	}

	writeGoodJSON(phase, http.StatusOK, w)
	return
}

func CertFetch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	phase, err := data.GetPhase(vars["subject_id"])
	if err != nil {
		httpError(err, "phase db pull", http.StatusInternalServerError, w)
		return
	}

	writeGoodJSON(phase, http.StatusOK, w)
	return
}

func CertKey(w http.ResponseWriter, r *http.Request) {
	// TODO: handle the posting of an exam key
}

func CreateSubject(w http.ResponseWriter, r *http.Request) {
	subject := &data.Subject{}
	err := json.NewDecoder(r.Body).Decode(&subject)
	if err != nil {
		httpError(err, "could not parse json", http.StatusBadRequest, w)
		return
	}

	pubKey := r.Header.Get("Public-Key")
	sigKey := r.Header.Get("Signing-Key")
	rSig := r.Header.Get("Signature-R")
	sSig := r.Header.Get("Signature-S")

	fmt.Println("PUBKEY: ", pubKey)
	tx := subject.Transaction.ConvertTx()
	fmt.Println("CallData: ", tx.Calldata)
	fmt.Println("Contract Addr: ", tx.ContractAddress)
	fmt.Println("Nonce: ", tx.Nonce)
	contentHash, err := data.StarkCurve.HashTx(caigo.HexToBN(pubKey), tx)
	if err != nil {
		httpError(err, "could not hash transaction", http.StatusBadRequest, w)
		return
	}
	fmt.Println("CONTENT HASH: ", contentHash)
	fmt.Println("CONTENT HASH HEX: ", caigo.BigToHex(contentHash))
	pubX, pubY := data.StarkCurve.XToPubKey(sigKey)

	valid := data.StarkCurve.Verify(
		contentHash,
		caigo.StrToBig(rSig),
		caigo.StrToBig(sSig),
		pubX,
		pubY,
	)
	if !valid {
		httpError(fmt.Errorf("invalid signature"), "signature invalid", http.StatusBadRequest, w)
		return 
	}
	fmt.Println("IS VALID: ", valid)

	resp, err := data.CreateSubject(subject)
	if err != nil {
		httpError(err, "could not insert", http.StatusBadRequest, w)
		return
	}

	writeGoodJSON(resp, http.StatusCreated, w)
	return
}
