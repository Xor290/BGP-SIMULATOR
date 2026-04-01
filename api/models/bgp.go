package models

import "gorm.io/gorm"

// AutonomousSystem représente un AS local géré par le système.
type AutonomousSystem struct {
	gorm.Model
	ASN         uint32 `gorm:"uniqueIndex;not null" json:"asn"`
	Name        string `gorm:"size:128;not null"   json:"name"`
	RouterID    string `gorm:"size:15;not null"    json:"router_id"`
	Description string `gorm:"size:255"            json:"description"`
	Peers       []Peer `gorm:"foreignKey:LocalASID" json:"peers,omitempty"`
}
