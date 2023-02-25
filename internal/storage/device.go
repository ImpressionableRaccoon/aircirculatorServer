package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AddedDevice struct {
	Device
	Token string `json:"token"`
}

func (st *PsqlStorage) AddDevice(ctx context.Context, user User, inputDevice Device) (
	device AddedDevice, err error) {
	company, err := st.GetUserCompany(ctx, user, inputDevice.Company)
	if err != nil {
		return AddedDevice{}, err
	}

	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	var exists bool
	row := st.db.QueryRow(timeoutCtx,
		"SELECT EXISTS (SELECT FROM devices WHERE company_id = $1 AND name = $2)",
		company.ID, inputDevice.Name)
	err = row.Scan(&exists)
	if err != nil {
		return AddedDevice{}, err
	}
	if exists {
		return AddedDevice{}, ErrDeviceAlreadyExists
	}

	var deviceID uuid.UUID
	var deviceToken string
	row = st.db.QueryRow(timeoutCtx,
		"INSERT INTO devices (company_id, name, resource) VALUES ($1, $2, $3) RETURNING id, token",
		inputDevice.Company, inputDevice.Name, inputDevice.Resource)
	err = row.Scan(&deviceID, &deviceToken)
	if err != nil {
		return AddedDevice{}, err
	}

	d, err := st.GetDeviceByID(ctx, deviceID)
	if err != nil {
		return AddedDevice{}, err
	}

	device = AddedDevice{
		Device: d,
		Token:  deviceToken,
	}

	return
}

func (st *PsqlStorage) GetDevice(ctx context.Context, user User, inputDevice Device) (device Device, err error) {
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	var exists, accessToDevice bool
	row := st.db.QueryRow(timeoutCtx,
		`SELECT 
    		EXISTS(SELECT FROM devices WHERE id = $1),
    		EXISTS(SELECT FROM companies WHERE id = (SELECT company_id FROM devices WHERE id = $1) AND owner_id = $2)`,
		inputDevice.ID, user.ID)
	err = row.Scan(&exists, &accessToDevice)
	if err != nil {
		return Device{}, err
	}

	if !exists {
		return Device{}, ErrDeviceNotFound
	}
	if !accessToDevice && !user.IsAdmin {
		return Device{}, ErrCompanyNoPermissions
	}

	return st.GetDeviceByID(ctx, inputDevice.ID)
}

func (st *PsqlStorage) GetDeviceByID(ctx context.Context, deviceID uuid.UUID) (device Device, err error) {
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	row := st.db.QueryRow(timeoutCtx,
		"SELECT id, company_id, name, resource, last_online FROM devices WHERE id = $1",
		deviceID)
	err = row.Scan(&device.ID, &device.Company, &device.Name, &device.Resource, &device.LastOnline)
	if errors.Is(err, pgx.ErrNoRows) {
		return Device{}, ErrDeviceNotFound
	}
	if err != nil {
		return Device{}, err
	}

	var sum int
	sum, err = st.GetJournalSum(ctx, device)
	if err != nil {
		return Device{}, err
	}
	device.MinutesRemaining = device.Resource - sum

	return device, nil
}

func (st *PsqlStorage) AuthDevice(ctx context.Context, deviceID uuid.UUID, token string) (device Device, err error) {
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	row := st.db.QueryRow(timeoutCtx,
		"SELECT id, company_id, name, resource, last_online FROM devices WHERE id = $1 AND token = $2",
		deviceID, token)
	err = row.Scan(&device.ID, &device.Company, &device.Name, &device.Resource, &device.LastOnline)
	if errors.Is(err, pgx.ErrNoRows) {
		return Device{}, ErrDeviceNotFound
	}
	if err != nil {
		return Device{}, err
	}

	var sum int
	sum, err = st.GetJournalSum(ctx, device)
	if err != nil {
		return Device{}, err
	}
	device.MinutesRemaining = device.Resource - sum

	return device, nil
}
