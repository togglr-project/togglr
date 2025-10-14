package algorithms

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
)

type Repository struct {
	db db.Tx
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		db: pool,
	}
}

func (r *Repository) GetBySlug(ctx context.Context, slug string) (domain.Algorithm, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM algorithms WHERE slug = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, slug)
	if err != nil {
		return domain.Algorithm{}, fmt.Errorf("query algorithm by slug: %w", err)
	}
	defer rows.Close()

	alg, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[algorithmModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Algorithm{}, domain.ErrEntityNotFound
		}

		return domain.Algorithm{}, fmt.Errorf("collect algorithm: %w", err)
	}

	return alg.toDomain(), nil
}

func (r *Repository) List(ctx context.Context) ([]domain.Algorithm, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM algorithms ORDER BY slug`

	rows, err := executor.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query algorithms: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[algorithmModel])
	if err != nil {
		return nil, fmt.Errorf("collect algorithms: %w", err)
	}

	result := make([]domain.Algorithm, 0, len(models))
	for _, m := range models {
		result = append(result, m.toDomain())
	}

	return result, nil
}

func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}

	return r.db
}
