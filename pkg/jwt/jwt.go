package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type Config struct {
	Secret              string
	AccessTokenExpiry   time.Duration
	RefreshTokenExpiry  time.Duration
}

type JWTManager struct {
	secret             []byte
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func NewJWTManager(cfg Config) *JWTManager {
	return &JWTManager{
		secret:             []byte(cfg.Secret),
		accessTokenExpiry:  cfg.AccessTokenExpiry,
		refreshTokenExpiry: cfg.RefreshTokenExpiry,
	}
}

// GenerateAccessToken creates a signed JWT access token.
func (m *JWTManager) GenerateAccessToken(userID, email string) (string, time.Time, error) {
	expiresAt := time.Now().Add(m.accessTokenExpiry)

	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiresAt, nil
}

// GenerateRefreshToken creates a random refresh token (not a JWT).
func (m *JWTManager) GenerateRefreshToken() (string, time.Time, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", time.Time{}, err
	}

	token := base64.URLEncoding.EncodeToString(bytes)
	expiresAt := time.Now().Add(m.refreshTokenExpiry)

	return token, expiresAt, nil
}

// ValidateAccessToken parses and validates a JWT access token.
func (m *JWTManager) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// RefreshTokenExpiry returns the refresh token expiry duration.
func (m *JWTManager) RefreshTokenExpiry() time.Duration {
	return m.refreshTokenExpiry
}

// AccessTokenExpiry returns the access token expiry duration.
func (m *JWTManager) AccessTokenExpiry() time.Duration {
	return m.accessTokenExpiry
}
