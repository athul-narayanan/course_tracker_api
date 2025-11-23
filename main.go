package main

import (
	"course-tracker/config"
	"course-tracker/internal/auth"
	"course-tracker/internal/subscription"
	"course-tracker/internal/university"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()

	dsn := "host=" + cfg.DBHost +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName +
		" port=" + cfg.DBPort +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	userService := &auth.AuthService{DB: db, CFG: &cfg}

	auth.RegisterRoutes(r, userService)

	universityService := &university.UniversityService{DB: db, CFG: &cfg}

	university.RegisterRoutes(r, universityService)

	subService := subscription.SubscriptionService{DB: db, CFG: &cfg}
	subscription.RegisterRoutes(r, &subService)

	// starting the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on 0.0.0.0:%s ...", port)
	log.Fatal(r.Run("0.0.0.0:" + port))

}
