package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/Bibhu20031/SchemaWatch/internal/auth"
	"github.com/Bibhu20031/SchemaWatch/internal/drift"
	"github.com/Bibhu20031/SchemaWatch/internal/observability"
	"github.com/Bibhu20031/SchemaWatch/internal/schema"
	"github.com/Bibhu20031/SchemaWatch/internal/snapshot"
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

	driftRepo := drift.NewRepository(pool)
	driftService := drift.NewService(driftRepo, schemaRepo)

	api.POST("/schemas", schemaHandler.Register)
	api.GET("/schemas", schemaHandler.List)
	api.GET("/schemas/:schema_id/latest", schemaHandler.GetLatest)
	api.GET("/schemas/:schema_id/versions", schemaHandler.ListVersions)

	startScheduler(pool, schemaRepo, driftService)

	log.Println("service started on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func startScheduler(
	pool *pgxpool.Pool,
	schemaRepo *schema.Repository,
	driftService *drift.Service,
) {
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			runSnapshotCycle(pool, schemaRepo, driftService)
		}
	}()
}

func runSnapshotCycle(
	pool *pgxpool.Pool,
	schemaRepo *schema.Repository,
	driftService *drift.Service,
) {
	ctx := context.Background()

	schemas, err := schemaRepo.ListSchemas(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	for _, sc := range schemas {
		schemaID := sc["id"].(int64)
		schemaName := sc["schema_name"].(string)
		tableName := sc["table_name"].(string)

		cols, err := snapshot.FetchTableSchema(ctx, pool, schemaName, tableName)
		if err != nil {
			log.Println(err)
			continue
		}

		data, _ := json.Marshal(cols)
		_, err = pool.Exec(ctx, `
			INSERT INTO schema_versions (schema_id, version, snapshot)
			SELECT $1, COALESCE(MAX(version), 0) + 1, $2
			FROM schema_versions
			WHERE schema_id = $1
		`, schemaID, data)
		if err != nil {
			log.Println(err)
			continue
		}

		fromV, prevRaw, toV, currRaw, err :=
			schemaRepo.GetLastTwoVersions(ctx, schemaID)
		if err != nil || prevRaw == nil {
			continue
		}

		prevCols := decodeSnapshot(prevRaw)
		currCols := decodeSnapshot(currRaw)

		_, err = driftService.Process(
			ctx,
			schemaID,
			fromV,
			toV,
			prevCols,
			currCols,
		)
		if err != nil {
			log.Println(err)
		}
	}
}

func decodeSnapshot(data []byte) []snapshot.Column {
	var cols []snapshot.Column
	_ = json.Unmarshal(data, &cols)
	return cols
}
