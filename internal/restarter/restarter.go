package restarter

import (
	"fmt"
	"math"
	"time"

	beservice "github.com/sik0-o/gorcon-restarter/v2/internal/battleye"
	"github.com/sik0-o/gorcon-restarter/v2/internal/service"
	act "github.com/sik0-o/gorcon-restarter/v2/internal/service/action"
	"go.uber.org/zap"
)

type Service struct {
	logger    *zap.Logger
	beService *beservice.Service
	cs        *service.CommandService
}

// TODO: remove beService
func NewService(logger *zap.Logger, beService *beservice.Service, cs *service.CommandService) (*Service, error) {
	return &Service{
		logger:    logger,
		beService: beService,
		cs:        cs,
	}, nil
}

func (s *Service) Restart() {
	checkTicker := time.NewTicker(time.Second)
	restartTime := time.Now().Add(s.beService.Config().Restart.Period)
	restart := false

	lockTime := s.beService.Config().Restart.ServerLock.Minutes()
	if lockTime <= 0 {
		lockTime = 1
	}

	checkCounter := 0

	// Ожидаем время до рестарта
	for t := range checkTicker.C {

		// Время до рестарта
		tDiff := restartTime.Sub(t)
		tdSec := math.Round(tDiff.Seconds())
		tdMin := math.Round(tDiff.Minutes())

		// Если до отключения осталось времени меньше или равное времени перед локом, то выходим из цикла и начинаем рестарт.
		if tdMin <= lockTime {
			s.logger.Info("SERVER NEED RESTART")
			restart = true
			break
		}
		checkCounter++

		s.logger.Debug("Counter", zap.Float64("seconds", tdSec), zap.Int("counter", checkCounter))

		var restartMessage string
		if tdSec >= 3600 && checkCounter%300 == 0 {
			restartMessage = fmt.Sprintf(s.beService.Config().Restart.Announcements.At, math.Floor(tDiff.Hours()))
		} else if tdSec < 3600 && tdSec > 30 && checkCounter%30 == 0 {
			restartMessage = fmt.Sprintf(s.beService.Config().Restart.Announcements.Min, math.Round(tdMin))
		} else if tdSec <= 30 {
			restartMessage = fmt.Sprintf(s.beService.Config().Restart.Announcements.Sec, math.Round(tdSec))
		}

		if len(restartMessage) > 0 {
			s.logger.Debug("send restart message to players", zap.String("messages", restartMessage))
			s.cs.Exec(act.NewAnnounce(restartMessage))
		}
	}

	if restart {
		s.restart()
	}
}

func (s *Service) restart() {
	s.logger.Info("Start server restart")
	s.cs.Exec(act.NewAnnounce("Server restart now"))

	s.logger.Info("Lock server")
	s.cs.Exec(act.NewLock())

	timer1 := time.NewTimer(10 * time.Second)
	<-timer1.C

	s.logger.Info("Kick players")
	if err := s.cs.Exec(act.NewKick(-1)); err != nil {
		//s.logger.Error("error", zap.Error(err))
		return
	}

	s.logger.Info("Server players kicked")

	timer2 := time.NewTimer(30 * time.Second)
	<-timer2.C

	s.logger.Info("Shutdown server")

	if err := s.cs.Exec(act.NewShutdown()); err != nil {
		return
	}

	s.logger.Info("Shutdown server performed")
}
