package categories

import (
	"context"
	"errors"
	"fmt"
	"time"

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

func (r *Repository) GetByID(ctx context.Context, id domain.CategoryID) (domain.Category, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM categories WHERE id = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, id)
	if err != nil {
		return domain.Category{}, fmt.Errorf("query category by ID: %w", err)
	}
	defer rows.Close()

	category, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[categoryModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Category{}, domain.ErrEntityNotFound
		}

		return domain.Category{}, fmt.Errorf("collect category: %w", err)
	}

	return category.toDomain(), nil
}

func (r *Repository) GetBySlug(ctx context.Context, slug string) (domain.Category, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM categories WHERE slug = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, slug)
	if err != nil {
		return domain.Category{}, fmt.Errorf("query category by slug: %w", err)
	}
	defer rows.Close()

	category, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[categoryModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Category{}, domain.ErrEntityNotFound
		}

		return domain.Category{}, fmt.Errorf("collect category: %w", err)
	}

	return category.toDomain(), nil
}

func (r *Repository) List(ctx context.Context) ([]domain.Category, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM categories ORDER BY name`

	rows, err := executor.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query categories: %w", err)
	}
	defer rows.Close()

	categories, err := pgx.CollectRows(rows, pgx.RowToStructByName[categoryModel])
	if err != nil {
		return nil, fmt.Errorf("collect categories: %w", err)
	}

	result := make([]domain.Category, len(categories))
	for i, category := range categories {
		result[i] = category.toDomain()
	}

	return result, nil
}

func (r *Repository) Create(ctx context.Context, category *domain.CategoryDTO) (domain.CategoryID, error) {
	executor := r.getExecutor(ctx)

	const query = `
		INSERT INTO categories (name, slug, description, color, kind)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id domain.CategoryID

	err := executor.QueryRow(
		ctx,
		query,
		category.Name,
		category.Slug,
		category.Description,
		category.Color,
		category.Kind,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("insert category: %w", err)
	}

	return id, nil
}

func (r *Repository) Update(
	ctx context.Context,
	id domain.CategoryID,
	name, slug string,
	description *string,
	color *string,
) error {
	executor := r.getExecutor(ctx)

	const query = `
		UPDATE categories
		SET name = $1, slug = $2, description = $3, color = $4, updated_at = NOW()
		WHERE id = $5
	`

	result, err := executor.Exec(ctx, query, name, slug, description, color, id)
	if err != nil {
		return fmt.Errorf("update category: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrEntityNotFound
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id domain.CategoryID) error {
	executor := r.getExecutor(ctx)

	const query = `DELETE FROM categories WHERE id = $1`

	result, err := executor.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrEntityNotFound
	}

	return nil
}

func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}

	return r.db
}

type categoryModel struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Slug        string    `db:"slug"`
	Description *string   `db:"description"`
	Color       *string   `db:"color"`
	Kind        string    `db:"kind"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (m *categoryModel) toDomain() domain.Category {
	return domain.Category{
		ID:          domain.CategoryID(m.ID),
		Name:        m.Name,
		Slug:        m.Slug,
		Description: m.Description,
		Color:       m.Color,
		Kind:        domain.CategoryKind(m.Kind),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
