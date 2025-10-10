package domain

import (
	"encoding/json"
	"time"
)

type UserNotificationID uint

// UserNotificationType represents the type of user notification.
type UserNotificationType string

const (
	UserNotificationTypeProjectAdded   UserNotificationType = "project_added"
	UserNotificationTypeProjectRemoved UserNotificationType = "project_removed"
	UserNotificationTypeRoleChanged    UserNotificationType = "role_changed"
	UserNotificationTypeNeedApprove    UserNotificationType = "need_approve"
)

type UserNotification struct {
	ID        UserNotificationID
	UserID    UserID
	Type      UserNotificationType
	Content   json.RawMessage
	IsRead    bool
	EmailSent bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserNotificationContent struct {
	UserAddedToProject     *UserAddedToProjectContent     `json:"userAddedToProject"`
	UserRemovedFromProject *UserRemovedFromProjectContent `json:"userRemovedFromProject"`
	UserRoleChanged        *UserRoleChangedContent        `json:"userRoleChanged"`
	NeedApproveChange      *NeedApproveChangeContent      `json:"needApproveChange"`
}

type UserAddedToProjectContent struct {
	ProjectName string `json:"projectName"`
	RoleName    string `json:"roleName"`
	ByUser      string `json:"byUser"`
}

type UserRoleChangedContent struct {
	ProjectName string `json:"projectName"`
	RoleNameOld string `json:"roleNameOld"`
	RoleNameNew string `json:"roleNameNew"`
	ByUser      string `json:"byUser"`
}

type UserRemovedFromProjectContent struct {
	ProjectName string `json:"projectName"`
	ByUser      string `json:"byUser"`
}

type NeedApproveChangeContent struct {
	ProjectName string `json:"projectName"`
	Entity      string `json:"entity"`
	RequestedBy string `json:"requestedBy"`
}
