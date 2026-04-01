package as

import (
	"bgp-manager/db"
	"bgp-manager/models"

	"github.com/gin-gonic/gin"
)

func CreateAutonomousSystem(c *gin.Context) {
	var req models.AutonomousSystem
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	database := c.MustGet("database").(*db.Database)
	if err := database.CreateAS(req); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, req)
}
