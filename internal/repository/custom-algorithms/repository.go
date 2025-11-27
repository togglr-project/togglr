package customalgorithms

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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

func (r *Repository) Create(ctx context.Context, dto domain.CustomAlgorithmDTO) (domain.CustomAlgorithm, error) {
	executor := r.getExecutor(ctx)

	hash := sha256.Sum256(dto.WASMBinary)
	wasmHash := hex.EncodeToString(hash[:])

	settingsJSON, err := json.Marshal(dto.DefaultSettings)
	if err != nil {
		return domain.CustomAlgorithm{}, fmt.Errorf("marshal settings: %w", err)
	}

	const query = `
INSERT INTO custom_algorithms (slug, name, description, kind, wasm_binary, wasm_hash, default_settings, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *`

	rows, err := executor.Query(ctx, query,
		dto.Slug, dto.Name, dto.Description, string(dto.Kind),
		dto.WASMBinary, wasmHash, settingsJSON, dto.CreatedBy)
	if err != nil {
		return domain.CustomAlgorithm{}, fmt.Errorf("insert custom algorithm: %w", err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[customAlgorithmModel])
	if err != nil {
		return domain.CustomAlgorithm{}, fmt.Errorf("collect row: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) Update(ctx context.Context, alg domain.CustomAlgorithm) error {
	executor := r.getExecutor(ctx)

	hash := sha256.Sum256(alg.WASMBinary)
	wasmHash := hex.EncodeToString(hash[:])

	settingsJSON, err := json.Marshal(alg.DefaultSettings)
	if err != nil {
		return fmt.Errorf("marshal settings: %w", err)
	}

	const query = `
UPDATE custom_algorithms SET
    name = $2,
    description = $3,
    kind = $4,
    wasm_binary = $5,
    wasm_hash = $6,
    default_settings = $7,
    updated_at = NOW()
WHERE id = $1`

	_, err = executor.Exec(ctx, query,
		alg.ID, alg.Name, alg.Description, string(alg.Kind),
		alg.WASMBinary, wasmHash, settingsJSON)
	if err != nil {
		return fmt.Errorf("update custom algorithm: %w", err)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id domain.CustomAlgorithmID) error {
	executor := r.getExecutor(ctx)

	const query = `DELETE FROM custom_algorithms WHERE id = $1`

	_, err := executor.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete custom algorithm: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id domain.CustomAlgorithmID) (domain.CustomAlgorithm, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM custom_algorithms WHERE id = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, id)
	if err != nil {
		return domain.CustomAlgorithm{}, fmt.Errorf("query custom algorithm by id: %w", err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[customAlgorithmModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.CustomAlgorithm{}, domain.ErrEntityNotFound
		}
		return domain.CustomAlgorithm{}, fmt.Errorf("collect custom algorithm: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) GetBySlug(ctx context.Context, slug string) (domain.CustomAlgorithm, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM custom_algorithms WHERE slug = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, slug)
	if err != nil {
		return domain.CustomAlgorithm{}, fmt.Errorf("query custom algorithm by slug: %w", err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[customAlgorithmModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.CustomAlgorithm{}, domain.ErrEntityNotFound
		}
		return domain.CustomAlgorithm{}, fmt.Errorf("collect custom algorithm: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) List(ctx context.Context) ([]domain.CustomAlgorithm, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM custom_algorithms ORDER BY name`

	rows, err := executor.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query custom algorithms: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[customAlgorithmModel])
	if err != nil {
		return nil, fmt.Errorf("collect custom algorithms: %w", err)
	}

	result := make([]domain.CustomAlgorithm, 0, len(models))
	for _, m := range models {
		result = append(result, m.toDomain())
	}

	return result, nil
}

func (r *Repository) ListByKind(ctx context.Context, kind domain.AlgorithmKind) ([]domain.CustomAlgorithm, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM custom_algorithms WHERE kind = $1 ORDER BY name`

	rows, err := executor.Query(ctx, query, string(kind))
	if err != nil {
		return nil, fmt.Errorf("query custom algorithms by kind: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[customAlgorithmModel])
	if err != nil {
		return nil, fmt.Errorf("collect custom algorithms: %w", err)
	}

	result := make([]domain.CustomAlgorithm, 0, len(models))
	for _, m := range models {
		result = append(result, m.toDomain())
	}

	return result, nil
}

func (r *Repository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT EXISTS(SELECT 1 FROM custom_algorithms WHERE slug = $1)`

	var exists bool
	err := executor.QueryRow(ctx, query, slug).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check custom algorithm exists: %w", err)
	}

	return exists, nil
}

func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}
	return r.db
}
