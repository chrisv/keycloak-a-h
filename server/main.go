package main

import (
	"log"

	"github.com/example-pipeline/keycloak/server/auth"

	"encoding/json"
	"net/http"
	"net/url"
)

var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		http.Error(w, "failed to get validated claims", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")
	enc.Encode(claims)
})

func main() {
	issuer, err := url.Parse("http://keycloak:41555/realms/file-api")
	if err != nil {
		log.Fatalf("failed to parse issuer URL: %v", err)
		return
	}
	authenticated, err := auth.NewMiddleware(issuer, handler)
	if err != nil {
		log.Fatalf("failed to create authentication middleware: %v", err)
		return
	}
	log.Printf("listening on port 3000...")
	http.ListenAndServe("127.0.0.1:3000", authenticated)
}
