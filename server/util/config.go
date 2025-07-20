package util

import (
	"os"
	"strconv"
)

type Config struct {
	NATSUrl              string
	NATSUsername         string
	NATSPassword         string
	AppDBPostgresURL     string
	MessageDBPostgresURL string
	SnowflakeNodeId      string
	Port                 string
	IsDev                bool
	JWTSecret            string
	AllowedOrigins       []string
}

var config *Config

func LoadConfig() *Config {
	IsDev, err := strconv.ParseBool(os.Getenv("IS_DEV"))
	if err != nil {
		// Defaults to prod if IsDev isn't set
		IsDev = false
	}

	allowedOrigins := []string{os.Getenv("CLIENT_ORIGIN")}
	if IsDev {
		allowedOrigins = append(allowedOrigins, "http://localhost:5173") // Frontend dev server port
	}

	config = &Config{
		NATSUrl:              os.Getenv("NATS_URL"),
		NATSUsername:         os.Getenv("NATS_USERNAME"),
		NATSPassword:         os.Getenv("NATS_PASSWORD"),
		AppDBPostgresURL:     os.Getenv("APP_DB_POSTGRES_URL"),
		MessageDBPostgresURL: os.Getenv("MESSAGE_DB_POSTGRES_URL"),
		SnowflakeNodeId:      os.Getenv("SNOWFLAKE_NODE_ID"),
		Port:                 os.Getenv("PORT"),
		JWTSecret:            os.Getenv("JWT_SECRET"),
		IsDev:                IsDev,
		AllowedOrigins:       allowedOrigins,
	}
	return config
}

func GetConfig() *Config {
	if config == nil {
		return LoadConfig()
	}
	return config
}
