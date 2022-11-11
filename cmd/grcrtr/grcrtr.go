package main

import (
	"fmt"
	"log"

	be "github.com/sik0-o/gorcon-restarter/v2/internal/battleye"
	"github.com/sik0-o/gorcon-restarter/v2/internal/car"
	"github.com/sik0-o/gorcon-restarter/v2/internal/config"
	ll "github.com/sik0-o/gorcon-restarter/v2/internal/logger"
	"github.com/sik0-o/gorcon-restarter/v2/internal/restarter"
	"github.com/sik0-o/gorcon-restarter/v2/internal/server/grpc"
	"github.com/sik0-o/gorcon-restarter/v2/internal/service"
	"go.uber.org/zap"
)

//nolint:gochecknoglobals
var (
	version   = "unknown"
	buildTime = "unknown"
	appConf   *config.Config
	addr      = "127.0.0.1:8080"
)

func main() {
	// t1()

	fmt.Println("Golang RCON server restarter")

	// Env setup
	logger, err := ll.New(version, "local", "Debug")
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
	appConf = conf

	grpc.RunnServe(logger, conf, addr)

	// beService, err := be.NewService(logger, "server1", conf.Servers["server1"])
	// if err != nil {
	// 	logger.Error("beService error", zap.Error(err))
	// 	return
	// }

	// beService.Shutdown()

	done := make(chan bool)

	for serverName, serverConf := range conf.Servers {
		logger.Info("handle server configuration", zap.String("serverName", serverName))
		go func(serverName string, serverConf config.ServerConf, logger *zap.Logger) {

			beService, err := be.NewService(logger, serverName, serverConf)
			if err != nil {
				logger.Error("beService error", zap.Error(err))
				return
			}

			restartService, err := restarter.NewService(logger, beService, service.New(beService))
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

func t1() {
	fmt.Println(car.CreateCode("MyCar", "EnitityAI", []string{"Shit"}))
}
