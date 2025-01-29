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
