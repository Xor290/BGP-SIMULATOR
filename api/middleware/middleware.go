package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"bgp-manager/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func clientSecret() []byte { return []byte(os.Getenv("USER_JWT_SECRET")) }

// validateClientToken valide un token client
func validateClientToken(tokenString string) (*models.ClientClaims, error) {
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return nil, fmt.Errorf("token vide")
	}

	token, err := jwt.ParseWithClaims(tokenString, &models.ClientClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("algorithme de signature inattendu: %v", token.Method.Alg())
		}
		return clientSecret(), nil
	})

	if err != nil {
		return nil, fmt.Errorf("erreur parsing JWT: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token invalide")
	}

	claims, ok := token.Claims.(*models.ClientClaims)
	if !ok {
		return nil, fmt.Errorf("type de claims invalide")
	}

	if claims.Role == "" {
		return nil, fmt.Errorf("role manquant dans le token")
	}

	if claims.Issuer != "api-client" {
		return nil, fmt.Errorf("issuer invalide")
	}

	if claims.ExpiresAt.Unix() < time.Now().Unix() {
		return nil, fmt.Errorf("token expiré")
	}

	return claims, nil
}

func ClientMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Printf("❌ [CLIENT-MWARE] Authorization header manquant")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token d'autorisation requis"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := validateClientToken(tokenStr)
		if err != nil {
			log.Printf("❌ [CLIENT-MWARE] Token invalide: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide"})
			c.Abort()
			return
		}

		c.Set("client_id", claims.ClientID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
