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
	JWTSecret            string
	AllowedOrigins       []string
}

var config *Config

func IsDevEnv() bool {
	isDev, err := strconv.ParseBool(os.Getenv("IS_DEV"))
	if err != nil {
		// Defaults to prod if IsDev isn't set
		isDev = false
	}
	return isDev
}

func LoadDotenvConfig() *Config {
	isDev := IsDevEnv()

	allowedOrigins := []string{os.Getenv("CLIENT_ORIGIN")}
	if isDev {
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
		AllowedOrigins:       allowedOrigins,
	}
	return config
}

func GetConfig() *Config {
	return config
}
