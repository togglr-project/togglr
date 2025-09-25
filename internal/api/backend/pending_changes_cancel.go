package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// CancelPendingChange handles POST /api/v1/pending_changes/{pending_change_id}/cancel.
func (r *RestAPI) CancelPendingChange(
	ctx context.Context,
	req *generatedapi.CancelPendingChangeRequest,
	params generatedapi.CancelPendingChangeParams,
) (generatedapi.CancelPendingChangeRes, error) {
	pendingChangeID := domain.PendingChangeID(params.PendingChangeID.String())

	// Cancel pending change
	err := r.pendingChangesUseCase.Cancel(
		ctx,
		pendingChangeID,
		req.CancelledBy,
	)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("pending change not found"),
			}}, nil
		}

		// Check for conflict (pending change is not in pending status)
		if err.Error() == "pending change is not in pending status" {
			return &generatedapi.ErrorConflict{Error: generatedapi.ErrorConflictError{
				Message: generatedapi.NewOptString("pending change is not in pending status"),
				Code:    generatedapi.NewOptString("CONFLICT"),
			}}, nil
		}

		slog.Error("cancel pending change failed", "error", err)

		return nil, err
	}

	return &generatedapi.SuccessResponse{
		Message: generatedapi.NewOptString("pending change cancelled successfully"),
	}, nil
}
