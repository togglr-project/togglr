package contract

import (
	"context"

	"github.com/togglr-project/togglr/internal/domain"
)

type FeedbackEventsRepository interface {
	AddEvent(ctx context.Context, event domain.FeedbackEventDTO) error
	AddEventsBatch(ctx context.Context, events []domain.FeedbackEventDTO) error
}
