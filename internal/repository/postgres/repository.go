package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/parseMachineReborn/url_shortener/internal/apperror"
	"github.com/parseMachineReborn/url_shortener/internal/model"
)

type repository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool: pool,
	}
}

func (r *repository) Add(ctx context.Context, shortenURL string, url *model.URL) error {
	sql := `
	INSERT INTO url (short_url, addr, redirect_count, creation_date)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (short_url) DO NOTHING;
	`

	_, err := r.pool.Exec(ctx, sql, shortenURL, url.Addr, url.RedirectCount, url.CreationDate)

	return err
}

func (r *repository) Get(ctx context.Context, shortURL string) (*model.URL, error) {
	sql := `
	SELECT addr, redirect_count, creation_date
	FROM url
	WHERE short_url = $1
	`
	row := r.pool.QueryRow(ctx, sql, shortURL)

	var url model.URL
	if err := row.Scan(
		&url.Addr,
		&url.RedirectCount,
		&url.CreationDate,
	); err != nil {
		return nil, err
	}

	return &url, nil
}

func (r *repository) GetAll(ctx context.Context) (map[string]*model.URL, error) {
	sql := `
	SELECT short_url, addr, redirect_count, creation_date
	FROM url
	`

	rows, err := r.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := make(map[string]*model.URL)
	for rows.Next() {
		var shortURL string
		var url model.URL
		err := rows.Scan(
			&shortURL,
			&url.Addr,
			&url.RedirectCount,
			&url.CreationDate,
		)

		if err != nil {
			return nil, err
		}

		res[shortURL] = &url
	}

	return res, nil
}

func (r *repository) Delete(ctx context.Context, shortURL string) error {
	sql := `
	DELETE FROM url WHERE short_url = $1
	`

	cmdTag, err := r.pool.Exec(ctx, sql, shortURL)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return apperror.ErrNotFound
	}

	return nil
}

func (r *repository) IncrementRedirectCount(ctx context.Context, shortURL string) error {
	sql := `
	UPDATE url SET redirect_count = redirect_count + 1 WHERE short_url = $1
	`

	cmdTag, err := r.pool.Exec(ctx, sql, shortURL)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return apperror.ErrNotFound
	}

	return nil
}
