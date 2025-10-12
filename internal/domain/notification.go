package domain

import (
	"encoding/json"
	"time"
)

// NotificationType represents the type of notification.
type NotificationType string

const (
	NotificationTypeEmail      NotificationType = "email"
	NotificationTypeTelegram   NotificationType = "telegram"
	NotificationTypeSlack      NotificationType = "slack"
	NotificationTypeMattermost NotificationType = "mattermost"
	NotificationTypeWebhook    NotificationType = "webhook"
	NotificationTypePachca     NotificationType = "pachca"
)

type NotificationStatus string

const (
	NotificationStatusPending NotificationStatus = "pending"
	NotificationStatusSent    NotificationStatus = "sent"
	NotificationStatusFailed  NotificationStatus = "failed"
)

type NotificationSettingID uint

// NotificationSetting represents a notification setting for a project.
type NotificationSetting struct {
	ID            NotificationSettingID
	ProjectID     ProjectID
	EnvironmentID EnvironmentID
	Type          NotificationType
	Config        json.RawMessage
	Enabled       bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type NotificationSettingDTO struct {
	ProjectID     ProjectID
	EnvironmentID EnvironmentID
	Type          NotificationType
	Config        json.RawMessage
	Enabled       bool
}

type FeatureNotificationID uint
type FeatureNotification struct {
	ID            FeatureNotificationID
	ProjectID     ProjectID
	EnvironmentID EnvironmentID
	FeatureID     FeatureID
	Payload       FeatureNotificationPayload
	SentAt        *time.Time
	Status        NotificationStatus
	FailReason    *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type FeatureNotificationWithSettings struct {
	FeatureNotification
	Settings []NotificationSetting
}

type FeatureNotificationPayload struct {
	State         *FeatureNotificationStatePayload         `json:"state,omitempty"`
	AutoDisabled  *FeatureNotificationAutoDisabledPayload  `json:"autoDisabled,omitempty"`
	ChangeRequest *FeatureNotificationChangeRequestPayload `json:"changeRequest,omitempty"`
}

type FeatureNotificationStatePayload struct {
	Enabled   bool   `json:"enabled"`
	ChangedBy string `json:"changedBy"`
}

type FeatureNotificationAutoDisabledPayload struct {
	DisabledAt time.Time `json:"disabledAt"`
}

type FeatureNotificationChangeRequestPayload struct {
	RequestedBy string `json:"requestedBy"`
}
