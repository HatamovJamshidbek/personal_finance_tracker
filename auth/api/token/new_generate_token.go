package token

import (
	"auth_service/api/models"
	"auth_service/configs"
	"auth_service/pkg/logger"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

func NewGenerateToken(logger logger.ILogger, tokenStr string) (*models.LoginResponse, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Error("this error is toke method is not matching to method of SigningMethodHMAC")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return configs.SignKey, nil
	})
	if err != nil {
		logger.Error("error parsing token:")
		return nil, fmt.Errorf("error parsing token: %v", err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := &models.LoginResponse{
			UserName:     claims["user_name"].(string),
			Email:        claims["email"].(string),
			PasswordHash: claims["password_hash"].(string),
			Role:         claims["role"].(string),
		}
		return user, nil
	} else {
		logger.Error("invalid token")
		return nil, fmt.Errorf("invalid token")
	}
}
func ExtractClaim(tokenStr string) (jwt.MapClaims, error) {
	var (
		token *jwt.Token
		err   error
	)

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.SignKey), nil
	}
	token, err = jwt.Parse(tokenStr, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, errors.New("token invalid token")
	}

	return claims, nil
}
