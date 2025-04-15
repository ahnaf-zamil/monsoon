package main

import "os"

type Config struct {
	NATSUrl      string
	NATSUsername string
	NATSPassword string
	PostgresURL  string
}

func LoadConfig() Config {
	return Config{
		NATSUrl:      os.Getenv("NATS_URL"),
		NATSUsername: os.Getenv("NATS_USERNAME"),
		NATSPassword: os.Getenv("NATS_PASSWORD"),
		PostgresURL:  os.Getenv("POSTGRES_URL"),
	}
}
