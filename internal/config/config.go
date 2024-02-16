package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

// Rest is a configuration for REST API.
type Rest struct {
	Port               int      `env:"PORT" env-default:"8080"`
	Host               string   `env:"HOST" env-default:"0.0.0.0"`
	AllowedCORSOrigins []string `env:"ALLOWED_CORS_ORIGINS" envSeparator:"," env-default:"*"`
}

// Project - project configuration.
type Project struct {
	Debug       bool   `env:"DEBUG" env-default:"false"`
	Name        string `env:"PROJECT_NAME" env-default:"Monolith"`
	Environment string `env:"ENVIRONMENT" env-default:"development"`
	ServiceName string `env:"SERVICE_NAME" env-default:"monolith"`
}

type Auth struct {
	JWTKey string `env:"JWT_KEY" env-default:"secret"`
}

// Database is a configuration for database.
type Database struct {
	DSN string `env:"DATABASE_DSN" env-default:"postgres://root:root@localhost:5433/har?sslmode=disable"`
}

// Logger is a configuration for logger.
type Logger struct {
	LogLevel string `env:"LOG_LEVEL" env-default:"info"`
	DevMode  bool   `env:"DEV_MODE" env-default:"false"`
	Encoder  string `env:"ENCODER" env-default:"console"`
}

type Swagger struct {
	Host    string   `env:"SWAGGER_HOST" env-default:"localhost:8080"`
	Schemes []string `env:"SWAGGER_SCHEMES" envSeparator:"," env-default:"http"`
	URL     string   `env:"SWAGGER_URL" env-default:"/v1/api/swagger"`
}

// Config - contains all configuration parameters in config package.
type Config struct {
	Project  Project
	Rest     Rest
	Database Database
	Logger   Logger
	Swagger  Swagger
	Auth     Auth
}

func Get() (Config, error) {
	config := Config{}

	if err := cleanenv.ReadEnv(&config.Rest); err != nil {
		return config, fmt.Errorf("error reading REST config: %w", err)
	}

	if err := cleanenv.ReadEnv(&config.Project); err != nil {
		return config, fmt.Errorf("error reading project config: %w", err)
	}

	if err := cleanenv.ReadEnv(&config.Database); err != nil {
		return config, fmt.Errorf("error reading database config: %w", err)
	}

	if err := cleanenv.ReadEnv(&config.Auth); err != nil {
		return config, fmt.Errorf("error reading auth config: %w", err)
	}

	if err := cleanenv.ReadEnv(&config.Logger); err != nil {
		return config, fmt.Errorf("error reading logger config: %w", err)
	}

	if err := cleanenv.ReadEnv(&config.Swagger); err != nil {
		return config, fmt.Errorf("error reading swagger config: %w", err)
	}

	return config, nil
}
