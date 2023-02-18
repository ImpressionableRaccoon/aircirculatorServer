package middlewares

import (
	"github.com/ImpressionableRaccoon/aircirculatorServer/configs"
)

type Middlewares struct {
	cfg *configs.Config
}

func NewMiddlewares(cfg *configs.Config) Middlewares {
	return Middlewares{
		cfg: cfg,
	}
}
