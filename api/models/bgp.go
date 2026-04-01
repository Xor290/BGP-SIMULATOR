package models

import (
	"time"

	"gorm.io/gorm"
)

// AutonomousSystem représente un AS local géré par le système.
type AutonomousSystem struct {
	gorm.Model
	ASN         uint32 `gorm:"uniqueIndex;not null" json:"asn"`
	Name        string `gorm:"size:128;not null"   json:"name"`
	RouterID    string `gorm:"size:15;not null"    json:"router_id"`   // ex: "10.0.0.1"
	Description string `gorm:"size:255"            json:"description"`
	Peers       []Peer `gorm:"foreignKey:LocalASID" json:"peers,omitempty"`
}

// Peer représente un voisin BGP configuré sur un AS local.
type Peer struct {
	gorm.Model
	LocalASID   uint             `gorm:"not null;index"      json:"local_as_id"`
	LocalAS     AutonomousSystem `gorm:"foreignKey:LocalASID" json:"local_as,omitempty"`
	RemoteASN   uint32           `gorm:"not null"            json:"remote_asn"`
	PeerIP      string           `gorm:"size:45;not null"    json:"peer_ip"`  // IPv4 ou IPv6
	Description string           `gorm:"size:255"            json:"description"`
	Password    string           `gorm:"size:128"            json:"-"`        // MD5 auth, jamais exposé
	Enabled     bool             `gorm:"default:true"        json:"enabled"`
	Session     *BGPSession      `gorm:"foreignKey:PeerID"   json:"session,omitempty"`
}

// BGPSessionState est l'état FSM d'une session BGP (RFC 4271).
type BGPSessionState string

const (
	StateIdle        BGPSessionState = "Idle"
	StateConnect     BGPSessionState = "Connect"
	StateActive      BGPSessionState = "Active"
	StateOpenSent    BGPSessionState = "OpenSent"
	StateOpenConfirm BGPSessionState = "OpenConfirm"
	StateEstablished BGPSessionState = "Established"
)

// BGPSession représente l'état en temps réel d'une session BGP.
type BGPSession struct {
	gorm.Model
	PeerID           uint            `gorm:"uniqueIndex;not null"           json:"peer_id"`
	State            BGPSessionState `gorm:"size:16;not null;default:'Idle'" json:"state"`
	UpSince          *time.Time      `json:"up_since,omitempty"`
	PrefixesReceived uint            `gorm:"default:0"                      json:"prefixes_received"`
	PrefixesSent     uint            `gorm:"default:0"                      json:"prefixes_sent"`
	LastError        string          `gorm:"size:255"                       json:"last_error,omitempty"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

// BGPRoute représente une route annoncée ou reçue via BGP.
type BGPRoute struct {
	gorm.Model
	PeerID      uint   `gorm:"not null;index"   json:"peer_id"`
	Peer        Peer   `gorm:"foreignKey:PeerID" json:"peer,omitempty"`
	Prefix      string `gorm:"size:49;not null" json:"prefix"`     // ex: "10.0.0.0/8"
	NextHop     string `gorm:"size:45;not null" json:"next_hop"`
	LocalPref   uint   `gorm:"default:100"      json:"local_pref"`
	MED         uint   `gorm:"default:0"        json:"med"`
	ASPath      string `gorm:"size:255"         json:"as_path"`    // ex: "65001 65002"
	Communities string `gorm:"size:255"         json:"communities"` // ex: "65001:100"
	Origin      string `gorm:"size:4"           json:"origin"`     // "IGP", "EGP", "?"
	Direction   string `gorm:"size:8;not null"  json:"direction"`  // "received" | "advertised"
	Best        bool   `gorm:"default:false"    json:"best"`
}
