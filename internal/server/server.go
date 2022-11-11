package server

import (
	"github.com/sik0-o/gorcon-restarter/v2/internal/service"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger
	cs     *service.CommandService
}
