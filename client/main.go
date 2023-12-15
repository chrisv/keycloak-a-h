package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Token struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	Scope            string `json:"scope"`
}

func getToken(endpoint, clientID, clientSecret, audience string) (token Token, err error) {
	resp, err := http.PostForm(endpoint, url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	})
	if err != nil {
		return token, fmt.Errorf("failed to get token from endpoint: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return token, fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return token, fmt.Errorf("failed to decode token response: %w", err)
	}
	return
}

func main() {
	endpoint := "http://keycloak:41555/realms/file-api/protocol/openid-connect/token"
	clientID := "file-api-sync"
	clientSecret := "ILyFLUdPtWSyl8jOZx8AIRKw5XjmJyGy"

	// Get the token.
	token, err := getToken(endpoint, clientID, clientSecret, "file-api")
	if err != nil {
		log.Fatalf("failed to get token: %v", err)
	}

	// Use it.
	url := "http://localhost:3000"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("failed to create request: %v", err)
	}
	req.Header.Add("authorization", "Bearer "+token.AccessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to make request: %v", err)
		return
	}
	fmt.Println("Printing body...")
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))
}
