package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// GetPendingChange handles GET /api/v1/pending_changes/{pending_change_id}.
func (r *RestAPI) GetPendingChange(
	ctx context.Context,
	params generatedapi.GetPendingChangeParams,
) (generatedapi.GetPendingChangeRes, error) {
	pendingChangeID := domain.PendingChangeID(params.PendingChangeID.String())

	// Get pending change
	change, err := r.pendingChangesUseCase.GetByID(ctx, pendingChangeID)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("pending change not found"),
			}}, nil
		}

		slog.Error("get pending change failed", "error", err)

		return nil, err
	}

	// Convert to response format
	response := convertPendingChangeToResponse(&change)

	return &response, nil
}
