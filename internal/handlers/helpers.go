package handlers

import (
	"fmt"
	"encoding/json"
	"net/http"
)

type HTTPError struct {
	Error string `json:"error"`
	Metadata string `json:"metadata"`
}

func writeGoodJSON(payload []byte, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
	return
}

func httpError(err error, metadata string, code int, w http.ResponseWriter) error {
	if code < 400 {
		return fmt.Errorf("%v not an http error code", code)
	}
	erry := HTTPError{
		Error: err.Error(),
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