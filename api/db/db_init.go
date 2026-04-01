package db

import (
	"fmt"
	"log"
	"os"

	"bgp-manager/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database est le type central du package db. Toutes les opérations sont
// définies comme méthodes sur ce type pour éviter de passer *gorm.DB partout.
type Database struct {
	*gorm.DB
}

// Connect ouvre la connexion PostgreSQL et retourne un *Database.
func Connect() *Database {
	dsn := buildDSN()

	logLevel := logger.Warn
	if os.Getenv("ENV") == "development" {
		logLevel = logger.Info
	}

	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("❌ Connexion DB échouée : %v", err)
	}

	log.Println("✅ Connecté à PostgreSQL")
	return &Database{gdb}
}

// Migrate applique les migrations automatiques pour tous les modèles.
func (d *Database) Migrate() {
	err := d.AutoMigrate(
		&models.AutonomousSystem{},
		&models.Peer{},
		&models.BGPSession{},
		&models.BGPRoute{},
	)
	if err != nil {
		log.Fatalf("❌ Migration échouée : %v", err)
	}
	log.Println("✅ Migrations appliquées")
}


// buildDSN construit la DSN depuis les variables d'environnement.
func buildDSN() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "bgp")
	password := getEnv("DB_PASSWORD", "bgp")
	dbname := getEnv("DB_NAME", "bgp_manager")
	sslmode := getEnv("DB_SSLMODE", "disable")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		host, port, user, password, dbname, sslmode,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
