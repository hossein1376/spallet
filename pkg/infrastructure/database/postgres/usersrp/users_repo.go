package usersrp

import (
	"context"
	"database/sql"

	"github.com/hossein1376/spallet/pkg/domain/model"
)

type UsersRepo struct {
	tx *sql.Tx
}

func NewRepo(tx *sql.Tx) *UsersRepo {
	return &UsersRepo{tx: tx}
}

func (r *UsersRepo) Create(
	ctx context.Context, username string,
) (*model.User, error) {
	var id model.UserID
	err := r.tx.QueryRowContext(
		ctx, "INSERT INTO users (username) VALUES ($1) RETURNING id;", username,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &model.User{ID: id, Username: username}, nil
}
