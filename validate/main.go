package main

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/MicahParks/keyfunc/v2"
)

func main() {
	// Get the JWKS URL.
	//
	// This is a local Keycloak JWKS endpoint for the master realm.
	jwksURL := "http://keycloak:41555/realms/file-api/protocol/openid-connect/certs"

	// Create the keyfunc options. Use an error handler that logs. Refresh the JWKS when a JWT signed by an unknown KID
	// is found or at the specified interval. Rate limit these refreshes. Timeout the initial JWKS refresh request after
	// 10 seconds. This timeout is also used to create the initial context.Context for keyfunc.Get.
	options := keyfunc.Options{
		RefreshErrorHandler: func(err error) {
			log.Printf("There was an error with the jwt.Keyfunc\nError: %s", err.Error())
		},
		RefreshInterval:   time.Hour,
		RefreshRateLimit:  time.Minute * 5,
		RefreshTimeout:    time.Second * 10,
		RefreshUnknownKID: true,
	}

	// Create the JWKS from the resource at the given URL.
	jwks, err := keyfunc.Get(jwksURL, options)
	if err != nil {
		log.Fatalf("Failed to create JWKS from resource at the given URL.\nError: %s", err.Error())
	}
	defer jwks.EndBackground()

	// Get a JWT to parse.
	jwtB64 := "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJoWnpoS2RweHo4UnhCUjRCLUVkUDNwV1RMaGlWemlUMkNkdnBxR3h2TGtvIn0.eyJleHAiOjE3MDI2NDk2ODIsImlhdCI6MTcwMjY0OTM4MiwianRpIjoiNGE0ZDhkNDMtODY5Ni00NzUzLTk3ZWQtYjdiZmMxNmVkMTVlIiwiaXNzIjoiaHR0cDovL2tleWNsb2FrOjQxNTU1L3JlYWxtcy9maWxlLWFwaSIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiJiYjVkZjZmNC0zZDQwLTRkYjYtYWNiMi03ZDJhZmIyYzVlNzkiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJmaWxlLWFwaS1zeW5jIiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyJodHRwOi8vZmlsZS1hcGk6ODA4MCJdLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsib2ZmbGluZV9hY2Nlc3MiLCJkZWZhdWx0LXJvbGVzLWZpbGUtYXBpIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJmaWxlLWFwaS1zeW5jIjp7InJvbGVzIjpbInVtYV9wcm90ZWN0aW9uIl19LCJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJzY29wZSI6InByb2ZpbGUgZW1haWwiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImNsaWVudEhvc3QiOiIxMjcuMC4wLjEiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJzZXJ2aWNlLWFjY291bnQtZmlsZS1hcGktc3luYyIsImNsaWVudEFkZHJlc3MiOiIxMjcuMC4wLjEiLCJjbGllbnRfaWQiOiJmaWxlLWFwaS1zeW5jIn0.i7V2SbHRcySwQa8CE_r5zsEygU317PmblOt-6uf3zqc1JNUfPUZaaSM0k-fkdXf7p9S30UuJuMeta0sH8UuopcRSiWkS1s2kLwrRf7MD1TcJXdzq_7Wo0-kXa0CAVKjkHpiyPLpFyfUXWEiXv575W-vhZs4mChCvvVGy1KMe4b88yjGbi4zIKWxpkTB88U1bfNRuI9dLV7oMX5XEiC65k_s7EGhrtgw0A6q2rB2vS775Vdk3aCCLvavLElaE2II--PErFXrvblmF8MUIQvtgccKb9y3uWJ0XaQuiR2GlxAoa1VZBni7LH1Vm-0JqgzDqVz8ty0oL0OvMariY9ZvYHg"

	// Parse the JWT.
	token, err := jwt.Parse(jwtB64, jwks.Keyfunc)
	if err != nil {
		log.Fatalf("Failed to parse the JWT.\nError: %s", err.Error())
	}

	// Check if the token is valid.
	if !token.Valid {
		log.Fatalf("The token is not valid.")
	}
	log.Println("The token is valid.")
}
