package main

import (
	"log"
	"net/http"
	"os"

	"bgp-manager/db"
	"bgp-manager/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Chargement des variables d'environnement
	if err := godotenv.Load(); err != nil {
		log.Println("Aucun fichier .env trouvé, utilisation des valeurs par défaut.")
	}

	// Initialisation de la base de données
	database := db.Connect()
	database.Migrate()
	log.Printf("✅ Database initialisée")

	// Créer le routeur Gin
	r := gin.Default()

	// Configuration des sessions
	store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	store.Options(sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	r.Use(sessions.Sessions("mysession", store))

	// Configuration CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:8080", "http://192.168.1.94", "https://limon-vitrine.sbs"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Fichiers statiques
	r.Static("/uploads", "./uploads")

	// Enregistrer toutes les routes
	routes.SetupRoutes(r, database)

	// Récupérer le port depuis les variables d'environnement
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("✅ Serveur démarré sur le port %s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Erreur au lancement du serveur : %v", err)
	}
}
