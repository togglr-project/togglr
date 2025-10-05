package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

const (
	pendingChangeNotPendingStatus = "pending change is not in pending status"
)

// ApprovePendingChange handles POST /api/v1/pending_changes/{pending_change_id}/approve.
func (r *RestAPI) ApprovePendingChange(
	ctx context.Context,
	req *generatedapi.ApprovePendingChangeRequest,
	params generatedapi.ApprovePendingChangeParams,
) (generatedapi.ApprovePendingChangeRes, error) {
	pendingChangeID := domain.PendingChangeID(params.PendingChangeID.String())

	// Get sessionID from auth if provided
	var sessionID string
	if req.Auth.SessionID.Set && req.Auth.SessionID.Value != "" {
		sessionID = req.Auth.SessionID.Value
	}

	// Approve pending change
	err := r.pendingChangesUseCase.Approve(
		ctx,
		pendingChangeID,
		int(req.ApproverUserID),
		req.ApproverName,
		string(req.Auth.Method),
		req.Auth.Credential,
		sessionID,
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
		if err.Error() == pendingChangeNotPendingStatus {
			return &generatedapi.ErrorConflict{Error: generatedapi.ErrorConflictError{
				Message: generatedapi.NewOptString(pendingChangeNotPendingStatus),
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
