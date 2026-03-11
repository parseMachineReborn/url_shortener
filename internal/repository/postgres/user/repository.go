package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/parseMachineReborn/url_shortener/internal/apperror"
	"github.com/parseMachineReborn/url_shortener/internal/model"
)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool: pool,
	}
}

func (r *repository) Register(ctx context.Context, email string, passHash string) error {
	sql := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		ON CONFLICT (email) DO NOTHING
	`

	cmdTag, err := r.pool.Exec(ctx, sql, email, passHash)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return apperror.ErrAlreadyRegistered
	}

	return nil
}

func (r *repository) GetUser(ctx context.Context, email string) (model.User, error) {
	sql := `
		SELECT id, email, password_hash
		FROM users
		WHERE email = $1
	`

	var user model.User
	row := r.pool.QueryRow(ctx, sql, email)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PassHash,
	)

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
