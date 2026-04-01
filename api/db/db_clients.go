package db

import (
	"bgp-manager/models"

	"golang.org/x/crypto/bcrypt"
)

func (db *Database) RegisterClient(req *models.RegisterClientRequest) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	client := models.Client{
		Username: req.Username,
		Password: string(hashed),
	}

	return db.Create(&client).Error
}

func (db *Database) ConnectClient(req *models.LoginRequest) (*models.Client, error) {
	var client models.Client
	if err := db.Where("username = ?", req.Username).First(&client).Error; err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(client.Password), []byte(req.Password)); err != nil {
		return nil, err
	}

	return &client, nil
}
