package app

import (
	"github.com/Bloxico/exchange-gateway/sofija/config"
	"github.com/Bloxico/exchange-gateway/sofija/database"
	"github.com/Bloxico/exchange-gateway/sofija/log"
)

type App struct {
	Config config.Config
	Logger log.Logger

	DB *database.DB
	// AMQP     *amqp.Client
}
