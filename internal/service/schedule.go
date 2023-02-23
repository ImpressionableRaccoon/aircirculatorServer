package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/storage"
)

func (s *Service) GetSchedules(ctx context.Context, user storage.User, deviceID uuid.UUID) (
	schedules []storage.Schedule, err error) {
	inputDevice := storage.Device{ID: deviceID}
	device, err := s.st.GetDevice(ctx, user, inputDevice)
	if err != nil {
		return nil, err
	}

	return s.st.GetSchedules(ctx, device)
}

func (s *Service) AddSchedules(ctx context.Context, user storage.User, deviceID uuid.UUID,
	inputSchedules []storage.Schedule) (schedules []storage.Schedule, err error) {
	inputDevice := storage.Device{ID: deviceID}
	device, err := s.st.GetDevice(ctx, user, inputDevice)
	if err != nil {
		return nil, err
	}

	for _, schedule := range inputSchedules {
		if schedule.Week < 0 || schedule.Week > 6 {
			return nil, ErrWrongScheduleWeek
		}
		if schedule.TimeStart < 0 || schedule.TimeStart >= 1440 {
			return nil, ErrWrongScheduleTimeStart
		}
		if schedule.TimeStop < 0 || schedule.TimeStop >= 1440 {
			return nil, ErrWrongScheduleTimeStop
		}
		if schedule.TimeStop <= schedule.TimeStart {
			return nil, ErrTimeStopNotMoreTimeStart
		}
	}

	return s.st.AddSchedules(ctx, device, inputSchedules)
}

func (s *Service) DeleteSchedule(ctx context.Context, user storage.User, deviceID uuid.UUID, scheduleID uuid.UUID) (
	err error) {
	inputDevice := storage.Device{ID: deviceID}
	device, err := s.st.GetDevice(ctx, user, inputDevice)
	if err != nil {
		return err
	}

	if scheduleID == uuid.Nil {
		return ErrWrongScheduleID
	}

	return s.st.DeleteSchedule(ctx, device, scheduleID)
}
