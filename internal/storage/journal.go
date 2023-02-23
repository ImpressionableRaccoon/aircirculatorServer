package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

func (st *PsqlStorage) GetJournal(ctx context.Context, device Device) (journal []Journal, err error) {
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	journal = make([]Journal, 0)

	var rows pgx.Rows
	rows, err = st.db.Query(timeoutCtx,
		`SELECT id, device_id, timestamp_start, timestamp_end, done FROM journals WHERE device_id = $1
            ORDER BY timestamp_start`,
		device.ID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		j := Journal{}

		err = rows.Scan(&j.ID, &j.Device, &j.Start, &j.End, &j.Done)
		if err != nil {
			return nil, err
		}

		journal = append(journal, j)
	}

	return
}
