package handlers

import (
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/gorilla/mux"
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

func ProtonsFetch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	protons, err := data.GetProtons(vars["element_id"])
	if err != nil {
		httpError(err, "protons db pull", http.StatusInternalServerError, w)
		return
	}

	writeGoodJSON(protons, http.StatusOK, w)
	return
}

func ProtonFetch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	proton, err := data.GetProton(vars["proton_id"])
	if err != nil {
		httpError(err, "proton db pull", http.StatusInternalServerError, w)
		return
	}

	writeGoodJSON(proton, http.StatusOK, w)
	return
}

func FetchCustomCert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	cert, err := data.GetCustomCert(vars["element_id"])
	if err != nil {
		httpError(err, "proton db pull", http.StatusInternalServerError, w)
		return
	}

	writeGoodJSON(cert, http.StatusOK, w)
	return
}

func UserFetch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	user, err := data.GetUser(vars["public_key"])
	if err != nil {
		httpError(err, "user db pull", http.StatusInternalServerError, w)
		return
	}

	writeGoodJSON(user, http.StatusOK, w)
	return
}

func CheckFactJobStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	fact, err := data.CheckFactJob(vars["job_id"])
	if err != nil {
		httpError(err, "fact job db pull", http.StatusInternalServerError, w)
		return
	}

	writeGoodJSON(fact, http.StatusOK, w)
	return
}

func CheckFactStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	fact, err := data.CheckFact(vars["fact"])
	if err != nil {
		httpError(err, "fact db pull", http.StatusInternalServerError, w)
		return
	}

	writeGoodJSON(fact, http.StatusOK, w)
	return
}

func ShipSharpStark(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	fact, err := data.ShipSharpStark(vars["fact"], vars["tx_hash"])
	if err != nil {
		httpError(err, "ship sharp db pull", http.StatusInternalServerError, w)
		return
	}

	writeGoodJSON(fact, http.StatusOK, w)
	return
}

func AddProton(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	proton := &data.Proton{}
	err := json.NewDecoder(r.Body).Decode(&proton)
	if err != nil {
		httpError(err, "could not parse json", http.StatusBadRequest, w)
		return
	}

	is_valid := proton.Verify(
		r.Header.Get("Public-Key"),
		r.Header.Get("Signing-Key"),
		r.Header.Get("Signature-R"),
		r.Header.Get("Signature-S"),
		vars["element_id"],
	)
	if !is_valid {
		httpError(fmt.Errorf("proton is invalid"), "invalid signature", http.StatusBadRequest, w)
		return
	}

	resp, err := proton.Create(vars["element_id"])
	if err != nil {
		httpError(err, "could not write proton to db", http.StatusInternalServerError, w)
		return
	}
	writeGoodJSON(resp, http.StatusCreated, w)
	return
}

func CreateElement(w http.ResponseWriter, r *http.Request) {
	element := &data.Element{}
	err := json.NewDecoder(r.Body).Decode(&element)
	if err != nil {
		httpError(err, "could not parse json", http.StatusBadRequest, w)
		return
	}

	hash, valid := VerifyTx(element.Transaction, r)
	fmt.Println("IS VALID: ", valid)
	if !valid {
		httpError(fmt.Errorf("invalid signature"), "signature invalid", http.StatusBadRequest, w)
		return
	}

	resp, err := data.CreateElement(element, hash)
	if err != nil {
		httpError(err, "could not insert", http.StatusBadRequest, w)
		return
	}

	writeGoodJSON(resp, http.StatusCreated, w)
	return
}

func CreateCustomCert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	questions := data.Attrs{}
	err := json.NewDecoder(r.Body).Decode(&questions)
	if err != nil {
		httpError(err, "could not parse json", http.StatusBadRequest, w)
		return
	}

	resp, err := data.CreateQuestions(questions, vars["element_id"])
	if err != nil {
		httpError(err, "could not write cert to db", http.StatusInternalServerError, w)
		return
	}

	writeGoodJSON(resp, http.StatusCreated, w)
	return
}

func ProposeCert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	cert := &data.Cert{}
	err := json.NewDecoder(r.Body).Decode(&cert)
	if err != nil {
		httpError(err, "could not parse json", http.StatusBadRequest, w)
		return
	}

	is_valid := cert.Verify(
		r.Header.Get("Public-Key"),
		r.Header.Get("Signing-Key"),
		r.Header.Get("Signature-R"),
		r.Header.Get("Signature-S"),
		vars["element_id"],
	)
	if !is_valid {
		httpError(fmt.Errorf("cert is invalid"), "invalid signature", http.StatusBadRequest, w)
		return
	}

	resp, err := cert.Create(vars["element_id"])
	if err != nil {
		httpError(err, "could not write cert to db", http.StatusInternalServerError, w)
		return
	}

	writeGoodJSON(resp, http.StatusCreated, w)
	return
}

func AddRubric(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	cert := &data.Cert{}
	err := json.NewDecoder(r.Body).Decode(&cert)
	if err != nil {
		httpError(err, "could not parse json", http.StatusBadRequest, w)
		return
	}

	is_valid := cert.Verify(
		r.Header.Get("Public-Key"),
		r.Header.Get("Signing-Key"),
		r.Header.Get("Signature-R"),
		r.Header.Get("Signature-S"),
		vars["element_id"],
	)
	if !is_valid {
		httpError(fmt.Errorf("cert is invalid"), "invalid signature", http.StatusBadRequest, w)
		return
	}

	resp, err := cert.AddRubric(vars["element_id"])
	if err != nil {
		httpError(err, "could not add rubric", http.StatusInternalServerError, w)
		return
	}

	writeGoodJSON(resp, http.StatusOK, w)
	return
}

func GradeCert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	cert := &data.Cert{}
	err := json.NewDecoder(r.Body).Decode(&cert)
	if err != nil {
		httpError(err, "could not parse json", http.StatusBadRequest, w)
		return
	}

	pubKey := r.Header.Get("Public-Key")

	is_valid := cert.VerifyAttempt(
		r.Header.Get("Public-Key"),
		r.Header.Get("Signing-Key"),
		r.Header.Get("Signature-R"),
		r.Header.Get("Signature-S"),
		vars["element_id"],
	)
	if !is_valid {
		httpError(fmt.Errorf("cert is invalid"), "invalid signature", http.StatusBadRequest, w)
		return
	}
	resp, err := cert.Grade(vars["element_id"], pubKey)
	if err != nil {
		httpError(err, "could not grade and submit exam", http.StatusInternalServerError, w)
		return
	}

	writeGoodJSON(resp, http.StatusCreated, w)
	return
}
