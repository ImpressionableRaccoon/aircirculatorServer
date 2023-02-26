package workers

import (
	"github.com/ImpressionableRaccoon/aircirculatorServer/configs"
	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/service"
)

type Workers struct {
	s    *service.Service
	cfg  *configs.Config
	done chan struct{}
}

func NewWorkers(s *service.Service, cfg *configs.Config) Workers {
	return Workers{
		s:    s,
		cfg:  cfg,
		done: make(chan struct{}),
	}
}

func (w *Workers) Stop() {
	close(w.done)
}
