package lib

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
}

var config *Config

func LoadConfig() *Config {
	IsDev, err := strconv.ParseBool(os.Getenv("IS_DEV"))
	if err != nil {
		IsDev = true
	}

	config = &Config{
		NATSUrl:              os.Getenv("NATS_URL"),
		NATSUsername:         os.Getenv("NATS_USERNAME"),
		NATSPassword:         os.Getenv("NATS_PASSWORD"),
		AppDBPostgresURL:     os.Getenv("APP_DB_POSTGRES_URL"),
		MessageDBPostgresURL: os.Getenv("MESSAGE_DB_POSTGRES_URL"),
		SnowflakeNodeId:      os.Getenv("SNOWFLAKE_NODE_ID"),
		Port:                 os.Getenv("PORT"),
		IsDev:                IsDev,
		JWTSecret:            os.Getenv("JWT_SECRET"),
	}
	return config
}

func GetConfig() *Config {
	if config == nil {
		return LoadConfig()
	}
	return config
}
