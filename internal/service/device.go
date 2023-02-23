package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/storage"
)

func (s *Service) AddDevice(ctx context.Context, user storage.User, companyID uuid.UUID, name string, resource int) (
	device storage.AddedDevice, err error) {
	if companyID == uuid.Nil || name == "" || resource == 0 {
		return storage.AddedDevice{}, ErrWrongDeviceFormat
	}
	inputDevice := storage.Device{
		Company:  companyID,
		Name:     name,
		Resource: resource,
	}
	return s.st.AddDevice(ctx, user, inputDevice)
}

func (s *Service) GetDevice(ctx context.Context, user storage.User, deviceID uuid.UUID) (
	device storage.Device, err error) {
	inputDevice := storage.Device{ID: deviceID}
	return s.st.GetDevice(ctx, user, inputDevice)
}
