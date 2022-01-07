package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

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
