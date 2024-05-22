package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtIssuerKey       = "JWT_ISSUER"
	jwtSignedStringKey = "JWT_SIGNED_STRING"
)

var (
	issuer = os.Getenv(jwtIssuerKey)
	key    = os.Getenv(jwtSignedStringKey)
	expire = 30 * 24 // 30 days
)

type Claims struct {
	ID uint
	jwt.RegisteredClaims
}

func GenerateToken(id uint) (string, error) {
	claims := Claims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expire) * time.Hour)),
			Issuer:    issuer,
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(key))
}

func ParseToken(token string) (uint, bool) {
	claims := new(Claims)
	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil || !t.Valid || claims == nil {
		return 0, false
	}
	return claims.ID, true
}
