package server

import (
	"database/sql"
	"net/http"

	"github.com/anand-kay/linkedin-clone/config"
	"github.com/go-redis/redis/v8"
)

// Server - Blueprint of a server
type Server struct {
	Config     *config.Config
	DB         *sql.DB
	RedisDB    *redis.Client
	HTTPServer *http.Server
}
