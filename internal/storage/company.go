package storage

import (
	"context"
	"errors"
	"time"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/utils"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"
)

func (st *PsqlStorage) GetUserCompanies(ctx context.Context, user User, ignoreUser bool) (companies []Company, err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	companies = make([]Company, 0)

	var rows pgx.Rows
	if ignoreUser {
		rows, err = st.db.Query(ctx, "SELECT id, owner_id, name, time_offset FROM companies")
	} else {
		rows, err = st.db.Query(ctx,
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
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var exists bool
	row := st.db.QueryRow(ctx,
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
	row = st.db.QueryRow(ctx,
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
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var exists, accessToDevice bool
	row := st.db.QueryRow(ctx,
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
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var d time.Duration
	row := st.db.QueryRow(ctx,
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

	return company, nil
}
