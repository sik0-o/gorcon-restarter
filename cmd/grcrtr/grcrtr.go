package main

import (
	"fmt"
	"log"

	be "github.com/sik0-o/gorcon-restarter/v2/internal/battleye"
	"github.com/sik0-o/gorcon-restarter/v2/internal/config"
	ll "github.com/sik0-o/gorcon-restarter/v2/internal/logger"
	"github.com/sik0-o/gorcon-restarter/v2/internal/restarter"
	"go.uber.org/zap"
)

//nolint:gochecknoglobals
var (
	version   = "unknown"
	buildTime = "unknown"
)

func main() {
	fmt.Println("Golang RCON server restarter")

	// Env setup
	logger, err := ll.New(version, "production", "Debug")
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Sync() //nolint:errcheck

	logger.Info("Initialize Golang RCON server restarter")

	// Config loading
	conf, err := config.Read()
	if err != nil {
		logger.Fatal("config error", zap.Error(err))
	}

	done := make(chan bool)

	for serverName, serverConf := range conf.Servers {
		logger.Info("handle server configuration", zap.String("serverName", serverName))
		go func(serverName string, serverConf config.ServerConf, logger *zap.Logger) {

			beService, err := be.NewService(logger, serverName, serverConf)
			if err != nil {
				logger.Error("beService error", zap.Error(err))
				return
			}

			restartService, err := restarter.NewService(beService)
			if err != nil {
				logger.Error("restarter service error", zap.Error(err))
				return
			}
			restartService.Restart()

		}(serverName, serverConf, logger)
	}

	// Await proc cancel
	for s := range done {
		if s {
			return
		}
	}

}
