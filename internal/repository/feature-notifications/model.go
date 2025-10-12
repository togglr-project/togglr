package feature_notifications

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/togglr-project/togglr/internal/domain"
)

type notificationModel struct {
	ID            uint            `db:"id"`
	ProjectID     string          `db:"project_id"`
	EnvironmentID int64           `db:"environment_id"`
	FeatureID     string          `db:"feature_id"`
	Payload       json.RawMessage `db:"payload"`
	SentAt        *time.Time      `db:"sent_at"`
	Status        string          `db:"status"`
	FailReason    *string         `db:"fail_reason"`
	CreatedAt     time.Time       `db:"created_at"`
	UpdatedAt     time.Time       `db:"updated_at"`
}

func (m *notificationModel) toDomain() domain.FeatureNotification {
	var payload domain.FeatureNotificationPayload
	if err := json.Unmarshal(m.Payload, &payload); err != nil {
		slog.Error("unmarshal notification payload", "payload", string(m.Payload), "error", err)
	}

	return domain.FeatureNotification{
		ID:            domain.FeatureNotificationID(m.ID),
		ProjectID:     domain.ProjectID(m.ProjectID),
		EnvironmentID: domain.EnvironmentID(m.EnvironmentID),
		FeatureID:     domain.FeatureID(m.FeatureID),
		Payload:       payload,
		SentAt:        m.SentAt,
		Status:        domain.NotificationStatus(m.Status),
		FailReason:    m.FailReason,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}
