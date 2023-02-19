package storage

import (
	"context"
	"time"
)

func (st *PsqlStorage) SignUp(ctx context.Context, login string, hash []byte, salt string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	res, err := st.db.Exec(ctx,
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

func (st *PsqlStorage) GetUser(ctx context.Context, login string) (user User, err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	row := st.db.QueryRow(ctx,
		"SELECT id, login, password_hash, password_salt, is_admin, last_online FROM users WHERE login = $1", login)
	err = row.Scan(&user.ID, &user.Login, &user.PasswordHash, &user.PasswordSalt, &user.IsAdmin, &user.LastOnline)
	return
}
