 package configs

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// read configs from env variables
type AuthRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
}

func GetGCPCreds() (*AuthRequest, error) {

	if err := godotenv.Load(); err != nil {
		return nil, errors.New("error loading .env file")
	}

	authReq := AuthRequest{
		ClientID:     os.Getenv("client_id"),
		ClientSecret: os.Getenv("client_secret"),
		RedirectURI:  os.Getenv("redirect_uri"),
	}

	if authReq.ClientID == "" || authReq.ClientSecret == "" || authReq.RedirectURI == "" {
		return nil, errors.New("missing required environment variables")
	}

	return &authReq, nil
}
