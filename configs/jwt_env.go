package configs

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type JWTKeys struct {
	Issuer   string
	Subject  string
	Id       string
	Audience string
	Secret   string
}

func GetJWTKeys() (*JWTKeys, error) {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		return nil, errors.New("error loading .env file")
	}

	// Read environment variables
	jwtKey := JWTKeys{
		Issuer:   os.Getenv("JWT_ISSUER"),
		Subject:  os.Getenv("JWT_SUBJECT"),
		Id:       os.Getenv("JWT_ID"),
		Audience: os.Getenv("JWT_AUDIENCE"),
		Secret:   os.Getenv("JWT_SECRET"),
	}

	// Check for missing values
	if jwtKey.Issuer == "" || jwtKey.Subject == "" || jwtKey.Id == "" || jwtKey.Audience == "" || jwtKey.Secret == "" {
		return nil, errors.New("missing required JWT environment variables")
	}

	return &jwtKey, nil
}
