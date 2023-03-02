package service

import (
	"github.com/ImpressionableRaccoon/aircirculatorServer/configs"
	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/storage"
)

type Service struct {
	*storage.PsqlStorage
	cfg *configs.Config
}

func NewService(st *storage.PsqlStorage, cfg *configs.Config) *Service {
	s := &Service{
		PsqlStorage: st,
		cfg:         cfg,
	}

	return s
}
