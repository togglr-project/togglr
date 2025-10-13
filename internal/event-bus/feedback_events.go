package event_bus

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
)

func (s *Service) PublishFeedbackEvent(ctx context.Context, event domain.FeedbackEventDTO) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal feedback event: %w", err)
	}

	return s.bus.Publish(ctx, TopicSDKFeedbackEvents, data)
}

func (s *Service) processSDKFeedbackEvents(ctx context.Context, messages [][]byte) error {
	batch := make([]domain.FeedbackEventDTO, 0, len(messages))

	for _, msg := range messages {
		var event domain.FeedbackEventDTO
		if err := json.Unmarshal(msg, &event); err != nil {
			slog.Error("failed to unmarshal feedback event", "err", err)

			continue
		}

		batch = append(batch, event)
	}

	return s.feedbackEventsRepo.AddEventsBatch(ctx, batch)
}
