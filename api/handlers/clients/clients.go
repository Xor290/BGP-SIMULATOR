package clients

import (
	"bgp-manager/db"
	"bgp-manager/models"

	"github.com/gin-gonic/gin"
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

	c.JSON(200, gin.H{"message": "Connexion réussie", "client": client})
}
