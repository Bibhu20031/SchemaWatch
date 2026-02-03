package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/Bibhu20031/SchemaWatch/internal/auth"
	"github.com/Bibhu20031/SchemaWatch/internal/observability"
	"github.com/Bibhu20031/SchemaWatch/internal/schema"
	database "github.com/Bibhu20031/SchemaWatch/internal/storage"
)

func main() {
	_ = godotenv.Load()

	pool, err := database.ConnectDB()
	if err != nil {
		log.Fatal("failed to connect database")
	}
	defer pool.Close()

	r := gin.Default()
	_ = r.SetTrustedProxies(nil)

	healthHandler := observability.NewHealthHandler(pool)
	r.GET("/health", healthHandler.Check)

	authRepo := auth.NewRepository(pool)
	authService := auth.NewService(authRepo)
	authMiddleware := auth.NewMiddleware(authService)

	api := r.Group("/v1")
	api.Use(authMiddleware.RequireAPIKey())

	schemaRepo := schema.NewRepository(pool)
	schemaService := schema.NewService(schemaRepo, pool)
	schemaHandler := schema.NewHandler(schemaService)

	api.POST("/schemas", schemaHandler.Register)
	api.GET("/schemas", schemaHandler.List)
	api.GET("/schemas/:schema_id/latest", schemaHandler.GetLatest)

	log.Println("service started on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
