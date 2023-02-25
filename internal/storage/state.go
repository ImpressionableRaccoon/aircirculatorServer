package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/utils"
)

type DeviceState struct {
	Time    time.Time                `json:"time"`
	Mode    string                   `json:"mode"`
	Enabled utils.ConvertibleBoolean `json:"enabled"`
}

func (st *PsqlStorage) PushDeviceState(ctx context.Context, device Device, state DeviceState) (err error) {
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	var (
		id  uuid.UUID
		end time.Time
	)
	row := st.db.QueryRow(timeoutCtx,
		`SELECT id, timestamp_end FROM journals
		WHERE device_id = $1 AND done = 'false' 
	    ORDER BY timestamp_start DESC LIMIT 1`,
		device.ID)
	err = row.Scan(&id, &end)

	if !state.Enabled {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		if err != nil {
			return err
		}
		return st.doneJournalWithTime(ctx, id, state)
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return st.createNewJournal(ctx, device, state)
	}
	if err != nil {
		return err
	}

	if state.Time.Sub(end) <= st.cfg.JournalTTL {
		return st.appendJournal(ctx, id, state)
	}

	err = st.doneJournal(ctx, id)
	if err != nil {
		return err
	}

	return st.createNewJournal(ctx, device, state)
}

func (st *PsqlStorage) createNewJournal(ctx context.Context, device Device, state DeviceState) (err error) {
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	_, err = st.db.Exec(timeoutCtx,
		`INSERT INTO journals (device_id, timestamp_start, timestamp_end) VALUES ($1, $2, $2)`,
		device.ID, state.Time)
	if err != nil {
		return err
	}

	return nil
}

func (st *PsqlStorage) appendJournal(ctx context.Context, journalID uuid.UUID, state DeviceState) (err error) {
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	_, err = st.db.Exec(timeoutCtx, `UPDATE journals SET timestamp_end = $1 WHERE id = $2`, state.Time, journalID)
	if err != nil {
		return err
	}

	return nil
}

func (st *PsqlStorage) doneJournal(ctx context.Context, journalID uuid.UUID) (err error) {
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	_, err = st.db.Exec(timeoutCtx, `UPDATE journals SET done = 'true' WHERE id = $1`, journalID)
	if err != nil {
		return err
	}

	return nil
}

func (st *PsqlStorage) doneJournalWithTime(ctx context.Context, journalID uuid.UUID, state DeviceState) (err error) {
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	_, err = st.db.Exec(timeoutCtx,
		`UPDATE journals SET done = 'true', timestamp_end = least(timestamp_end + $1, $2) WHERE id = $3`,
		st.cfg.JournalTTL, state.Time, journalID)
	if err != nil {
		return err
	}

	return nil
}
