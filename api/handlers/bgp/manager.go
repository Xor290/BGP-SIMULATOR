package bgp

import (
	"bgp-manager/db"
	"bgp-manager/models"
)

func ApplyPeerConfig(database *db.Database, peer *models.Peer) error {
	return ApplyASConfig(database, peer.LocalASID, &peer.LocalAS)
}

func ApplyASConfig(database *db.Database, asID uint, localAS *models.AutonomousSystem) error {
	peers, err := database.GetPeersByASID(asID)
	if err != nil {
		return err
	}
	prefixes, err := database.GetPrefixesByASID(asID)
	if err != nil {
		return err
	}
	return RunAnsiblePlaybook(localAS, peers, prefixes)
}

func SyncSessionStates(database *db.Database) error {
	return nil
}

func RemoveAllPeer(database *db.Database, peerID uint) error {
	return database.DeletePeers()
}
