package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// InitiateTOTPApproval handles POST /api/v1/pending_changes/{pending_change_id}/initiate-totp
func (r *RestAPI) InitiateTOTPApproval(
	ctx context.Context,
	req *generatedapi.InitiateTOTPApprovalRequest,
	params generatedapi.InitiateTOTPApprovalParams,
) (generatedapi.InitiateTOTPApprovalRes, error) {
	pendingChangeID := domain.PendingChangeID(params.PendingChangeID.String())

	// Initiate TOTP approval session
	sessionID, err := r.pendingChangesUseCase.InitiateTOTPApproval(
		ctx,
		pendingChangeID,
		int(req.ApproverUserID),
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

		slog.Error("initiate TOTP approval failed", "error", err)

		return nil, err
	}

	return &generatedapi.InitiateTOTPApprovalResponse{
		SessionID: sessionID,
		Message:   "TOTP approval session initiated successfully",
	}, nil
}
