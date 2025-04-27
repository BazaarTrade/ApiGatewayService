package tokenManager

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidTokenClaims      = errors.New("invalid token claims")
)

func GenerateAccessToken(userID int, ACCESS_TOKEN_SECRET_PHRASE string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(time.Minute * 1).Unix(), // 15 minutes expiration
	})

	signedToken, err := token.SignedString([]byte(ACCESS_TOKEN_SECRET_PHRASE))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateAccessToken(accessToken, ACCESS_TOKEN_SECRET_PHRASE string) (bool, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		return []byte(ACCESS_TOKEN_SECRET_PHRASE), nil
	})

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check expiration time
		exp, ok := claims["exp"].(float64)
		if !ok {
			return false, ErrInvalidTokenClaims
		}
		expTime := time.Unix(int64(exp), 0)
		if time.Now().After(expTime) {
			return false, nil
		}

		return true, nil
	}

	return false, ErrInvalidTokenClaims
}
