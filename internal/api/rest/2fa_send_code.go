package rest

import (
	"context"
	"log/slog"

	etogglcontext "github.com/rom8726/etoggl/internal/context"
	generatedapi "github.com/rom8726/etoggl/internal/generated/server"
)

func (r *RestAPI) Send2FACode(ctx context.Context) (generatedapi.Send2FACodeRes, error) {
	err := r.usersUseCase.Send2FACode(ctx, etogglcontext.UserID(ctx), "disable")
	if err != nil {
		slog.Error("send 2fa code failed", "error", err)

		return nil, err
	}

	return &generatedapi.Send2FACodeNoContent{}, nil
}
