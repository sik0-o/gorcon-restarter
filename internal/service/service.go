package service

import (
	"github.com/sik0-o/gorcon-restarter/v2/internal/battleye"
	"github.com/sik0-o/gorcon-restarter/v2/internal/service/action"
)

type CommandService struct {
	beService *battleye.Service
}

func New(beService *battleye.Service) *CommandService {
	return &CommandService{
		beService: beService,
	}
}

func (cs *CommandService) Exec(action action.Action) error {
	return cs.beService.ExecCmd(action.Command())
}
