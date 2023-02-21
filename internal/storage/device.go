package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/google/uuid"
)

type AddedDevice struct {
	Device
	Token string `json:"token"`
}

func (st *PsqlStorage) AddDevice(ctx context.Context, user User, companyID uuid.UUID, name string, resource int) (
	device AddedDevice, err error) {
	company, err := st.GetUserCompany(ctx, user, companyID)
	if err != nil {
		return AddedDevice{}, err
	}

	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	var exists bool
	row := st.db.QueryRow(timeoutCtx,
		"SELECT EXISTS (SELECT FROM devices WHERE company_id = $1 AND name = $2)",
		company.ID, name)
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
		companyID, name, resource)
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

func (st *PsqlStorage) GetDevice(ctx context.Context, user User, deviceID uuid.UUID) (device Device, err error) {
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	var exists, accessToDevice bool
	row := st.db.QueryRow(timeoutCtx,
		`SELECT 
    		EXISTS(SELECT FROM devices WHERE id = $1),
    		EXISTS(SELECT FROM companies WHERE id = (SELECT company_id FROM devices WHERE id = $1) AND owner_id = $2)`,
		deviceID, user.ID)
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

	return st.GetDeviceByID(ctx, deviceID)
}

func (st *PsqlStorage) GetDeviceByID(ctx context.Context, deviceID uuid.UUID) (device Device, err error) {
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	row := st.db.QueryRow(timeoutCtx,
		"SELECT id, company_id, name, resource, last_online FROM devices WHERE id = $1",
		deviceID)
	err = row.Scan(&device.ID, &device.Company, &device.Name, &device.Resource, &device.LastOnline)
	if errors.Is(err, pgx.ErrNoRows) {
		return Device{}, ErrCompanyNotFound
	}
	if err != nil {
		return Device{}, err
	}

	device.MinutesRemaining, err = st.CalcDeviceMinutesRemaining(ctx, device.ID)
	if err != nil {
		return Device{}, err
	}

	return
}

func (st *PsqlStorage) CalcDeviceMinutesRemaining(ctx context.Context, id uuid.UUID) (
	minutesRemaining int, err error) {
	return 0, nil // todo: implement
}
