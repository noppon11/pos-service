package main

import (
	"database/sql"
	"log"
	"os"

	"pos-service/internal/handler"
	"pos-service/internal/repository"
	"pos-service/internal/routes"
	"pos-service/internal/service"
	"pos-service/internal/validator"


	_ "github.com/lib/pq"
	"github.com/gin-gonic/gin"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("db unreachable: %v", err)
	}

	posValidator := &validator.PosValidator{}
	branchRepo := repository.NewInMemoryBranchRepository()
	posService := service.NewPosService(db, branchRepo)
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