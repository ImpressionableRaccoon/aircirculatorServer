package middlewares

import (
	"github.com/ImpressionableRaccoon/aircirculatorServer/configs"
	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/handlers"
	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/service"
)

type Middlewares struct {
	s   *service.Service
	h   *handlers.Handler
	cfg *configs.Config
}

func NewMiddlewares(s *service.Service, h *handlers.Handler, cfg *configs.Config) Middlewares {
	return Middlewares{
		s:   s,
		h:   h,
		cfg: cfg,
	}
}
