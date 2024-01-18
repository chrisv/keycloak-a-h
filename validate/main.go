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
	jwtB64 := "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJoWnpoS2RweHo4UnhCUjRCLUVkUDNwV1RMaGlWemlUMkNkdnBxR3h2TGtvIn0.eyJleHAiOjE3MDI5MDYyNzcsImlhdCI6MTcwMjkwNTk3NywianRpIjoiMWU0ZDQ3NjAtNGJkYi00NGExLWFjZWYtNGY4NjJiZjA0ZjRkIiwiaXNzIjoiaHR0cDovL2tleWNsb2FrOjQxNTU1L3JlYWxtcy9maWxlLWFwaSIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiJiYjVkZjZmNC0zZDQwLTRkYjYtYWNiMi03ZDJhZmIyYzVlNzkiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJmaWxlLWFwaS1zeW5jIiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyJodHRwOi8vZmlsZS1hcGk6ODA4MCJdLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsib2ZmbGluZV9hY2Nlc3MiLCJkZWZhdWx0LXJvbGVzLWZpbGUtYXBpIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJmaWxlLWFwaS1zeW5jIjp7InJvbGVzIjpbInVtYV9wcm90ZWN0aW9uIl19LCJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJzY29wZSI6InByb2ZpbGUgZW1haWwiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImNsaWVudEhvc3QiOiIxMjcuMC4wLjEiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJzZXJ2aWNlLWFjY291bnQtZmlsZS1hcGktc3luYyIsImNsaWVudEFkZHJlc3MiOiIxMjcuMC4wLjEiLCJjbGllbnRfaWQiOiJmaWxlLWFwaS1zeW5jIn0.XzWO7po6CllN9AP2xEBDE2QjfZ4z_DICDmgXe2aHQSbnXLTfNv5Gs5fYE3CGWdXgN6wBYZGeNIj7Ig_IZBe0EGwva7YQEMIyaIzdrkPW_RHHo2yZMg9GmbK540KgXxm3B1cglSJ2VdkywK5wFN_n-88tIzGI0ueQ4WfR-kY8rW597jtjBi1S6iZmNYJv3Seu93oDc26VwpgHTahZxaV7OCijhcpqIApqOmqkz9KdYD68arWBbkeXyGXLyB1pA6ZtHqkdAqN7sLb7Q0V9wVfRGDSy8M6xdix8svCG_9Upw-A_qHuC285upDS3KsIq893tppXC4nX16UHUxB6zGL9BdQ"

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
