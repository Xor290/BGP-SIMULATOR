package clients

import (
	"bgp-manager/db"
	"bgp-manager/models"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RegisterClient(c *gin.Context) {
	var req models.RegisterClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	database := c.MustGet("database").(*db.Database)

	if err := database.RegisterClient(&req); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Client enregistré"})
}

func LoginClient(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	database := c.MustGet("database").(*db.Database)

	client, err := database.ConnectClient(&req)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	claims := models.ClientClaims{
		ClientID: int(client.ID),
		Username: client.Username,
		Role:     "client",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "api-client",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(os.Getenv("USER_JWT_SECRET")))
	if err != nil {
		c.JSON(500, gin.H{"error": "génération du token échouée"})
		return
	}

	c.JSON(200, models.LoginResponse{
		AccessToken: signed,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		User:        client,
	})
}
