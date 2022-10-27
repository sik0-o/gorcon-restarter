package main

import (
	"fmt"
	"log"

	"github.com/multiplay/go-battleye"
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
	// Я хочу выполнять событие A сервера B через каждые C секунд (часов)

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

			restartService.Restart()

			// beService.Login("sasipisu")
			// beService.Players()
			// beService.Admins()
			// beService.Missions()

			// beService.MonitorEnable()

			beService.Say("Hello friend")
			return
			// client, err := battleye.NewClient(
			// 	serverConf.Host+":"+strconv.Itoa(serverConf.Port),
			// 	serverConf.Password,
			// )
			// if err != nil {
			// 	logger.Error(
			// 		"client initialize error",
			// 		zap.String("serverName", serverName),
			// 		zap.String("serverHost", serverConf.Host),
			// 		zap.Int("serverPort", serverConf.Port),
			// 		zap.Error(err),
			// 	)
			// 	return
			// }
			// say -1 Server restart at

			// go func() {
			// 	for serverMessage := range client.Messages() {
			// 		logger.Info(
			// 			"server message",
			// 			zap.String("message", serverMessage),
			// 			zap.String("serverName", serverName),
			// 			zap.String("serverHost", serverConf.Host),
			// 			zap.Int("serverPort", serverConf.Port),
			// 		)
			// 	}
			// }()

			// message := "Your AD can be here."

			// if err := sendMessage(client, message); err != nil {
			// 	logger.Error(
			// 		"client command execute error",
			// 		zap.String("serverName", serverName),
			// 		zap.String("serverHost", serverConf.Host),
			// 		zap.Int("serverPort", serverConf.Port),
			// 		zap.Error(err),
			// 	)
			// }

			// checkTicker := time.NewTicker(15 * time.Second)
			// restartTime := time.Now().Add(serverConf.Restart.Period)
			// restart := false

			// for t := range checkTicker.C {

			// 	tDiff := restartTime.Sub(t)

			// 	if tDiff.Minutes() <= 2 {
			// 		logger.Info("SERVER NEED RESTART")
			// 		restart = true
			// 		break
			// 	}

			// 	if err := sendMessage(client, fmt.Sprintf("Server restart %f", tDiff.Minutes())); err != nil {
			// 		logger.Error(
			// 			"client command execute error",
			// 			zap.String("serverName", serverName),
			// 			zap.String("serverHost", serverConf.Host),
			// 			zap.Int("serverPort", serverConf.Port),
			// 			zap.Error(err),
			// 		)
			// 	}
			// }

			// if restart {
			// 	logger.Info("Start server restart")
			// 	if err := sendMessage(client, "Server restart now"); err != nil {
			// 		logger.Error(
			// 			"client command execute error",
			// 			zap.String("serverName", serverName),
			// 			zap.String("serverHost", serverConf.Host),
			// 			zap.Int("serverPort", serverConf.Port),
			// 			zap.Error(err),
			// 		)
			// 	}

			// 	logger.Info("Lock server")
			// 	if resp, err := client.Exec("#lock Dro4im"); err != nil {
			// 		logger.Error(
			// 			"client command execute error",
			// 			zap.String("serverName", serverName),
			// 			zap.String("serverHost", serverConf.Host),
			// 			zap.Int("serverPort", serverConf.Port),
			// 			zap.Error(err),
			// 		)
			// 	} else {
			// 		logger.Info("Server locked", zap.Any("response", resp))
			// 	}

			// 	logger.Info("Get available commands")
			// 	if resp, err := client.Exec("commands"); err != nil {
			// 		logger.Error(
			// 			"client command execute error",
			// 			zap.String("serverName", serverName),
			// 			zap.String("serverHost", serverConf.Host),
			// 			zap.Int("serverPort", serverConf.Port),
			// 			zap.Error(err),
			// 		)
			// 	} else {
			// 		logger.Info("Commands received", zap.Any("response", resp))
			// 	}

			// 	logger.Info("Kick players")
			// 	if resp, err := client.Exec("admins"); err != nil {
			// 		logger.Error(
			// 			"client command execute error",
			// 			zap.String("serverName", serverName),
			// 			zap.String("serverHost", serverConf.Host),
			// 			zap.Int("serverPort", serverConf.Port),
			// 			zap.Error(err),
			// 		)
			// 	} else {
			// 		logger.Info("Server players received", zap.Any("response", resp))
			// 	}

			// 	if resp, err := client.Exec("#kick -1"); err != nil {
			// 		logger.Error(
			// 			"client command execute error",
			// 			zap.String("serverName", serverName),
			// 			zap.String("serverHost", serverConf.Host),
			// 			zap.Int("serverPort", serverConf.Port),
			// 			zap.Error(err),
			// 		)
			// 	} else {
			// 		logger.Info("Server players kicked", zap.Any("response", resp))
			// 		timer1 := time.NewTimer(30 * time.Second)

			// 		<-timer1.C

			// 		logger.Info("Shutdown server")
			// 		if resp, err := client.Exec("#shutdown"); err != nil {
			// 			logger.Error(
			// 				"client command execute error",
			// 				zap.String("serverName", serverName),
			// 				zap.String("serverHost", serverConf.Host),
			// 				zap.Int("serverPort", serverConf.Port),
			// 				zap.Error(err),
			// 			)
			// 		} else {
			// 			logger.Info("Shutdown server performed", zap.Any("response", resp))
			// 		}
			// 	}
			// }

			// return
		}(serverName, serverConf, logger)
	}

	for s := range done {
		if s {
			return
		}
	}

	// // Create service
	// restarterService := restarter.NewService()

	// // Start timer
	// restarterService.Start()
}

func sendMessage(client *battleye.Client, message string) error {
	cmd := fmt.Sprintf("say -1 %s", message)

	_, err := client.Exec(cmd)
	if err != nil {
		return err
	}

	return nil
}
