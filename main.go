package main

import (
	"github.com/nicitapa/cooking-service/internal/configs"
	"github.com/nicitapa/cooking-service/internal/controller"
	"github.com/nicitapa/cooking-service/internal/repository"
	"github.com/nicitapa/cooking-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"os"

	_ "github.com/nicitapa/cooking-service/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Recipe Catalog API
// @version 1.0
// @description API for managing cooking recipes
// @host localhost:8080
// @BasePath /
func main() {
	// Zerolog config (pretty in dev)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	cfg := configs.Load()
	if cfg.DatabaseURL == "" {
		log.Fatal().Msg("DATABASE_URL is required")
	}

	db, err := sqlx.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("db connect failed")
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal().Err(err).Msg("db ping failed")
	}

	repo := repository.NewRecipeRepository(db)
	svc := service.NewRecipeService(repo)

	router := gin.Default()
	ctrl := controller.NewController(router, svc)
	ctrl.RegisterRoutes()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	log.Info().Msgf("server listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal().Err(err).Msg("server run failed")
	}
}
