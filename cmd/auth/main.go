package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Vighnesh-V-H/sync/internal/config"
	"github.com/Vighnesh-V-H/sync/internal/db"
	"github.com/Vighnesh-V-H/sync/internal/logger"
	"github.com/Vighnesh-V-H/sync/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	
	logCfg := logger.Config{
		Level:       cfg.Logging.Level,
		Format:      "json",
		ServiceName: cfg.Observability.ServiceName,
		Environment: cfg.Primary.Env,
		IsProd:      cfg.Primary.Env == "prod",
	}
	if cfg.Logging.Pretty {
		logCfg.Format = "console"
	}
	log := logger.New(logCfg)
	log.Info().Msg("Starting auth service")

	
	database, err := db.NewDB(cfg.Database.URL, log)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	defer database.Close()


	if cfg.Primary.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	routes.SetupAuthRoutes(router, log)

	
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Info().Str("address", addr).Msg("Starting HTTP server")

	go func() {
		if err := router.Run(addr); err != nil {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = ctx

	log.Info().Msg("Server exited")
}
