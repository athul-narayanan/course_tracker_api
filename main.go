package main

import (
	"context"
	"course-tracker/config"
	"course-tracker/internal/auth"
	"course-tracker/internal/kafka"
	"course-tracker/internal/notification"
	"course-tracker/internal/subscription"
	"course-tracker/internal/university"
	"log"
	"os"

	"course-tracker/internal/util"

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

	producer := kafka.NewProducer()

	universityService := &university.UniversityService{DB: db, CFG: &cfg, Producer: producer}

	university.RegisterRoutes(r, universityService)

	subService := subscription.SubscriptionService{DB: db, CFG: &cfg}
	subscription.RegisterRoutes(r, &subService)

	ctx := context.Background()
	emailService := util.NewEmailService()
	consumer := kafka.NewConsumer(universityService, emailService)

	notificationService := notification.NotificationService{DB: db, Config: &cfg}
	notification.RegisterRoutes(r, &notificationService)
	notification_consumer := kafka.NewNotificationConsumer(universityService, &notificationService)

	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Printf("Kafka consumer stopped: %v", err)
		}
	}()

	go func() {
		if err := notification_consumer.Start(ctx); err != nil {
			log.Printf("Kafka consumer stopped: %v", err)
		}
	}()

	// starting the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on 0.0.0.0:%s ...", port)
	log.Fatal(r.Run("0.0.0.0:" + port))

}
