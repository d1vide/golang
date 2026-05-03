package config

import "os"

type Config struct {
	Addr     string
	CertFile string
	KeyFile  string
	DSN      string
}

func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func New() Config {
	return Config{
		Addr:     getEnv("APP_ADDR", ":8443"),
		CertFile: getEnv("CERT_FILE", "certs/server.crt"),
		KeyFile:  getEnv("KEY_FILE", "certs/server.key"),
		DSN:      getEnv("APP_DSN", "postgres://postgres:postgres@localhost:5432/study_security?sslmode=disable"),
	}
}