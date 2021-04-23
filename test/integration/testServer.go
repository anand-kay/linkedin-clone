package integration

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis"
	"github.com/anand-kay/linkedin-clone/core"
	"github.com/anand-kay/linkedin-clone/server"
	"github.com/go-redis/redis/v8"
)

var srv *server.Server
var s *miniredis.Miniredis

func setupServer() sqlmock.Sqlmock {
	srv = &server.Server{}

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalln(err)
	}
	srv.DB = db

	s, err = miniredis.Run()
	if err != nil {
		log.Fatalln(err)
	}
	// defer s.Close()
	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	srv.RedisDB = rdb

	srv.HTTPServer = &http.Server{}
	srv.HTTPServer.Addr = "localhost:3000"
	srv.HTTPServer.Handler = core.GetHandler(core.NewRouter(srv))

	go func() {
		if err := core.StartServer(srv.HTTPServer); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatalln(err)
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)

	return mock
}

func teardownServer() {
	srv.DB.Close()
	s.Close()
	srv.RedisDB.Close()
	srv.HTTPServer.Shutdown(context.Background())
}
