package main

import (
	"database/sql"
	"log"
	"os"

	"pos-service/internal/handler"
	"pos-service/internal/routes"
	"pos-service/internal/service"
	"pos-service/internal/validator"

	_ "github.com/lib/pq" // 👈 ตัวนี้สำคัญมาก
	"github.com/gin-gonic/gin"
)

func main() {
	// =========================
	// 1. Load config (simple version)
	// =========================
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// =========================
	// 2. Init DB
	// =========================
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	// optional: ping DB on startup
	if err := db.Ping(); err != nil {
		log.Fatalf("db unreachable: %v", err)
	}

	// =========================
	// 3. Init dependencies
	// =========================
	posValidator := &validator.PosValidator{}
	posService := service.NewPosService(db, posValidator)
	posHandler := handler.NewPosHandler(posService, posValidator)

	// =========================
	// 4. Init router
	// =========================
	r := gin.Default()

	// =========================
	// 5. Register routes
	// =========================
	routes.SetupRoutes(r, posHandler)

	// =========================
	// 6. Run server
	// =========================
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server running on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}