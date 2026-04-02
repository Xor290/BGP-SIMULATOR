package db

import (
	"bgp-manager/models"
)

func (db *Database) CreateAS(as models.AutonomousSystem) error {
	return db.DB.Create(&as).Error
}

func (db *Database) GetAS(id uint) (models.AutonomousSystem, error) {
	var as models.AutonomousSystem
	err := db.DB.First(&as, id).Error
	return as, err
}

func (db *Database) GetASByASN(asn uint32) (models.AutonomousSystem, error) {
	var as models.AutonomousSystem
	err := db.DB.Where("asn = ?", asn).First(&as).Error
	return as, err
}

func (db *Database) GetPrefixesByASID(asID uint) ([]models.PrefixSinceAS, error) {
	var prefixes []models.PrefixSinceAS
	return prefixes, db.DB.Where("as_id = ? AND active = ?", asID, true).Find(&prefixes).Error
}
