package bgp

import (
	"bgp-manager/db"
	"bgp-manager/models"
)

func ApplyPeerConfig(database *db.Database, peer *models.Peer) error {
	return nil
}

func SyncSessionStates(database *db.Database) error {
	return nil
}

func RemovePeer(database *db.Database, peerID uint) error {
	return nil
}
