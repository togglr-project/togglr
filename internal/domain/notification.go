package domain

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
