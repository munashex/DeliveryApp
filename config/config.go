package config

import (
	"time"
)

type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	Auth        AuthConfig
	Delivery    DeliveryConfig
}

type ServerConfig struct {
	Port            int
	Protocol        string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	CertFile        string
	KeyFile         string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type AuthConfig struct {
	JWTSecret      string
	JWTTTL         time.Duration
	APIKey         string
}

type DeliveryConfig struct {
	DefaultRadius    float64
	MaxRadius        float64
	BaseFee          float64
	FeePerKm         float64
	DispatchTimeout  time.Duration
}

// Load loads configuration from environment variables and files
func Load() (*Config, error) {
	// Implementation would read from env vars, config files, etc.
	// and return a populated Config struct
}

// Sanitized returns a safe version of config for logging (without secrets)
func (c *Config) Sanitized() Config {
	sanitized := *c
	sanitized.Database.Password = "*****"
	sanitized.Auth.JWTSecret = "*****"
	sanitized.Auth.APIKey = "*****"
	return sanitized
}