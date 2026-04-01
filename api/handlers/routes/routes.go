package routes

import (
	"bgp-manager/db"
	"bgp-manager/handlers/clients"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, database *db.Database) {
	r.Use(func(c *gin.Context) {
		c.Set("database", database)
		c.Next()
	})

	authClient := r.Group("/api/v1/auth")
	{
		authClient.POST("/register", clients.RegisterClient)
		authClient.POST("/login", clients.LoginClient)
	}

}
