package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dontpanicdao/jibe-api/internal/data"
	"github.com/gorilla/mux"
)

func SubjectsFetch(w http.ResponseWriter, r *http.Request) {
	subjects, err := data.GetSubjects()
	if err != nil {
		httpError(err, "subjects db pull", http.StatusInternalServerError, w)
		return
	}

	writeGoodJSON(subjects, w)
	return
}

func SubjectFetch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	subject, err := data.GetSubject(vars["subject_id"])
	if err != nil {
		httpError(err, "subject db pull", http.StatusInternalServerError, w)
		return
	}

	writeGoodJSON(subject, w)
	return
}

func PhasesFetch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	phases, err := data.GetPhases(vars["subject_id"])
	if err != nil {
		httpError(err, "phases db pull", http.StatusInternalServerError, w)
		return
	}

	writeGoodJSON(phases, w)
	return
}

func PhaseFetch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	phase, err := data.GetPhase(vars["phase_id"])
	if err != nil {
		httpError(err, "phase db pull", http.StatusInternalServerError, w)
		return
	}

	writeGoodJSON(phase, w)
	return
}

func CertFetch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	phase, err := data.GetPhase(vars["subject_id"])
	if err != nil {
		httpError(err, "phase db pull", http.StatusInternalServerError, w)
		return
	}

	writeGoodJSON(phase, w)
	return
}

func CertKey(w http.ResponseWriter, r *http.Request) {
	// TODO: handle the posting of an exam key
}

func CreateSubject(w http.ResponseWriter, r *http.Request) {
	subject := &data.TypedSubject{}
	json.NewDecoder(r.Body).Decode(&subject)

	err := data.CreateSubject(subject)
	if err != nil {
		httpError(err, "could not insert", http.StatusBadRequest, w)
		return
	}
	fmt.Println("R BODY: ", r.Body)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "congrats homie")
	return
}
