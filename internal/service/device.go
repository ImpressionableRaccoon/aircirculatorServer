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
	return s.AddUserDevice(ctx, user, inputDevice)
}
