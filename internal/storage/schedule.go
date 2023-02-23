package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (st *PsqlStorage) GetSchedules(ctx context.Context, device Device) (schedules []Schedule, err error) {
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	schedules = make([]Schedule, 0)

	var rows pgx.Rows
	rows, err = st.db.Query(timeoutCtx,
		`SELECT id, device_id, week_day, time_start, time_stop
			FROM schedules WHERE device_id = $1
			ORDER BY week_day, time_start`,
		device.ID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		schedule := Schedule{}

		err = rows.Scan(&schedule.ID, &schedule.Device, &schedule.Week, &schedule.TimeStart, &schedule.TimeStop)
		if err != nil {
			return nil, err
		}

		schedules = append(schedules, schedule)
	}

	return
}

func (st *PsqlStorage) AddSchedules(ctx context.Context, device Device, inputSchedules []Schedule) (
	schedules []Schedule, err error) {
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*60)
	defer timeoutCancel()

	weekDays := make([]int, 0, len(inputSchedules))
	timeStarts := make([]int, 0, len(inputSchedules))
	timeStops := make([]int, 0, len(inputSchedules))

	for _, schedule := range inputSchedules {
		weekDays = append(weekDays, schedule.Week)
		timeStarts = append(timeStarts, schedule.TimeStart)
		timeStops = append(timeStops, schedule.TimeStop)
	}

	rows, err := st.db.Query(timeoutCtx,
		`INSERT INTO schedules (device_id, week_day, time_start, time_stop)
		SELECT $1, data_table.week_day, data_table.time_start, data_table.time_stop
		FROM (SELECT unnest($2::integer[]) AS week_day,
		             unnest($3::integer[]) AS time_start,
		             unnest($4::integer[]) AS time_stop) AS data_table
		RETURNING id, device_id, week_day, time_start, time_stop`,
		device.ID, weekDays, timeStarts, timeStops)
	if err != nil {
		return nil, err
	}

	schedules = make([]Schedule, 0, len(inputSchedules))

	for rows.Next() {
		schedule := Schedule{}

		err = rows.Scan(&schedule.ID, &schedule.Device, &schedule.Week, &schedule.TimeStart, &schedule.TimeStop)
		if err != nil {
			return nil, err
		}

		schedules = append(schedules, schedule)
	}

	return
}

func (st *PsqlStorage) DeleteSchedule(ctx context.Context, device Device, scheduleID uuid.UUID) (err error) {
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	res, err := st.db.Exec(timeoutCtx,
		`DELETE FROM schedules WHERE id = $1 AND device_id = $2`,
		scheduleID, device.ID)
	if err != nil {
		return err
	}
	if res.RowsAffected() != 1 {
		return ErrDeleteScheduleFailed
	}

	return nil
}
