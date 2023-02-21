package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/storage"
)

func (s *Service) AddDevice(ctx context.Context, user storage.User, companyID uuid.UUID, name string, resource int) (
	device storage.AddedDevice, err error) {
	return s.st.AddDevice(ctx, user, companyID, name, resource)
}

func (s *Service) GetDevice(ctx context.Context, user storage.User, deviceID uuid.UUID) (
	device storage.Device, err error) {
	return s.st.GetDevice(ctx, user, deviceID)
}
