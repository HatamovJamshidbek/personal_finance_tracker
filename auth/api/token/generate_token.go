package token

import (
	"auth_service/api/models"
	"auth_service/configs"
	"auth_service/pkg/logger"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/golang-jwt/jwt/v4"
	"time"
)

func GenerateJWTToken(logger logger.ILogger, user *models.LoginResponse) (*models.Token, error) {
	var (
		tokenValidateTime float64
		token             models.Token
	)
	accessToken := *jwt.New(jwt.SigningMethodHS256)
	accessClaims := accessToken.Claims.(jwt.MapClaims)
	accessClaims["id"] = user.Id
	accessClaims["user_name"] = user.UserName
	accessClaims["email"] = user.Email
	accessClaims["password_hash"] = user.PasswordHash
	accessClaims["role"] = user.Role
	accessClaims["iat"] = time.Now().Unix()
	accessClaims["exp"] = time.Now().Add(configs.AccessExpireTime).Unix()

	newAccessToken, err := accessToken.SignedString(configs.SignKey)
	if err != nil {
		logger.Error("this error is signed up to signed key is %")
		return nil, fmt.Errorf("this error is SignedString to SignKey")
	}
	refreshToken := *jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["id"] = user.Id
	refreshClaims["user_name"] = user.UserName
	refreshClaims["email"] = user.Email
	refreshClaims["password_hash"] = user.PasswordHash
	refreshClaims["role"] = user.Role
	refreshClaims["iat"] = time.Now().Unix()
	refreshClaims["exp"] = time.Now().Add(configs.RefreshExpireTime).Unix()
	tokenValidateTime = float64(refreshClaims["exp"].(int64) - accessClaims["iat"].(int64))
	newRefreshToken, err := refreshToken.SignedString(configs.SignKey)
	if err != nil {
		logger.Error("this error is signed to refresh token is ")
		return nil, fmt.Errorf("this error signed refresh token ")
	}
	token.AccessToken = newAccessToken
	token.RefreshToken = newRefreshToken
	token.ExpiredTime = tokenValidateTime
	return &token, nil
}
