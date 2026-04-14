package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"pos-service/internal/handler"
	"pos-service/internal/repository"
	"pos-service/internal/routes"
	"pos-service/internal/service"
	"pos-service/internal/validator"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env for local development
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system environment")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("db unreachable: %v", err)
	}

	log.Println("database connected")

	posValidator := &validator.PosValidator{}
	branchRepo := repository.NewPostgresBranchRepository(db)
	posService := service.NewPosService(db, branchRepo, posValidator)
	posHandler := handler.NewPosHandler(posService, posValidator)

	r := gin.Default()
	routes.SetupRoutes(r, posHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server running on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}