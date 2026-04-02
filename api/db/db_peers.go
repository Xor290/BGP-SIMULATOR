package db

import (
	"bgp-manager/models"

	"gorm.io/gorm"
)

func (db *Database) CreatePeer(peer *models.Peer) error {
	return db.DB.Create(peer).Error
}

func (db *Database) DeletePeerById(peerID string) error {
	return db.DB.Delete(&models.Peer{}, peerID).Error
}

func (db *Database) GetSessions(peerID string) ([]models.BGPSession, error) {
	var sessions []models.BGPSession
	return sessions, db.DB.Where("peer_id = ?", peerID).Find(&sessions).Error
}

func (db *Database) GetAllPeers() ([]models.Peer, error) {
	var peers []models.Peer
	return peers, db.DB.Find(&peers).Error
}

func (db *Database) GetPeerById(peerID string) (*models.Peer, error) {
	var peer models.Peer
	return &peer, db.DB.Preload("LocalAS").First(&peer, peerID).Error
}

func (db *Database) GetPeersByASID(asID uint) ([]models.Peer, error) {
	var peers []models.Peer
	return peers, db.DB.Where("local_as_id = ?", asID).Find(&peers).Error
}

func (db *Database) DeletePeers() error {
	return db.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Peer{}).Error
}

func (db *Database) CreatePrefixAS(prefix *models.PrefixSinceAS) error {
	return db.DB.Create(prefix).Error
}
