package bgp

import (
	"bgp-manager/db"
	"bgp-manager/models"
)

func ApplyPeerConfig(database *db.Database, peer *models.Peer) error {
	peers, err := database.GetPeersByASID(peer.LocalASID)
	if err != nil {
		return err
	}
	return RunAnsiblePlaybook(&peer.LocalAS, peers)
}

func SyncSessionStates(database *db.Database) error {
	return nil
}

func RemoveAllPeer(database *db.Database, peerID uint) error {
	return database.DeletePeers()
}
