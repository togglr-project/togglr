package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// RejectPendingChange handles POST /api/v1/pending_changes/{pending_change_id}/reject
func (r *RestAPI) RejectPendingChange(
	ctx context.Context,
	req *generatedapi.RejectPendingChangeRequest,
	params generatedapi.RejectPendingChangeParams,
) (generatedapi.RejectPendingChangeRes, error) {
	pendingChangeID := domain.PendingChangeID(params.PendingChangeID.String())

	// Reject pending change
	err := r.pendingChangesUseCase.Reject(
		ctx,
		pendingChangeID,
		req.RejectedBy,
		req.Reason,
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

		slog.Error("reject pending change failed", "error", err)
		return nil, err
	}

	return &generatedapi.SuccessResponse{
		Message: generatedapi.NewOptString("pending change rejected successfully"),
	}, nil
}
