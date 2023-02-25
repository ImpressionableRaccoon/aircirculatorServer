package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/utils"
)

func (st *PsqlStorage) GetUserCompanies(ctx context.Context, user User, ignoreUser bool) (companies []Company, err error) {
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	companies = make([]Company, 0)

	var rows pgx.Rows
	if ignoreUser {
		rows, err = st.db.Query(timeoutCtx, "SELECT id, owner_id, name, time_offset FROM companies")
	} else {
		rows, err = st.db.Query(timeoutCtx,
			"SELECT id, owner_id, name, time_offset FROM companies WHERE owner_id = $1",
			user.ID)
	}
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		company := Company{}
		var d time.Duration

		err = rows.Scan(&company.ID, &company.Owner, &company.Name, &d)
		if err != nil {
			return nil, err
		}

		company.TimeOffset = utils.Offset{Duration: d}

		companies = append(companies, company)
	}

	return
}

func (st *PsqlStorage) AddCompany(ctx context.Context, user User, name string, offset utils.Offset) (
	company Company, err error) {
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	var exists bool
	row := st.db.QueryRow(timeoutCtx,
		"SELECT EXISTS (SELECT FROM companies WHERE owner_id = $1 AND name = $2)",
		user.ID, name)
	err = row.Scan(&exists)
	if err != nil {
		return Company{}, err
	}
	if exists {
		return Company{}, ErrCompanyAlreadyExists
	}

	var companyID uuid.UUID
	row = st.db.QueryRow(timeoutCtx,
		"INSERT INTO companies (owner_id, name, time_offset) VALUES ($1, $2, $3) RETURNING id",
		user.ID, name, offset.Duration)
	err = row.Scan(&companyID)
	if err != nil {
		return Company{}, err
	}

	return st.GetCompanyByID(ctx, companyID)
}

func (st *PsqlStorage) GetUserCompany(ctx context.Context, user User, id uuid.UUID) (
	company Company, err error) {
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	var exists, accessToDevice bool
	row := st.db.QueryRow(timeoutCtx,
		`SELECT 
    		EXISTS(SELECT FROM companies WHERE id = $1),
    		EXISTS(SELECT FROM companies WHERE id = $1 AND owner_id = $2)`,
		id, user.ID)
	err = row.Scan(&exists, &accessToDevice)
	if err != nil {
		return Company{}, err
	}

	if !exists {
		return Company{}, ErrCompanyNotFound
	}
	if !accessToDevice && !user.IsAdmin {
		return Company{}, ErrCompanyNoPermissions
	}

	return st.GetCompanyByID(ctx, id)
}

func (st *PsqlStorage) GetCompanyByID(ctx context.Context, id uuid.UUID) (company Company, err error) {
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	var d time.Duration
	row := st.db.QueryRow(timeoutCtx,
		"SELECT id, owner_id, name, time_offset FROM companies WHERE id = $1",
		id)
	err = row.Scan(&company.ID, &company.Owner, &company.Name, &d)
	if errors.Is(err, pgx.ErrNoRows) {
		return Company{}, ErrCompanyNotFound
	}
	if err != nil {
		return Company{}, err
	}
	company.TimeOffset = utils.Offset{Duration: d}

	return
}

func (st *PsqlStorage) GetCompanyDevices(ctx context.Context, user User, id uuid.UUID) (
	devices []Device, err error) {
	company, err := st.GetUserCompany(ctx, user, id)
	if err != nil {
		return nil, err
	}

	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	devices = make([]Device, 0)

	var rows pgx.Rows
	rows, err = st.db.Query(timeoutCtx,
		"SELECT id, company_id, name, resource, last_online FROM devices WHERE company_id = $1",
		company.ID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		device := Device{}

		err = rows.Scan(&device.ID, &device.Company, &device.Name, &device.Resource, &device.LastOnline)
		if err != nil {
			return nil, err
		}

		var sum int
		sum, err = st.GetJournalSum(ctx, device)
		device.MinutesRemaining = device.Resource - sum

		devices = append(devices, device)
	}

	return
}
