package config

import (
	"strconv"
	"strings"
)

// Config - Blueprint of app's config
type Config struct {
	AppHost                 string
	AppPort                 string
	DbHost                  string
	DbPort                  string
	DbUser                  string
	DbPassword              string
	DbName                  string
	DbSslmode               string
	RedisHost               string
	RedisPort               string
	JwtSecret               string
	WriteTimeout            int64
	ReadTimeout             int64
	IdleTimeout             int64
	GracefulShutdownTimeout int64
}

// NewConfig - Creates a new config of the server
func NewConfig(env string) (*Config, error) {
	config := Config{}

	for _, v := range strings.Split(env, "\n") {
		c := strings.Split(v, "=")

		switch c[0] {
		case "APP_HOST":
			config.AppHost = c[1]
		case "APP_PORT":
			config.AppPort = c[1]
		case "DB_HOST":
			config.DbHost = c[1]
		case "DB_PORT":
			config.DbPort = c[1]
		case "DB_USER":
			config.DbUser = c[1]
		case "DB_PASSWORD":
			config.DbPassword = c[1]
		case "DB_NAME":
			config.DbName = c[1]
		case "DB_SSLMODE":
			config.DbSslmode = c[1]
		case "REDIS_HOST":
			config.RedisHost = c[1]
		case "REDIS_PORT":
			config.RedisPort = c[1]
		case "JWT_SECRET":
			config.JwtSecret = c[1]
		case "WRITE_TIMEOUT":
			n, err := strconv.ParseInt(c[1], 10, 8)
			if err != nil {
				return &config, err
			}
			config.WriteTimeout = n
		case "READ_TIMEOUT":
			n, err := strconv.ParseInt(c[1], 10, 8)
			if err != nil {
				return &config, err
			}
			config.ReadTimeout = n
		case "IDLE_TIMEOUT":
			n, err := strconv.ParseInt(c[1], 10, 8)
			if err != nil {
				return &config, err
			}
			config.IdleTimeout = n
		case "GRACEFUL_SHUTDOWN_TIMEOUT":
			n, err := strconv.ParseInt(c[1], 10, 8)
			if err != nil {
				return &config, err
			}
			config.GracefulShutdownTimeout = n
		}
	}

	return &config, nil
}
