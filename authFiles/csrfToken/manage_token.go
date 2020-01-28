package csrfToken

import (
	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func Generate(signingKey []byte, claims jwt.Claims) (string, error) {
	tn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := tn.SignedString(signingKey)
	return signedString, err
}

func Valid(signedToken string, signingKey []byte) (bool, error) {
	token, err := jwt.ParseWithClaims(signedToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil {
		return false, err
	}

	if _, ok := token.Claims.(*CustomClaims); !ok || !token.Valid {
		return false, err
	}

	return true, nil
}

func Claims(email string, tokenExpires int64) jwt.Claims {
	return CustomClaims{
		email,
		jwt.StandardClaims{
			ExpiresAt: tokenExpires,
		},
	}
}

func CSRFToken(signingKey []byte) (string, error) {
	tn := jwt.New(jwt.SigningMethodHS256)
	signedString, err := tn.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return signedString, nil
}

func ValidCSRF(signedToken string, signingKey []byte) (bool, error) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil || !token.Valid {
		return false, err
	}

	return true, nil
}
