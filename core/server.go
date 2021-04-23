package core

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anand-kay/linkedin-clone/config"
	"github.com/anand-kay/linkedin-clone/server"
)

// NewServer - Creates a new instance of the server
func NewServer(config *config.Config) *server.Server {
	server := &server.Server{}

	server.Config = config

	// Create database instance
	db, err := newDB(config)
	if err != nil {
		log.Fatalln(err)
	}
	server.DB = db

	server.RedisDB = newRedisDB(config)

	server.HTTPServer = &http.Server{}
	server.HTTPServer.Addr = config.AppHost + ":" + config.AppPort
	server.HTTPServer.Handler = GetHandler(NewRouter(server))
	server.HTTPServer.WriteTimeout = time.Second * time.Duration(rand.Int63n(config.WriteTimeout))
	server.HTTPServer.ReadTimeout = time.Second * time.Duration(rand.Int63n(config.ReadTimeout))
	server.HTTPServer.IdleTimeout = time.Second * time.Duration(rand.Int63n(config.IdleTimeout))

	return server
}

// StartServer - Starts the server
func StartServer(server *http.Server) error {
	return server.ListenAndServe()
}

// ShutdownServer - Shuts down the server gracefully or otherwise depending on the OS signal received
func ShutdownServer(server *server.Server) {
	// Create channel to receive OS signal
	c := make(chan os.Signal, 1)

	// SIGINT (Ctrl+C) - Graceful shutdown
	// SIGQUIT (Ctrl+\) - Quit process without graceful shutdown
	// SIGKILL or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT)

	// Block until we receive our signal.
	sig := <-c

	if sig == syscall.SIGINT {
		// Create a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(rand.Int63n(server.Config.GracefulShutdownTimeout)))
		defer cancel()
		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.

		// Close database connection
		err := server.DB.Close()
		if err != nil {
			log.Println("Error while trying to close database connection. Quiting process.")
			log.Println(err)
			os.Exit(1)
		}

		// Close Redis client
		err = server.RedisDB.Close()
		if err != nil {
			log.Println("Error while trying to close Redis client. Quiting process.")
			log.Println(err)
			os.Exit(1)
		}

		// Shutdown server gracefully
		err = server.HTTPServer.Shutdown(ctx)
		if err != nil {
			log.Println("Error while trying to shutdown server gracefully. Quiting process.")
			log.Println(err)
			os.Exit(1)
		}
		// Optionally, you could run server.httpServer.Shutdown in a goroutine and block on
		// <-ctx.Done() if your application should wait for other services
		// to finalize based on context cancellation.

		fmt.Println()
		log.Println("Shutting down server gracefully")
		os.Exit(0)
	} else if sig == syscall.SIGQUIT {
		log.Println("Quiting process. No graceful shutdown.")
		os.Exit(1)
	}
}
