package middleware

import (
	"api_gateway/api/token"
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, fmt.Errorf("authorization header is required"))
			return
		}

		valid, err := token.ValidateToken(auth)
		if err != nil || !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, fmt.Errorf("invalid token: %s", err))
			return
		}

		claims, err := token.ExtractClaims(auth)
		if err != nil || !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, fmt.Errorf("invalid token claims: %s", err))
			return
		}

		c.Set("claims", claims)
		fmt.Println(claims)
		c.Next()
	}
}

func CasbinMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		if enforcer == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Authorization service not initialized"})
			return
		}

		claims, exists := c.Get("claims")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		jwtClaims, ok := claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			return
		}

		sub, ok := jwtClaims["role"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid role"})
			return
		}

		obj := c.FullPath()
		act := c.Request.Method

		allowed, err := enforcer.Enforce(sub, obj, act)
		if err != nil {
			fmt.Println("++++++++++++++++++++---------000000000000", sub)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error occurred during authorization"})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}

		c.Next()
	}
}
