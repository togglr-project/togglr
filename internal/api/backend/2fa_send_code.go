package apibackend

import (
	"context"
	"log/slog"

	appcontext "github.com/togglr-project/togglr/internal/context"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) Send2FACode(ctx context.Context) (generatedapi.Send2FACodeRes, error) {
	err := r.usersUseCase.Send2FACode(ctx, appcontext.UserID(ctx), "disable")
	if err != nil {
		slog.Error("send 2fa code failed", "error", err)

		return nil, err
	}

	return &generatedapi.Send2FACodeNoContent{}, nil
}
