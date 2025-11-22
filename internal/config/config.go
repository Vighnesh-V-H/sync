package config

import (
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	Primary       Primary              `koanf:"primary" validate:"required"`
	Server        ServerConfig         `koanf:"server" validate:"required"`
	Redis         RedisConfig          `koanf:"redis" validate:"required"`
	ClickHouse    ClickHouseConfig     `koanf:"clickhouse" validate:"required"`
	Database      DatabaseConfig       `koanf:"database" validate:"required"`
	Logging       LoggingConfig        `koanf:"logging" validate:"required"`
	App           AppConfig            `koanf:"app" validate:"required"`
	JWT           JWTConfig            `koanf:"jwt" validate:"required"`
	Observability *ObservabilityConfig `koanf:"observability"`
}

type Primary struct {
	Env string `koanf:"env" validate:"required,oneof=dev staging prod"`
}

type JWTConfig struct {
	Secret string `koanf:"secret" validate:"required"`
}

type ServerConfig struct {
	Host               string   `koanf:"host" validate:"required"`
	Port               int      `koanf:"port" validate:"required,min=1,max=65535"`
	AuthPort           int      `koanf:"auth_port" validate:"required,min=1,max=65535"`
	EventsPort         int      `koanf:"events_port" validate:"required,min=1,max=65535"`
	ReadTimeout        int      `koanf:"read_timeout" validate:"required,min=1"`
	WriteTimeout       int      `koanf:"write_timeout" validate:"required,min=1"`
	IdleTimeout        int      `koanf:"idle_timeout" validate:"required,min=1"`
	CORSAllowedOrigins []string `koanf:"cors_allowed_origins" validate:"required,dive,http_url"`
}

type RedisConfig struct {
	URL      string `koanf:"url" validate:"required,uri"` 
	// Old fields (commented out):
	// URL      string `koanf:"url" validate:"required,uri"`
	// Password string `koanf:"password"`
	// DB       int    `koanf:"db" validate:"omitempty,min=0,max=15"`
	Timeout  int    `koanf:"timeout" validate:"required,min=1"`
}

type ClickHouseConfig struct {
	DSN             string `koanf:"dsn" validate:"required"`
	Timeout         int    `koanf:"timeout" validate:"required,min=1"`
	MaxOpen         int    `koanf:"max_open" validate:"required,min=1,max=1000"`
	MaxIdle         int    `koanf:"max_idle" validate:"omitempty,min=0"`
	ConnMaxLifetime int    `koanf:"conn_max_lifetime" validate:"omitempty,min=1"`
}

type DatabaseConfig struct {
	URL string `koanf:"url" validate:"required"`
}

type LoggingConfig struct {
	Level  string `koanf:"level" validate:"required,oneof=debug info warn error fatal panic"`
	Pretty bool   `koanf:"pretty"`
}

type AppConfig struct {
	TenantDefault string `koanf:"tenant_default" validate:"required"`
	WindowSecs    int    `koanf:"window_secs" validate:"required,min=1,max=3600"`
	BatchSize     int    `koanf:"batch_size" validate:"required,min=10,max=10000"`
}

type ObservabilityConfig struct {
	ServiceName    string `koanf:"service_name" validate:"required"`
	Environment    string `koanf:"environment" validate:"required,oneof=dev staging prod"`
	PrometheusPort int    `koanf:"prometheus_port" validate:"omitempty,min=9090,max=65535"`
}

func DefaultObservabilityConfig() *ObservabilityConfig {
	return &ObservabilityConfig{
		ServiceName:    "analytics-engine",
		Environment:    "dev",
		PrometheusPort: 9090,
	}
}

func (o *ObservabilityConfig) Validate() error {
	v := validator.New()
	return v.Struct(o)
}

func LoadConfig() (*Config, error) {
	tempLogger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	k := koanf.New(".")
	err := k.Load(env.Provider("SYNC_", ".", func(s string) string {
		// SYNC_PRIMARY_ENV -> primary.env
		return strings.ToLower(strings.Replace(strings.TrimPrefix(s, "SYNC_"), "_", ".", 1))
	}), nil)
	if err != nil {
		tempLogger.Fatal().Err(err).Msg("could not load initial env variables")
	}

	mainConfig := &Config{}
	err = k.Unmarshal("", mainConfig)
	if err != nil {
		tempLogger.Fatal().Err(err).Msg("could not unmarshal main config")
	}

	if mainConfig.Observability == nil {
		mainConfig.Observability = DefaultObservabilityConfig()
	}

	mainConfig.Observability.ServiceName = "analytics-engine"
	mainConfig.Observability.Environment = mainConfig.Primary.Env

	if err := mainConfig.Observability.Validate(); err != nil {
		tempLogger.Fatal().Err(err).Msg("invalid observability config")
	}

	if mainConfig.Server.Host == "" {
		mainConfig.Server.Host = "0.0.0.0"
	}
	if mainConfig.Server.Port == 0 {
		mainConfig.Server.Port = 8080
	}
	if mainConfig.Server.AuthPort == 0 {
		mainConfig.Server.AuthPort = 8081
	}
	if mainConfig.Server.EventsPort == 0 {
		mainConfig.Server.EventsPort = 8083
	}
	if mainConfig.Server.ReadTimeout == 0 {
		mainConfig.Server.ReadTimeout = 10
	}
	if mainConfig.Server.WriteTimeout == 0 {
		mainConfig.Server.WriteTimeout = 10
	}
	if mainConfig.Server.IdleTimeout == 0 {
		mainConfig.Server.IdleTimeout = 300
	}
	if mainConfig.Server.CORSAllowedOrigins == nil {
		mainConfig.Server.CORSAllowedOrigins = []string{"*"}
	}
	// Redis DB default (commented out for Upstash URL mode)
	// if mainConfig.Redis.DB == 0 {
	// 	mainConfig.Redis.DB = 0
	// }
	if mainConfig.Redis.Timeout == 0 {
		mainConfig.Redis.Timeout = 5
	}
	if mainConfig.ClickHouse.MaxIdle == 0 {
		mainConfig.ClickHouse.MaxIdle = 10
	}
	if mainConfig.ClickHouse.ConnMaxLifetime == 0 {
		mainConfig.ClickHouse.ConnMaxLifetime = 3600
	}
	if mainConfig.Logging.Level == "" {
		mainConfig.Logging.Level = "info"
	}
	if mainConfig.App.TenantDefault == "" {
		mainConfig.App.TenantDefault = "default"
	}
	if mainConfig.App.WindowSecs == 0 {
		mainConfig.App.WindowSecs = 60
	}
	if mainConfig.App.BatchSize == 0 {
		mainConfig.App.BatchSize = 100
	}

	return mainConfig, nil
}
