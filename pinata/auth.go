package pinata

import (
	"net/http"
)

// Auth represents the authentication credentials for the Pinata API.
// It can be used to authenticate requests with either an API key and secret,
// or a JWT token.
type Auth struct {
	apiKey    string
	apiSecret string
	jwt       string
}

// NewAuth creates a new Auth instance with the provided API key, API secret, and JWT token.
// The returned Auth instance can be used to authenticate requests to the Pinata API.
// If both an API key/secret and a JWT token are provided, the JWT token will take precedence.
func NewAuth(apiKey, apiSecret, jwt string) *Auth {
	return &Auth{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		jwt:       jwt,
	}
}

// NewAuthWithJWT creates a new Auth instance with the provided JWT token.
// The returned Auth instance can be used to authenticate requests to the Pinata API.
// If both an API key/secret and a JWT token are provided, the JWT token will take precedence.
func NewAuthWithJWT(jwt string) *Auth {
	return &Auth{
		jwt: jwt,
	}
}

// setAuthHeader sets the appropriate authentication headers on the provided HTTP request.
// If a JWT token is provided, it sets the Authorization header to "Bearer <JWT>".
// Otherwise, it sets the pinata_api_key and pinata_secret_api_key headers with the provided API key and secret.
func (a *Auth) setAuthHeader(req *http.Request) {
	if a.jwt != "" {
		req.Header.Set("Authorization", "Bearer "+a.jwt)
		return
	}
	req.Header.Set("pinata_api_key", a.apiKey)
	req.Header.Set("pinata_secret_api_key", a.apiSecret)
}
