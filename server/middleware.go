package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/iceokoli/get-crypto-balance/broker"
)

func loggingMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Printf("%s request for the endpoint %s", r.Method, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func authMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check Credentials
		apiKey := r.Header.Get("X-Auth-Key")
		confirmKey := strings.ToUpper(broker.GenerateSignature(perms.user, perms.secret))

		if apiKey == "" || confirmKey != apiKey {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"invalid_key"}`))
			return
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
