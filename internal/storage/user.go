package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (st *PsqlStorage) SignUp(ctx context.Context, login string, hash []byte, salt string) error {
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	res, err := st.db.Exec(timeoutCtx,
		"INSERT INTO users (login, password_hash, password_salt) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
		login, hash, salt)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrUserAlreadyExists
	}

	return nil
}

func (st *PsqlStorage) GetUserByID(ctx context.Context, id uuid.UUID) (user User, err error) {
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	row := st.db.QueryRow(timeoutCtx,
		"SELECT id, login, password_hash, password_salt, is_admin, last_online FROM users WHERE id = $1",
		id)
	err = row.Scan(&user.ID, &user.Login, &user.PasswordHash, &user.PasswordSalt, &user.IsAdmin, &user.LastOnline)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrUserNotFound
	}
	return
}

func (st *PsqlStorage) GetUserByLogin(ctx context.Context, login string) (user User, err error) {
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	row := st.db.QueryRow(timeoutCtx,
		"SELECT id, login, password_hash, password_salt, is_admin, last_online FROM users WHERE login = $1",
		login)
	err = row.Scan(&user.ID, &user.Login, &user.PasswordHash, &user.PasswordSalt, &user.IsAdmin, &user.LastOnline)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrUserNotFound
	}
	return
}

func (st *PsqlStorage) UpdateUserLastOnline(ctx context.Context, user User) (err error) {
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Second*10)
	defer timeoutCancel()

	_, err = st.db.Exec(timeoutCtx, "UPDATE users SET last_online = now() WHERE id = $1", user.ID)
	return err
}
