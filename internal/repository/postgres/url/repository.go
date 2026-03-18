package url

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

func (r *repository) Add(ctx context.Context, shortURL string, url *model.URL, userId int) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	sql := `
	INSERT INTO url (short_url, addr, redirect_count, creation_date)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (short_url) DO NOTHING;
	`
	_, err = tx.Exec(ctx, sql, shortURL, url.Addr, url.RedirectCount, url.CreationDate)

	if err != nil {
		return err
	}

	sql = `
	INSERT INTO user_urls (user_id, short_url)
	VALUES($1, $2)
	ON CONFLICT (user_id, short_url) DO NOTHING
	`

	_, err = tx.Exec(ctx, sql, userId, shortURL)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *repository) Get(ctx context.Context, shortURL string, userId int) (*model.URL, error) {
	sql := `
	SELECT addr, redirect_count, creation_date
	FROM url u
	INNER JOIN user_urls uu ON u.short_url = uu.short_url
	WHERE uu.user_id = $1 AND uu.short_url = $2
	`
	row := r.pool.QueryRow(ctx, sql, userId, shortURL)

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

func (r *repository) GetAll(ctx context.Context, userId int) (map[string]*model.URL, error) {
	sql := `
	SELECT short_url, addr, redirect_count, creation_date
	FROM url u
	INNER JOIN user_urls uu ON u.short_url = uu.short_url
	WHERE uu.user_id = $1
	`

	rows, err := r.pool.Query(ctx, sql, userId)
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

func (r *repository) Delete(ctx context.Context, shortURL string, userId int) error {
	sql := `
	DELETE FROM user_urls WHERE short_url = $1 AND user_id = $2 
	`

	cmdTag, err := r.pool.Exec(ctx, sql, shortURL, userId)

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
