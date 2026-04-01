package db

import (
	"bgp-manager/models"
)

func (db *Database) CreateAS(as models.AutonomousSystem) error {
	return db.DB.Create(&as).Error
}

func (db *Database) GetAS(asn uint) (models.AutonomousSystem, error) {
	var as models.AutonomousSystem
	err := db.DB.First(&as, asn).Error
	return as, err
}
