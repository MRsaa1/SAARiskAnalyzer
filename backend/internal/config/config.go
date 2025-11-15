package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	AppEnv   string
	AppName  string
	Port     string
	DB       DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
	Admin    AdminConfig
	Perf     PerfConfig
	Log      LogConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret             string
	ExpiryMinutes      int
	RefreshExpiryHours int
}

type CORSConfig struct {
	AllowOrigins string
}

type AdminConfig struct {
	Email    string
	Password string
}

type PerfConfig struct {
	MaxSimulations int
	WorkerPoolSize int
}

type LogConfig struct {
	Level  string
	Format string
}

func Load() (*Config, error) {
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	
	// KEY FIX: explicitly bind all environment variables
	viper.BindEnv("DB_NAME")
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_SSLMODE")
	viper.BindEnv("JWT_SECRET")
	viper.BindEnv("ADMIN_EMAIL")
	viper.BindEnv("ADMIN_PASSWORD")
	
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_NAME", "SAA Risk Analyzer")
	viper.SetDefault("PORT", "8083")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_NAME", "risk_db")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("JWT_SECRET", "changeme")
	viper.SetDefault("JWT_EXPIRY_MINUTES", 15)
	viper.SetDefault("JWT_REFRESH_EXPIRY_HOURS", 720)
	viper.SetDefault("VAR_MAX_SIMULATIONS", 100000)
	viper.SetDefault("WORKER_POOL_SIZE", 4)
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "json")
	viper.SetDefault("ADMIN_EMAIL", "admin@example.com")
	viper.SetDefault("ADMIN_PASSWORD", "Admin123456!")

	if err := viper.ReadInConfig(); err != nil {
		// OK if no config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config: %w", err)
		}
	}

	cfg := &Config{
		AppEnv:  viper.GetString("APP_ENV"),
		AppName: viper.GetString("APP_NAME"),
		Port:    viper.GetString("PORT"),
		DB: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		JWT: JWTConfig{
			Secret:             viper.GetString("JWT_SECRET"),
			ExpiryMinutes:      viper.GetInt("JWT_EXPIRY_MINUTES"),
			RefreshExpiryHours: viper.GetInt("JWT_REFRESH_EXPIRY_HOURS"),
		},
		CORS: CORSConfig{
			AllowOrigins: viper.GetString("CORS_ALLOW_ORIGINS"),
		},
		Admin: AdminConfig{
			Email:    viper.GetString("ADMIN_EMAIL"),
			Password: viper.GetString("ADMIN_PASSWORD"),
		},
		Perf: PerfConfig{
			MaxSimulations: viper.GetInt("VAR_MAX_SIMULATIONS"),
			WorkerPoolSize: viper.GetInt("WORKER_POOL_SIZE"),
		},
		Log: LogConfig{
			Level:  viper.GetString("LOG_LEVEL"),
			Format: viper.GetString("LOG_FORMAT"),
		},
	}

	return cfg, nil
}
