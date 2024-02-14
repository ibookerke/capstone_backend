package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

// Rest is a configuration for REST API.
type Rest struct {
	Port               int      `env:"PORT" env-default:"8003"`
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

// Database is a configuration for database.
type Database struct {
	DSN string `env:"DATABASE_DSN" env-default:"postgres://postgres:postgres@localhost:5434/postgres?sslmode=disable"`
}

// Logger is a configuration for logger.
type Logger struct {
	LogLevel string `env:"LOG_LEVEL" env-default:"info"`
	DevMode  bool   `env:"DEV_MODE" env-default:"false"`
	Encoder  string `env:"ENCODER" env-default:"console"`
}

type AwsS3 struct {
	Region string `env:"AWS_REGION" env-default:"eu-central-1"`
	Bucket string `env:"AWS_BUCKET_NAME" env-default:"product"`
}

type InvestmentOpportunity struct {
	Permissions []string `env:"INVESTMENT_OPPORTUNITY_PERMISSIONS" envSeparator:"," env-default:""`
	Audience    []string `env:"INVESTMENT_OPPORTUNITY_AUDIENCE" envSeparator:"," env-default:""`
	Issuer      string   `env:"INVESTMENT_OPPORTUNITY_ISSUER" env-default:"*"`
}

type Swagger struct {
	Host    string   `env:"SWAGGER_HOST" env-default:"localhost:8003"`
	Schemes []string `env:"SWAGGER_SCHEMES" envSeparator:"," env-default:"http"`
	URL     string   `env:"SWAGGER_URL" env-default:"/v1/api/swagger"`
}

type Authorization struct {
	InvestmentOpportunity InvestmentOpportunity
}

type Kafka struct {
	Brokers []string `env:"KAFKA_BROKERS" envSeparator:"," env-default:"localhost:9092"`
}

// Config - contains all configuration parameters in config package.
type Config struct {
	Project       Project
	Rest          Rest
	AwsS3         AwsS3
	Database      Database
	Logger        Logger
	Authorization Authorization
	Swagger       Swagger
	Kafka         Kafka
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

	if err := cleanenv.ReadEnv(&config.Logger); err != nil {
		return config, fmt.Errorf("error reading logger config: %w", err)
	}

	if err := cleanenv.ReadEnv(&config.AwsS3); err != nil {
		return config, fmt.Errorf("error reading aws s3 config: %w", err)
	}

	if err := cleanenv.ReadEnv(&config.Authorization.InvestmentOpportunity); err != nil {
		return config, fmt.Errorf("error reading investment opportunity config: %w", err)
	}

	if err := cleanenv.ReadEnv(&config.Swagger); err != nil {
		return config, fmt.Errorf("error reading swagger config: %w", err)
	}

	if err := cleanenv.ReadEnv(&config.Kafka); err != nil {
		return config, fmt.Errorf("error reading swagger config: %w", err)
	}

	return config, nil
}
