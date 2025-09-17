package segments

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/rom8726/etoggle/internal/domain"
)

type segmentModel struct {
	ID          string         `db:"id"`
	ProjectID   string         `db:"project_id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	Conditions  []byte         `db:"conditions"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}

func (m *segmentModel) toDomain() domain.Segment {
	var conditions []domain.Condition
	if err := json.Unmarshal(m.Conditions, &conditions); err != nil {
		slog.Error("unmarshal segment conditions", "conditions", string(m.Conditions), "error", err)
	}

	desc := ""
	if m.Description.Valid {
		desc = m.Description.String
	}

	return domain.Segment{
		ID:          domain.SegmentID(m.ID),
		ProjectID:   domain.ProjectID(m.ProjectID),
		Name:        m.Name,
		Description: desc,
		Conditions:  conditions,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
