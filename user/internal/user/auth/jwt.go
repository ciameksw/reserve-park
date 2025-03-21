package auth

import (
	"time"

	"github.com/ciameksw/reserve-park/user/internal/user/mongodb"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Username string
	Role     mongodb.RoleType
	jwt.RegisteredClaims
}

func GenerateJWT(username string, role mongodb.RoleType, key string) (string, error) {
	claims := UserClaims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedJWT, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return signedJWT, nil
}
