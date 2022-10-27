package battleye

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/multiplay/go-battleye"
	"github.com/sik0-o/gorcon-restarter/v2/internal/config"
	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
	client *battleye.Client

	serverName    string
	serverConf    config.ServerConf
	messageBuffer []string
}

func NewService(logger *zap.Logger, serverName string, serverConf config.ServerConf) (*Service, error) {
	var host string
	var port string

	hsarr := strings.Split(serverConf.Host, ":")
	if len(hsarr) > 0 {
		host = hsarr[0]
	}
	if len(hsarr) > 1 {
		port = hsarr[1]
	}
	if serverConf.Port > 0 {
		port = strconv.Itoa(serverConf.Port)
	}

	// create beclient
	client, err := battleye.NewClient(
		host+":"+port,
		serverConf.Password,
	)
	if err != nil {
		logger.Error(
			"client initialize error",
			zap.String("serverName", serverName),
			zap.String("serverHost", serverConf.Host),
			zap.Int("serverPort", serverConf.Port),
			zap.Error(err),
		)
		return nil, err
	}

	s := &Service{
		logger:        logger,
		client:        client,
		serverName:    serverName,
		serverConf:    serverConf,
		messageBuffer: []string{},
	}

	logger.Info("BattlEye version", zap.String("response", s.BEVersion()))

	//messageBuffer := []string{}
	go func() {
		for message := range s.client.Messages() {
			s.logger.Debug("got message", zap.String("message", message))
		}
	}()

	return s, nil
}

func (s *Service) BEVersion() string {
	resp, err := s.client.Exec("version")
	if err != nil {
		s.logger.Error(
			"retrive BE version error",
			zap.String("serverName", s.serverName),
			zap.String("serverHost", s.serverConf.Host),
			zap.Int("serverPort", s.serverConf.Port),
			zap.Error(err),
		)
		return ""
	}

	return resp
}

func (s *Service) Config() config.ServerConf {
	return s.serverConf
}

func (s *Service) Logger() *zap.Logger {
	return s.logger
}

func (s *Service) Login(password string) {
	cmd := "#login"
	if len(password) > 0 {
		cmd += " " + password
	}

	resp, err := s.client.Exec(cmd)
	if err != nil {
		s.logger.Error(
			"login command error",
			zap.String("serverName", s.serverName),
			zap.String("serverHost", s.serverConf.Host),
			zap.Int("serverPort", s.serverConf.Port),
			zap.Error(err),
		)
	}

	s.logger.Info("login", zap.String("response", resp))
}

func (s *Service) ExecCmd(cmd string) error {
	resp, err := s.client.Exec(cmd)
	if err != nil {
		s.logger.Error(
			"exec command error",
			zap.String("command", cmd),
			zap.String("serverName", s.serverName),
			zap.String("serverHost", s.serverConf.Host),
			zap.Int("serverPort", s.serverConf.Port),
			zap.Error(err),
		)
		return err
	}

	s.logger.Info(
		"exec command",
		zap.String("command", cmd),
		zap.String("serverName", s.serverName),
		zap.String("serverHost", s.serverConf.Host),
		zap.Int("serverPort", s.serverConf.Port),
		zap.String("response", resp),
	)
	return nil
}

func (s *Service) DebugConsole() {
	s.ExecCmd("#debug console")
}

func (s *Service) MonitorEnable() {
	s.ExecCmd("#monitor 1")
}

func (s *Service) MonitorDisable() {
	s.ExecCmd("#monitor 0")
}

func (s *Service) KickAll() error {
	return s.ExecCmd("#kick -1")
}

func (s *Service) Lock() error {
	return s.ExecCmd("#lock")
}

func (s *Service) Unlock() error {
	return s.ExecCmd("#unlock")
}

func (s *Service) Shutdown() error {
	return s.ExecCmd("#shutdown")
}

func (s *Service) KickPlayer(playerId int, reason string) error {
	cmd := fmt.Sprintf("kick %d", playerId)
	if len(reason) > 0 {
		cmd += " " + reason
	}

	return s.ExecCmd(cmd)
}

func (s *Service) KickWithReason(reason string, playerId ...int) {
	if len(playerId) == 0 {
		playerId = append(playerId, -1)
	}

	for _, plId := range playerId {
		s.KickPlayer(plId, reason)
	}
}

// TODO: extend with avg.ping
func (s *Service) Players() {
	s.ExecCmd("players")
}

func (s *Service) Admins() {
	s.ExecCmd("admins")
}

func (s *Service) Missions() {
	s.ExecCmd("missions")
}

func (s *Service) Bans() {
	s.ExecCmd("bans")
}

func (s *Service) Say(message string, playerId ...int) error {
	if len(playerId) == 0 {
		playerId = append(playerId, -1)
	}

	for _, plId := range playerId {
		cmd := fmt.Sprintf("say %d %s", plId, message)

		if err := s.ExecCmd(cmd); err != nil {
			return err
		}

	}

	return nil
}
