// Main egw entry point: initializes and starts the server
package main

import (
	"fmt"
	"github.com/Bloxico/exchange-gateway/sofija/app"
	"github.com/Bloxico/exchange-gateway/sofija/config"
	"github.com/Bloxico/exchange-gateway/sofija/server"
	"github.com/pkg/errors"
)

func runServer() error {
	egwApp := app.MustInitializeApp()

	cfg := config.ServerConfig{
		Port:   egwApp.Config.Http.Port,
		Logger: egwApp.Logger,
	}

	srv := server.NewServer(cfg, egwApp.DB)

	egwApp.Logger.Infof("Server started at %d", cfg.Port)
	err := srv.ListenAndServe("local", "domain")
	if err != nil {
		return errors.Wrap(err, "listen and serve")
	}
	return nil
}

func main() {
	err := runServer()
	if err != nil {
		fmt.Printf("failed starting app: %s", err)
		panic(err)
	}
}
