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
	"github.com/Vighnesh-V-H/sync/internal/handler"
	"github.com/Vighnesh-V-H/sync/internal/logger"
	"github.com/Vighnesh-V-H/sync/internal/repositories"
	"github.com/Vighnesh-V-H/sync/internal/routes"
	"github.com/Vighnesh-V-H/sync/internal/service"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
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
	log.Info().Msg("Starting events service")

	database, err := db.NewDB(cfg.Database.URL, log)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	defer database.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.URL,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		DialTimeout:  time.Duration(cfg.Redis.Timeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Redis.Timeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Redis.Timeout) * time.Second,
	})
	defer redisClient.Close()

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()
	if err := redisClient.Ping(pingCtx).Err(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	log.Info().Msg("Successfully connected to Redis")

	
	eventRepo := repositories.NewEventRepository(database, redisClient, log)
	eventSvc := service.NewEventService(eventRepo, log)
	eventHandler := handler.NewEventHandler(eventSvc, log)

	if cfg.Primary.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	api := router.Group("/api/v1")

	routes.SetupEventRoutes(api, eventHandler)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.EventsPort)
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
