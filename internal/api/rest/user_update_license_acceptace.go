package rest

import (
	"context"
	"log/slog"

	etogglcontext "github.com/rom8726/etoggle/internal/context"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) UpdateLicenseAcceptance(
	ctx context.Context,
	req *generatedapi.UpdateLicenseAcceptanceRequest,
) (generatedapi.UpdateLicenseAcceptanceRes, error) {
	userID := etogglcontext.UserID(ctx)

	err := r.usersUseCase.UpdateLicenseAcceptance(ctx, userID, req.Accepted)
	if err != nil {
		slog.Error("update license acceptance failed", "error", err, "user_id", userID)

		return nil, err
	}

	return &generatedapi.UpdateLicenseAcceptanceNoContent{}, nil
}
