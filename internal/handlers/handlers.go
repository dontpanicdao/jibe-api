package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/jibe-api/internal/data"
)

func ElementsFetch(w http.ResponseWriter, r *http.Request) {
	elements, err := data.GetElements()
	if err != nil {
		httpError(err, "elements db pull", http.StatusInternalServerError, w)
		return
	}

	writeGoodJSON(elements, http.StatusOK, w)
	return
}

func ElementFetch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	element, err := data.GetElement(vars["element_id"])
	if err != nil {
		httpError(err, "element db pull", http.StatusInternalServerError, w)
		return
	}

	writeGoodJSON(element, http.StatusOK, w)
	return
}

func PhasesFetch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	phases, err := data.GetPhases(vars["element_id"])
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

	phase, err := data.GetPhase(vars["element_id"])
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

func CreateElement(w http.ResponseWriter, r *http.Request) {
	element := &data.Element{}
	err := json.NewDecoder(r.Body).Decode(&element)
	if err != nil {
		httpError(err, "could not parse json", http.StatusBadRequest, w)
		return
	}

	pubKey := r.Header.Get("Public-Key")
	sigKey := r.Header.Get("Signing-Key")
	rSig := r.Header.Get("Signature-R")
	sSig := r.Header.Get("Signature-S")

	tx := element.Transaction.ConvertTx()
	contentHash, err := data.StarkCurve.HashTx(caigo.HexToBN(pubKey), tx)
	if err != nil {
		httpError(err, "could not hash transaction", http.StatusBadRequest, w)
		return
	}

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

	resp, err := data.CreateElement(element)
	if err != nil {
		httpError(err, "could not insert", http.StatusBadRequest, w)
		return
	}

	writeGoodJSON(resp, http.StatusCreated, w)
	return
}
