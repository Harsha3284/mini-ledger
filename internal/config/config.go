package config

import "os"

type Config struct {
	HTTPPort string
	DBURL    string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://ledger:ledger123@127.0.0.1:5433/mini_ledger?sslmode=disable"
	}

	return Config{
		HTTPPort: port,
		DBURL:    dbURL,
	}
}
