package tags

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

func (r *Repository) GetByID(ctx context.Context, id domain.TagID) (domain.Tag, error) {
	executor := r.getExecutor(ctx)

	const query = `
		SELECT t.*, c.id as cat_id, c.name as cat_name, c.slug as cat_slug, 
		       c.description as cat_description, c.color as cat_color, c.kind as cat_kind,
		       c.created_at as cat_created_at, c.updated_at as cat_updated_at
		FROM tags t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.id = $1
	`

	rows, err := executor.Query(ctx, query, id)
	if err != nil {
		return domain.Tag{}, fmt.Errorf("query tag by ID: %w", err)
	}
	defer rows.Close()

	tag, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[tagWithCategoryModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Tag{}, domain.ErrEntityNotFound
		}

		return domain.Tag{}, fmt.Errorf("collect tag: %w", err)
	}

	return tag.toDomain(), nil
}

func (r *Repository) GetByProjectAndSlug(
	ctx context.Context,
	projectID domain.ProjectID,
	slug string,
) (domain.Tag, error) {
	executor := r.getExecutor(ctx)

	const query = `
		SELECT t.*, c.id as cat_id, c.name as cat_name, c.slug as cat_slug, 
		       c.description as cat_description, c.color as cat_color, c.kind as cat_kind,
		       c.created_at as cat_created_at, c.updated_at as cat_updated_at
		FROM tags t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.project_id = $1::uuid AND t.slug = $2
	`

	rows, err := executor.Query(ctx, query, projectID, slug)
	if err != nil {
		return domain.Tag{}, fmt.Errorf("query tag by project and slug: %w", err)
	}
	defer rows.Close()

	tag, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[tagWithCategoryModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Tag{}, domain.ErrEntityNotFound
		}

		return domain.Tag{}, fmt.Errorf("collect tag: %w", err)
	}

	return tag.toDomain(), nil
}

func (r *Repository) ListByProject(
	ctx context.Context,
	projectID domain.ProjectID,
	categoryID *domain.CategoryID,
) ([]domain.Tag, error) {
	executor := r.getExecutor(ctx)

	query := `
		SELECT t.*, c.id as cat_id, c.name as cat_name, c.slug as cat_slug, 
		       c.description as cat_description, c.color as cat_color, c.kind as cat_kind,
		       c.created_at as cat_created_at, c.updated_at as cat_updated_at
		FROM tags t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.project_id = $1::uuid
	`
	args := []interface{}{projectID}

	if categoryID != nil {
		query += ` AND t.category_id = $2`

		args = append(args, *categoryID)
	}

	query += ` ORDER BY t.name`

	rows, err := executor.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query project tags: %w", err)
	}
	defer rows.Close()

	tags, err := pgx.CollectRows(rows, pgx.RowToStructByName[tagWithCategoryModel])
	if err != nil {
		return nil, fmt.Errorf("collect project tags: %w", err)
	}

	result := make([]domain.Tag, len(tags))
	for i, tag := range tags {
		result[i] = tag.toDomain()
	}

	return result, nil
}

func (r *Repository) Create(ctx context.Context, tag *domain.TagDTO) (domain.TagID, error) {
	executor := r.getExecutor(ctx)

	const query = `
		INSERT INTO tags (project_id, category_id, name, slug, description, color)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id domain.TagID

	err := executor.QueryRow(
		ctx,
		query,
		tag.ProjectID,
		tag.CategoryID,
		tag.Name,
		tag.Slug,
		tag.Description,
		tag.Color,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("insert tag: %w", err)
	}

	return id, nil
}

func (r *Repository) Update(
	ctx context.Context,
	id domain.TagID,
	categoryID *domain.CategoryID,
	name, slug string,
	description *string,
	color *string,
) error {
	executor := r.getExecutor(ctx)

	const query = `
		UPDATE tags
		SET category_id = $1, name = $2, slug = $3, description = $4, color = $5, updated_at = NOW()
		WHERE id = $6
	`

	result, err := executor.Exec(ctx, query, categoryID, name, slug, description, color, id)
	if err != nil {
		return fmt.Errorf("update tag: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrEntityNotFound
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id domain.TagID) error {
	executor := r.getExecutor(ctx)

	const query = `DELETE FROM tags WHERE id = $1`

	result, err := executor.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete tag: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrEntityNotFound
	}

	return nil
}

func (r *Repository) CreateFromCategories(ctx context.Context, projectID domain.ProjectID) error {
	executor := r.getExecutor(ctx)

	const query = `
		INSERT INTO tags (project_id, category_id, name, slug, description, color)
		SELECT $1, c.id, c.name, c.slug, c.description, c.color
		FROM categories c
		WHERE c.kind IN ('system', 'user')
		ON CONFLICT (project_id, slug) DO NOTHING
	`

	_, err := executor.Exec(ctx, query, projectID)
	if err != nil {
		return fmt.Errorf("create tags from categories: %w", err)
	}

	return nil
}

func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}

	return r.db
}

type tagWithCategoryModel struct {
	ID             string     `db:"id"`
	ProjectID      string     `db:"project_id"`
	CategoryID     *string    `db:"category_id"`
	Name           string     `db:"name"`
	Slug           string     `db:"slug"`
	Description    *string    `db:"description"`
	Color          *string    `db:"color"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
	CatID          *string    `db:"cat_id"`
	CatName        *string    `db:"cat_name"`
	CatSlug        *string    `db:"cat_slug"`
	CatDescription *string    `db:"cat_description"`
	CatColor       *string    `db:"cat_color"`
	CatKind        *string    `db:"cat_kind"`
	CatCreatedAt   *time.Time `db:"cat_created_at"`
	CatUpdatedAt   *time.Time `db:"cat_updated_at"`
}

func (m *tagWithCategoryModel) toDomain() domain.Tag {
	var categoryID *domain.CategoryID
	if m.CategoryID != nil {
		categoryID = (*domain.CategoryID)(m.CategoryID)
	}

	tag := domain.Tag{
		ID:          domain.TagID(m.ID),
		ProjectID:   domain.ProjectID(m.ProjectID),
		CategoryID:  categoryID,
		Name:        m.Name,
		Slug:        m.Slug,
		Description: m.Description,
		Color:       m.Color,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		Category:    nil,
	}

	if m.CatID != nil {
		category := domain.Category{
			ID:          domain.CategoryID(*m.CatID),
			Name:        *m.CatName,
			Slug:        *m.CatSlug,
			Description: m.CatDescription,
			Color:       m.CatColor,
			Kind:        domain.CategoryKind(*m.CatKind),
			CreatedAt:   *m.CatCreatedAt,
			UpdatedAt:   *m.CatUpdatedAt,
		}
		tag.Category = &category
	}

	return tag
}
