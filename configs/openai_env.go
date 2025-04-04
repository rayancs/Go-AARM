package configs

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type OpenAiKeys struct {
	Key string
}

func GetOpenAiKeys() (*OpenAiKeys, error) {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		return nil, errors.New("error loading .env file")
	}

	// Read environment variable
	openAiKey := OpenAiKeys{
		Key: os.Getenv("OPENAI"),
	}

	// Check for missing value
	if openAiKey.Key == "" {
		return nil, errors.New("missing required OpenAI API key environment variable")
	}

	return &openAiKey, nil
}
