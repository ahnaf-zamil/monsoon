package lib

import "os"

type Config struct {
	NATSUrl              string
	NATSUsername         string
	NATSPassword         string
	AppDBPostgresURL     string
	MessageDBPostgresURL string
	SnowflakeNodeId      string
}

var config *Config

func LoadConfig() *Config {
	config = &Config{
		NATSUrl:              os.Getenv("NATS_URL"),
		NATSUsername:         os.Getenv("NATS_USERNAME"),
		NATSPassword:         os.Getenv("NATS_PASSWORD"),
		AppDBPostgresURL:     os.Getenv("APP_DB_POSTGRES_URL"),
		MessageDBPostgresURL: os.Getenv("MESSAGE_DB_POSTGRES_URL"),
		SnowflakeNodeId:      os.Getenv("SNOWFLAKE_NODE_ID"),
	}
	return config
}

func GetConfig() *Config {
	if config == nil {
		return LoadConfig()
	}
	return config
}
