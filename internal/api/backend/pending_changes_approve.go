package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// ApprovePendingChange handles POST /api/v1/pending_changes/{pending_change_id}/approve
func (r *RestAPI) ApprovePendingChange(
	ctx context.Context,
	req *generatedapi.ApprovePendingChangeRequest,
	params generatedapi.ApprovePendingChangeParams,
) (generatedapi.ApprovePendingChangeRes, error) {
	pendingChangeID := domain.PendingChangeID(params.PendingChangeID.String())

	// Approve pending change
	err := r.pendingChangesUseCase.Approve(
		ctx,
		pendingChangeID,
		int(req.ApproverUserID),
		req.ApproverName,
		string(req.Auth.Method),
		req.Auth.Credential,
	)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("pending change not found"),
			}}, nil
		}

		if errors.Is(err, domain.ErrPermissionDenied) {
			return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
				Message: generatedapi.NewOptString("permission denied"),
			}}, nil
		}

		// Check for conflict (pending change is not in pending status)
		if err.Error() == "pending change is not in pending status" {
			return &generatedapi.ErrorConflict{Error: generatedapi.ErrorConflictError{
				Message: generatedapi.NewOptString("pending change is not in pending status"),
				Code:    generatedapi.NewOptString("CONFLICT"),
			}}, nil
		}

		slog.Error("approve pending change failed", "error", err)
		return nil, err
	}

	return &generatedapi.SuccessResponse{
		Message: generatedapi.NewOptString("pending change approved successfully"),
	}, nil
}
