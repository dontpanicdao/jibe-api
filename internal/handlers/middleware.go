package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dontpanicdao/jibe-api/pkg/caigo"
)

type Auth struct {
}

func NewAuth() *Auth {
	return &Auth{}
}

func (auth *Auth) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Println("GOT REQUEST FOR GET")
			next.ServeHTTP(w, r)
		} else {
			//headers for "pubkey" "signature-r" "signature-s"
			sigKey := r.Header.Get("Signing-Key")
			rSig := r.Header.Get("Signature-R")
			sSig := r.Header.Get("Signature-S")
			hash := r.Header.Get("Content-Hash")
			if sigKey == "" || rSig == "" || sSig == "" || hash == "" {
				httpError(errors.New("missing required post headers (pub key, rSig, sSig, hash)"), sigKey, http.StatusUnauthorized, w)
				return
			}
			fmt.Println("GOT REQUEST FOR POST: ", sigKey, rSig, sSig, hash)
			valid := caigo.Verify(
				caigo.HexToBN(hash),
				caigo.StrToBig(rSig),
				caigo.StrToBig(sSig),
				caigo.XToPubKey(sigKey),
				caigo.SC(),
			)
			fmt.Println("VALID: ", valid)
			if !valid {
				httpError(errors.New("invalid signature"), sigKey, http.StatusUnauthorized, w)
				return
			} else {
				next.ServeHTTP(w, r)
			}
		}
	})
}
