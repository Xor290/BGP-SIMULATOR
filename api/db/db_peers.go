package db

import "bgp-manager/models"

func (db *Database) CreatePeer(peer *models.Peer) error {
	return db.DB.Create(peer).Error
}
