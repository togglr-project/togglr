package apibackend

import (
	"context"
	"log/slog"

	appcontext "github.com/togglr-project/togglr/internal/context"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) UpdateLicenseAcceptance(
	ctx context.Context,
	req *generatedapi.UpdateLicenseAcceptanceRequest,
) (generatedapi.UpdateLicenseAcceptanceRes, error) {
	userID := appcontext.UserID(ctx)

	err := r.usersUseCase.UpdateLicenseAcceptance(ctx, userID, req.Accepted)
	if err != nil {
		slog.Error("update license acceptance failed", "error", err, "user_id", userID)

		return nil, err
	}

	return &generatedapi.UpdateLicenseAcceptanceNoContent{}, nil
}
