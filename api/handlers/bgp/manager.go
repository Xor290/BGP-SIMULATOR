package bgp

import (
	"bgp-manager/db"
	"bgp-manager/models"
	"fmt"
)

func ApplyPeerConfig(database *db.Database, peer *models.Peer) error {
	if err := RunVtyshCommand(peer.LocalAS.ASN, "configure terminal\n"+
		"router bgp "+fmt.Sprint(peer.LocalAS.ASN)+"\n"+
		"neighbor "+peer.PeerIP+" remote-as "+fmt.Sprint(peer.RemoteASN)+"\n"+
		"exit\n"+
		"write memory"); err != nil {
		return err
	}
	return nil
}

func SyncSessionStates(database *db.Database) error {
	return nil
}

func RemovePeer(database *db.Database, peerID uint) error {
	return nil
}
