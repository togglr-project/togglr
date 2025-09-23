package feature_tags

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rom8726/etoggle/internal/domain"
	"github.com/rom8726/etoggle/pkg/db"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		db: pool,
	}
}

func (r *Repository) ListFeatureTags(ctx context.Context, featureID domain.FeatureID) ([]domain.Tag, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT t.*, c.id as cat_id, c.name as cat_name, c.slug as cat_slug, 
	   c.description as cat_description, c.color as cat_color, c.kind as cat_kind,
	   c.created_at as cat_created_at, c.updated_at as cat_updated_at
FROM tags t
LEFT JOIN categories c ON t.category_id = c.id
JOIN feature_tags ft ON t.id = ft.tag_id
WHERE ft.feature_id = $1
ORDER BY t.name
`

	rows, err := executor.Query(ctx, query, featureID)
	if err != nil {
		return nil, fmt.Errorf("query feature tags: %w", err)
	}
	defer rows.Close()

	tags, err := pgx.CollectRows(rows, pgx.RowToStructByName[tagWithCategoryModel])
	if err != nil {
		return nil, fmt.Errorf("collect feature tags: %w", err)
	}

	result := make([]domain.Tag, len(tags))
	for i, tag := range tags {
		result[i] = tag.toDomain()
	}

	return result, nil
}

func (r *Repository) AddFeatureTag(ctx context.Context, featureID domain.FeatureID, tagID domain.TagID) error {
	executor := r.getExecutor(ctx)

	const query = `
INSERT INTO feature_tags (feature_id, tag_id)
VALUES ($1, $2)
ON CONFLICT (feature_id, tag_id) DO NOTHING
`

	_, err := executor.Exec(ctx, query, featureID, tagID)
	if err != nil {
		return fmt.Errorf("add feature tag: %w", err)
	}

	return nil
}

func (r *Repository) RemoveFeatureTag(ctx context.Context, featureID domain.FeatureID, tagID domain.TagID) error {
	executor := r.getExecutor(ctx)

	const query = `
DELETE FROM feature_tags
WHERE feature_id = $1 AND tag_id = $2
`

	result, err := executor.Exec(ctx, query, featureID, tagID)
	if err != nil {
		return fmt.Errorf("remove feature tag: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrEntityNotFound
	}

	return nil
}

func (r *Repository) HasFeatureTag(ctx context.Context, featureID domain.FeatureID, tagID domain.TagID) (bool, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT EXISTS(
	SELECT 1 FROM feature_tags
	WHERE feature_id = $1 AND tag_id = $2
)
`

	var exists bool
	err := executor.QueryRow(ctx, query, featureID, tagID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check feature tag: %w", err)
	}

	return exists, nil
}

func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}
	return r.db
}

// Reuse models from tags repository
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
