package models

import (
	"time"

	"gorm.io/gorm"
)

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

// Peer représente un voisin BGP configuré sur un AS local.
type Peer struct {
	gorm.Model
	LocalASID   uint             `gorm:"not null;index"       json:"local_as_id"`
	LocalAS     AutonomousSystem `gorm:"foreignKey:LocalASID" json:"local_as,omitempty"`
	RemoteASN   uint32           `gorm:"not null"             json:"remote_asn"`
	PeerIP      string           `gorm:"size:45;not null"     json:"peer_ip"`
	Description string           `gorm:"size:255"             json:"description"`
	Password    string           `gorm:"size:128"             json:"-"`
	Enabled     bool             `gorm:"default:true"         json:"enabled"`
	Session     *BGPSession      `gorm:"foreignKey:PeerID"    json:"session,omitempty"`
}

// BGPSession représente l'état en temps réel d'une session BGP.
type BGPSession struct {
	gorm.Model
	PeerID           uint            `gorm:"uniqueIndex;not null"            json:"peer_id"`
	State            BGPSessionState `gorm:"size:16;not null;default:'Idle'" json:"state"`
	UpSince          *time.Time      `json:"up_since,omitempty"`
	PrefixesReceived uint            `gorm:"default:0"                       json:"prefixes_received"`
	PrefixesSent     uint            `gorm:"default:0"                       json:"prefixes_sent"`
	LastError        string          `gorm:"size:255"                        json:"last_error,omitempty"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

// BGPRoute représente une route annoncée ou reçue via BGP.
type BGPRoute struct {
	gorm.Model
	PeerID      uint   `gorm:"not null;index"    json:"peer_id"`
	Peer        Peer   `gorm:"foreignKey:PeerID" json:"peer,omitempty"`
	Prefix      string `gorm:"size:49;not null"  json:"prefix"`
	NextHop     string `gorm:"size:45;not null"  json:"next_hop"`
	LocalPref   uint   `gorm:"default:100"       json:"local_pref"`
	MED         uint   `gorm:"default:0"         json:"med"`
	ASPath      string `gorm:"size:255"          json:"as_path"`
	Communities string `gorm:"size:255"          json:"communities"`
	Origin      string `gorm:"size:4"            json:"origin"`    // "IGP", "EGP", "?"
	Direction   string `gorm:"size:8;not null"   json:"direction"` // "received" | "advertised"
	Best        bool   `gorm:"default:false"     json:"best"`
}

// --- DTOs ---

// CreatePeerRequest est le body attendu pour créer un nouveau peer BGP.
type CreatePeerRequest struct {
	LocalASID   uint   `json:"local_as_id" binding:"required"`
	RemoteASN   uint32 `json:"remote_asn"  binding:"required"`
	PeerIP      string `json:"peer_ip"     binding:"required,ip"`
	Description string `json:"description"`
	Password    string `json:"password"`
	Enabled     *bool  `json:"enabled"` // nil → défaut true (cohérent avec gorm:"default:true")
}

// UpdatePeerRequest est le body attendu pour modifier un peer existant.
type UpdatePeerRequest struct {
	Description string `json:"description"`
	Password    string `json:"password"`
	Enabled     *bool  `json:"enabled"` // pointeur pour distinguer false et absent
}

// PeerResponse est la représentation d'un peer retournée par l'API.
type PeerResponse struct {
	ID          uint             `json:"id"`
	LocalASN    uint32           `json:"local_asn"`
	RemoteASN   uint32           `json:"remote_asn"`
	PeerIP      string           `json:"peer_ip"`
	Description string           `json:"description"`
	Enabled     bool             `json:"enabled"`
	Session     *SessionResponse `json:"session,omitempty"`
}

// SessionResponse est l'état d'une session BGP retourné par l'API.
type SessionResponse struct {
	State            BGPSessionState `json:"state"`
	UpSince          *string         `json:"up_since,omitempty"` // ISO 8601
	PrefixesReceived uint            `json:"prefixes_received"`
	PrefixesSent     uint            `json:"prefixes_sent"`
	LastError        string          `json:"last_error,omitempty"`
}
