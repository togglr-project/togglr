package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type ErrorReportEvent struct {
	RequestID    string                       `json:"requestID"`
	ProjectID    domain.ProjectID             `json:"projectID"`
	EnvKey       string                       `json:"envKey"`
	FeatureKey   string                       `json:"featureKey"`
	Context      map[domain.RuleAttribute]any `json:"context"`
	ErrorType    string                       `json:"errorType"`
	ErrorMessage string                       `json:"errorMessage"`
}

type EventsBus interface {
	PublishErrorReport(ctx context.Context, event ErrorReportEvent) error
}
