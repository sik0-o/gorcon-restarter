package restarter

import (
	"fmt"
	"time"

	beservice "github.com/sik0-o/gorcon-restarter/v2/internal/battleye"
)

type Service struct {
	beService *beservice.Service
}

func NewService(beService *beservice.Service) (*Service, error) {
	return &Service{
		beService: beService,
	}, nil
}

func (s *Service) Restart() {
	checkTicker := time.NewTicker(15 * time.Second)
	restartTime := time.Now().Add(s.beService.Config().Restart.Period)
	restart := false

	// Check
	for t := range checkTicker.C {

		tDiff := restartTime.Sub(t)

		if tDiff.Minutes() <= 2 {
			s.beService.Logger().Info("SERVER NEED RESTART")
			restart = true
			break
		}

		s.beService.Say(fmt.Sprintf("Server restart %f", tDiff.Minutes()))
	}

	if restart {
		s.beService.Logger().Info("Start server restart")
		s.beService.Say("Server restart now")

		s.beService.Logger().Info("Lock server")
		s.beService.Lock()

		s.beService.Logger().Info("Kick players")
		if err := s.beService.KickAll(); err != nil {
			//s.beService.Logger().Error("error", zap.Error(err))
			return
		}

		s.beService.Logger().Info("Server players kicked")
		timer1 := time.NewTimer(30 * time.Second)

		<-timer1.C

		s.beService.Logger().Info("Shutdown server")

		if err := s.beService.Shutdown(); err != nil {
			return
		}

		s.beService.Logger().Info("Shutdown server performed")

	}
}
