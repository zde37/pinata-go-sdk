package pinata

import (
	"net/http"
)

type Auth struct {
	APIKey    string
	APISecret string
	JWT       string
}

func NewAuth(apiKey, apiSecret, jwt string) *Auth {
	return &Auth{
		APIKey:    apiKey,
		APISecret: apiSecret,
		JWT:       jwt,
	}
}

func NewAuthWithJWT(jwt string) *Auth {
	return &Auth{
		JWT: jwt,
	}
}

func (a *Auth) setAuthHeader(req *http.Request) {
	if a.JWT != "" {
		req.Header.Set("Authorization", "Bearer "+a.JWT)
		return
	}
	req.Header.Set("pinata_api_key", a.APIKey)
	req.Header.Set("pinata_secret_api_key", a.APISecret)
}