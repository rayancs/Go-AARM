package services

import (
	"app/configs"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTservice struct {
}
type Claims struct {
	UserID   string `json:"id"`
	Username string `json:"name"`
	Emoji    string `json:"emoji"`
	Email    string `json:"email"`

	jwt.RegisteredClaims
}

func NewJWTService() *JWTservice {
	return &JWTservice{}
}
func (j *JWTservice) CreateAuthToken() (string, error) {
	return "", nil
}
func CreateToken(uId, name, emoji, email string) (string, error) {
	jwtKeyEnv, err := configs.GetJWTKeys()
	if err != nil {
		return "", err
	}

	secretKey := []byte(jwtKeyEnv.Secret)
	claims := Claims{
		UserID:   uId,
		Emoji:    emoji,
		Email:    email,
		Username: name,
		RegisteredClaims: jwt.RegisteredClaims{
			// Set standard claims
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    jwtKeyEnv.Issuer,
			Subject:   jwtKeyEnv.Subject,
			ID:        jwtKeyEnv.Id,
			Audience:  []string{"App"},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// VerifyToken validates a JWT token and returns the claims if valid
func VerifyToken(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, errors.New("token is empty")
	}

	jwtKeyEnv, err := configs.GetJWTKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to get JWT keys: %w", err)
	}

	secretKey := []byte(jwtKeyEnv.Secret)

	// Parse the token
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			// Validate the algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Check if token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("failed to extract claims")
	}

	// Check if token is expired
	expirationTime, err := claims.GetExpirationTime()
	if err != nil {
		return nil, fmt.Errorf("failed to get expiration time: %w", err)
	}

	if expirationTime.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	// Validate audience if needed
	aud, err := claims.GetAudience()
	if err != nil {
		return nil, fmt.Errorf("failed to get audience: %w", err)
	}

	validAudience := false
	for _, a := range aud {
		if a == "App" {
			validAudience = true
			break
		}
	}

	if !validAudience {
		return nil, errors.New("invalid audience")
	}

	// Validate issuer
	if iss, err := claims.GetIssuer(); err != nil || iss != jwtKeyEnv.Issuer {
		return nil, errors.New("invalid issuer")
	}

	return claims, nil
}
