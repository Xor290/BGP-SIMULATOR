package bgp

import (
	"bgp-manager/db"
	"bgp-manager/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPeers(c *gin.Context) {

}

func CreatePeer(c *gin.Context) {
	var req models.CreatePeerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	peer := models.Peer{
		LocalASID:   req.LocalASID,
		RemoteASN:   req.RemoteASN,
		PeerIP:      req.PeerIP,
		Description: req.Description,
		Password:    req.Password,
		Enabled:     enabled,
	}

	database := c.MustGet("database").(*db.Database)
	if err := database.CreatePeer(&peer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, peer)
}

func DeletePeer(c *gin.Context) {
}

func GetSessions(c *gin.Context) {
}

func SyncSessions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "sync déclenché"})
}
