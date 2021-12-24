package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func SubjectsFetch(w http.ResponseWriter, r *http.Request) {
	subjects, err := data.GetSubjects()
	if err != nil {
		httpError(err, "subjects db pull", http.StatusInternalServerError)
		return
	}

	writeGoodJSON(subjects, w)
	return
}

func SubjectFetch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	subject, err := data.GetSubject(vars["subject_id"])
	if err != nil {
		httpError(err, "subject db pull", http.StatusInternalServerError)
		return
	}

	writeGoodJSON(subjects, w)
	return
}

func PhasesFetch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	phases, err := data.GetPhases(vars["subject_id"])
	if err != nil {
		httpError(err, "phases db pull", http.StatusInternalServerError)
		return
	}

	writeGoodJSON(phases, w)
	return
}

func PhaseFetch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	phase, err := data.GetPhase(vars["phase_id"])
	if err != nil {
		httpError(err, "phase db pull", http.StatusInternalServerError)
		return
	}

	writeGoodJSON(phases, w)
	return
}

func CertFetch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	phase, err := data.GetPhase(vars["subject_id"])
	if err != nil {
		httpError(err, "phase db pull", http.StatusInternalServerError)
		return
	}

	writeGoodJSON(phases, w)
	return
}

func CertKey(w http.ResponseWriter, r *http.Request) {
	// TODO: handle the posting of an exam key
}