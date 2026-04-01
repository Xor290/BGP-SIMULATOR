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

func (PrefixSinceAS) TableName() string { return "prefix_since_as" }

// PrefixSinceAS représente un préfixe annoncé depuis un AS local vers ses voisins BGP.
type PrefixSinceAS struct {
	gorm.Model
	Prefix    string           `gorm:"size:49;not null" json:"prefix"`    // ex: "192.168.1.0/24"
	ASID      uint             `gorm:"not null;index"   json:"as_id"`
	AS        AutonomousSystem `gorm:"foreignKey:ASID"  json:"as,omitempty"`
	NextHop   string           `gorm:"size:45;not null" json:"next_hop"` // ex: "10.0.0.11"
	LocalPref uint             `gorm:"default:100"      json:"local_pref"`
	MED       uint             `gorm:"default:0"        json:"med"`
	Active    bool             `gorm:"default:true"     json:"active"`   // false = préfixe retiré (withdraw)
}

type CreatePrefixSinceAS struct {
	Prefix    string `json:"prefix"` // ex: "192.168.1.0/24"
	ASN       uint32 `json:"asn"`
	NextHop   string `json:"next_hop"` // ex: "10.0.0.11"
	LocalPref uint   `json:"local_pref"`
}
