package ruleattributes

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rom8726/etoggle/internal/domain"
	"github.com/rom8726/etoggle/pkg/db"
)

type Repository struct {
	// db allows executing queries either on a transaction or a connection pool
	// depending on context (see getExecutor)
	db db.Tx
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{db: pool}
}

type attributeModel struct {
	ID          uint           `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
}

func (m attributeModel) toDomain() domain.RuleAttributeEntity {
	var desc *string
	if m.Description.Valid {
		desc = &m.Description.String
	}

	return domain.RuleAttributeEntity{
		ID:          m.ID,
		Name:        domain.RuleAttribute(m.Name),
		Description: desc,
	}
}

func (r *Repository) Create(
	ctx context.Context,
	name domain.RuleAttribute,
	description *string,
) (domain.RuleAttributeEntity, error) {
	exec := r.getExecutor(ctx)

	const query = `
INSERT INTO rule_attributes (name, description)
VALUES ($1, $2)
RETURNING id, name, description`

	var model attributeModel

	var desc any
	if description != nil {
		desc = *description
	} else {
		desc = sql.NullString{}
	}

	err := exec.QueryRow(ctx, query, name, desc).Scan(
		&model.ID,
		&model.Name,
		&model.Description,
	)
	if err != nil {
		return domain.RuleAttributeEntity{}, fmt.Errorf("insert rule attribute: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) Delete(ctx context.Context, name domain.RuleAttribute) error {
	exec := r.getExecutor(ctx)

	const query = `DELETE FROM rule_attributes WHERE name = $1`

	_, err := exec.Exec(ctx, query, name)
	if err != nil {
		return fmt.Errorf("delete rule attribute: %w", err)
	}

	return nil
}

func (r *Repository) List(ctx context.Context) ([]domain.RuleAttributeEntity, error) {
	exec := r.getExecutor(ctx)
	const query = `SELECT id, name, description FROM rule_attributes ORDER BY name`

	rows, err := exec.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query rule attributes: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[attributeModel])
	if err != nil {
		return nil, fmt.Errorf("collect rule attributes: %w", err)
	}

	result := make([]domain.RuleAttributeEntity, 0, len(models))
	for i := range models {
		result = append(result, models[i].toDomain())
	}

	return result, nil
}

//nolint:ireturn // repository pattern requires interface return
func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}

	return r.db
}
