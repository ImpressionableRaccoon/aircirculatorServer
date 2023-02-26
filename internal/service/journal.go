package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/storage"
)

func (s *Service) GetJournal(ctx context.Context, user storage.User, deviceID uuid.UUID) (
	journal []storage.Journal, err error) {
	inputDevice := storage.Device{ID: deviceID}
	device, err := s.st.GetDevice(ctx, user, inputDevice)
	if err != nil {
		return nil, err
	}

	return s.st.GetJournal(ctx, device)
}

func (s *Service) DropShortJournals(ctx context.Context) (deleted int, err error) {
	return s.st.DropShortJournals(ctx)
}
