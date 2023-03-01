package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type NotificationDevice struct {
	ID                   uuid.UUID
	Owner                string
	Company              string
	Name                 string
	Resource             int
	MinutesRemaining     int
	LastOnline           time.Time
	TimeOffset           time.Duration
	LastOnlineWithOffset time.Time
	LastOnlineDuration   time.Duration
}

func (st *PsqlStorage) GetAllDevices(ctx context.Context) (devices []NotificationDevice, err error) {
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	devices = make([]NotificationDevice, 0)

	var rows pgx.Rows
	rows, err = st.db.Query(timeoutCtx,
		`SELECT
		id,
		(SELECT login FROM users WHERE id = (SELECT owner_id FROM companies WHERE id = company_id)) AS owner_login,
		(SELECT name FROM companies WHERE id = company_id) AS company_name,
		name AS device_name,
		resource,
		(SELECT resource-EXTRACT(EPOCH FROM COALESCE(SUM(timestamp_end - timestamp_start), '0 days')::INTERVAL)/60
				FROM journals WHERE device_id = devices.id)::INTEGER as minutes_remaining,
		last_online,
		(SELECT time_offset FROM companies WHERE id = company_id)
		FROM devices`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		device := NotificationDevice{}

		err = rows.Scan(
			&device.ID,
			&device.Owner,
			&device.Company,
			&device.Name,
			&device.Resource,
			&device.MinutesRemaining,
			&device.LastOnline,
			&device.TimeOffset,
		)
		if err != nil {
			return nil, err
		}

		locationName := fmt.Sprintf(
			"%02d:%02d:%02d",
			int(device.TimeOffset.Hours())%24,
			int(device.TimeOffset.Minutes())%60,
			int(device.TimeOffset.Seconds())%60)
		loc := time.FixedZone(locationName, int(device.TimeOffset.Seconds()))
		device.LastOnlineWithOffset = device.LastOnline.In(loc)

		device.LastOnlineDuration = time.Now().UTC().Sub(device.LastOnline)

		devices = append(devices, device)
	}

	return devices, nil
}
