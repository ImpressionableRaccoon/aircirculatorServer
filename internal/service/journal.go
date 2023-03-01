package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/storage"
)

func (s *Service) GetJournalWithCheck(ctx context.Context, user storage.User, deviceID uuid.UUID) (
	journal []storage.Journal, err error) {
	inputDevice := storage.Device{ID: deviceID}
	device, err := s.GetDevice(ctx, user, inputDevice)
	if err != nil {
		return nil, err
	}

	return s.GetJournal(ctx, device)
}
