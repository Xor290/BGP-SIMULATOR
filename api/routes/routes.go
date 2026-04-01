package routes

import (
	"bgp-manager/db"
	asHandler "bgp-manager/handlers/as"
	"bgp-manager/handlers/bgp"
	"bgp-manager/handlers/clients"
	"bgp-manager/middleware"

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

	peers := r.Group("/api/v1/peers")
	peers.Use(middleware.ClientMiddleware())
	{
		peers.GET("/all", bgp.GetPeers)
		peers.POST("/create", bgp.CreatePeer)
		peers.GET("/:peerID", bgp.GetPeerWithId)
		peers.DELETE("/:peerID", bgp.DeletePeer)
		peers.GET("/:peerID/sessions", bgp.GetSessions)
		peers.POST("/sync", bgp.SyncSessions)
	}

	asGroup := r.Group("/api/v1/as")
	asGroup.Use(middleware.ClientMiddleware())
	{
		asGroup.POST("/", asHandler.CreateAutonomousSystem)
	}
}
