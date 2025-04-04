package configs

import (
	"os"

	"github.com/joho/godotenv"
)

// GetMongoDBURI loads the MongoDB URI from environment variables
// Panics if the URI is not found
func GetMongoDBURI() string {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		panic("Failed to load .env file")
	}

	// Read MongoDB URI from environment variable
	mongoURI := os.Getenv("MONGO")

	// Crash the app if the MongoDB URI is missing
	if mongoURI == "" {
		panic("MONGODB_URI environment variable is not set")
	}

	return mongoURI
}
