package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/simonaditia/nyayurin/user-service/config"
)

type JwtServiceInterface interface {
	GenerateToken(userID int64) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
	secretKey string
	issuer    int64
}

// GenerateToken implements JwtServiceInterface.
func (j *jwtService) GenerateToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"iss":     j.issuer,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires after 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// ValidateToken implements JwtServiceInterface.
func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return []byte(j.secretKey), nil
	})
}

func NewJwtService(cfg *config.Config) JwtServiceInterface {
	return &jwtService{
		secretKey: cfg.App.JwtSecretKey,
		issuer:    cfg.App.JwtIssuer,
	}
}
