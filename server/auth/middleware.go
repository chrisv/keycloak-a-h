package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

type Access struct {
	Roles []string `json:"roles"`
}

// KeycloakClaims contains custom data we want from the token.
type KeycloakClaims struct {
	Sub               string            `json:"sub"`
	RealmAccess       Access            `json:"realm_access"`
	ResourceAccess    map[string]Access `json:"resource_access"`
	ClientAddress     string            `json:"clientAddress"`
	ClientID          string            `json:"client_id"`
	Scope             string            `json:"scope"`
	PreferredUsername string            `json:"preferred_username"`
}

func GetClaims(ctx context.Context) (claims *KeycloakClaims, ok bool) {
	vc := ctx.Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	claims, ok = vc.CustomClaims.(*KeycloakClaims)
	return
}

// HasScope checks whether our claims have a specific scope.
func (c *KeycloakClaims) HasScope(expectedScope string) bool {
	result := strings.Split(c.Scope, " ")
	for i := range result {
		if result[i] == expectedScope {
			return true
		}
	}

	return false
}

// Satisfy validator.CustomClaims interface.
func (c *KeycloakClaims) Validate(ctx context.Context) error {
	return nil
}

// EnsureValidToken is middleware that will check the validity of the JWT.
func NewMiddleware(issuer *url.URL, next http.Handler) (h http.Handler, err error) {
	jwtValidator, err := validator.New(
		jwks.NewCachingProvider(issuer, 30*time.Minute).KeyFunc,
		validator.RS256,
		issuer.String(),
		[]string{"account"},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &KeycloakClaims{}
			},
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create new JWT validator: %w", err)
	}
	return jwtmiddleware.New(jwtValidator.ValidateToken).CheckJWT(next), nil
}
