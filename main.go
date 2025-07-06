package main

import (
	"log"
	"net/http"
	"shortb/config"

	// "shortb/models"
	"shortb/routes"
)

func main() {
	config.ConnectDB()

	// Nyalakan jika ingin Migrate

	// err := config.DB.AutoMigrate(&models.User{}, &models.Session{}, &models.Clicklogs{}, &models.Shortlink{})
	// if err != nil {
	// 	log.Fatal("❌ Gagal migrasi:", err)
	// } else {
	// 	log.Println("✅ Migrasi sukses!")
	// }

	r := routes.SetupRoutes()

	log.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
