package bgp

import (
	"bgp-manager/db"
	"bgp-manager/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPeers(c *gin.Context) {
	database := c.MustGet("database").(*db.Database)
	peers, err := database.GetAllPeers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, peers)
}

func GetPeerWithId(c *gin.Context) {
	peerID := c.Param("peerID")
	database := c.MustGet("database").(*db.Database)
	peer, err := database.GetPeerById(peerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, peer)
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

	if err := database.DB.Preload("LocalAS").First(&peer, peer.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := ApplyPeerConfig(database, &peer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, peer)
}

func DeletePeer(c *gin.Context) {
	peerID := c.Param("peerID")
	database := c.MustGet("database").(*db.Database)
	if err := database.DeletePeerById(peerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "peer deleted"})
}

func GetSessions(c *gin.Context) {
	peerID := c.Param("peerID")
	database := c.MustGet("database").(*db.Database)
	sessions, err := database.GetSessions(peerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sessions)
}

func SyncSessions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "sync déclenché"})
}

func CreatePrefixSinceAS(c *gin.Context) {
	var req models.CreatePrefixSinceAS
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database := c.MustGet("database").(*db.Database)

	as, err := database.GetASByASN(req.ASN)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "AS introuvable"})
		return
	}

	prefix := models.PrefixSinceAS{
		Prefix:    req.Prefix,
		ASID:      as.ID,
		NextHop:   req.NextHop,
		LocalPref: req.LocalPref,
		Active:    true,
	}

	if err := database.CreatePrefixAS(&prefix); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := ApplyASConfig(database, as.ID, &as); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, prefix)
}
