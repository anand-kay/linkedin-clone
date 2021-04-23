package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/anand-kay/linkedin-clone/config"
	"github.com/anand-kay/linkedin-clone/core"
	"github.com/anand-kay/linkedin-clone/libs"
	"github.com/anand-kay/linkedin-clone/utils"
)

func main() {
	// Create a new config for the server from .env file
	envStr, err := utils.ReadEnvFile(".env")
	if err != nil {
		log.Fatalln(err)
	}
	config, err := config.NewConfig(envStr)
	if err != nil {
		log.Fatalln(err)
	}

	// Create server
	server := core.NewServer(config)

	// Start http server
	go func() {
		if err := core.StartServer(server.HTTPServer); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatalln(err)
			}
		}
	}()

	// Populate graph in redis db
	err = libs.PopulateGraph(context.Background(), server.DB, server.RedisDB)
	if err != nil {
		log.Fatalln(err)
	}

	// Shutdown server
	core.ShutdownServer(server)
}
