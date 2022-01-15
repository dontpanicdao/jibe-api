package handlers

import (
	"errors"
	"fmt"
	"net/http"
)

type Auth struct {
}

func NewAuth() *Auth {
	return &Auth{}
}

func (auth *Auth) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Printf("%v %v: %v\n", r.Method, r.RequestURI, r.Host)
			next.ServeHTTP(w, r)
		} else {
			pubKey := r.Header.Get("Public-Key")
			sigKey := r.Header.Get("Signing-Key")
			rSig := r.Header.Get("Signature-R")
			sSig := r.Header.Get("Signature-S")
			if sigKey == "" || sigKey == "null" || rSig == "" || sSig == "" || pubKey == "" {
				httpError(errors.New("missing required post headers (pub key, rSig, sSig)"), "need header", http.StatusUnauthorized, w)
				return
			} else {
				fmt.Printf("%v %v: %v pub %v\n", r.Method, r.RequestURI, r.Host, pubKey)
				next.ServeHTTP(w, r)
			}
		}
	})
}
