package configs

import (
	"os"
)

// FCMConfig holds the configuration for FCM
type FCMConfig struct {
	APIKey string
}

// LoadFCMConfig loads the FCM configuration from environment variables
func LoadFCMConfig() FCMConfig {
	return FCMConfig{
		APIKey: os.Getenv("FCM_API_KEY"),
	}
}

// DBConfig holds the configuration for the database
type DBConfig struct {
	User     string
	Password string
	DBName   string
	Host     string
	Port     string
}

// LoadDBConfig loads the database configuration from environment variables
func LoadDBConfig() DBConfig {
	return DBConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
	}
}
