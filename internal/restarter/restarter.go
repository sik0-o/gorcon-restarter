package restarter

import (
	"fmt"
	"math"
	"time"

	beservice "github.com/sik0-o/gorcon-restarter/v2/internal/battleye"
	"go.uber.org/zap"
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

		// Если до отключения осталось времени меньше или равное времени перед локом, то выходим из цикла и начинаем рестарт.
		if tDiff.Minutes() <= lockTime {
			s.beService.Logger().Info("SERVER NEED RESTART")
			restart = true
			break
		}
		checkCounter++

		s.beService.Logger().Debug("Counter", zap.Float64("seconds", tDiff.Seconds()), zap.Int("counter", checkCounter))

		if tDiff.Seconds() >= 3600 && checkCounter%3600 == 0 {
			s.beService.Say(fmt.Sprintf(s.beService.Config().Restart.Announcements.At, restartTime.String()))
		} else if tDiff.Seconds() > 30 && checkCounter%30 == 0 {
			s.beService.Say(fmt.Sprintf(s.beService.Config().Restart.Announcements.Min, math.Round(tDiff.Minutes())))
		} else if tDiff.Seconds() <= 30 {
			s.beService.Say(fmt.Sprintf(s.beService.Config().Restart.Announcements.Sec, math.Round(tDiff.Seconds())))
		}
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
