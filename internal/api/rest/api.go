package rest

import (
	"context"

	generatedapi "github.com/rom8726/etoggl/internal/generated/server"
)

var _ (generatedapi.Handler) = (*RestAPI)(nil)

type RestAPI struct {
}

func New() *RestAPI {
	return &RestAPI{}
}

func (r RestAPI) Ping(ctx context.Context) (generatedapi.PingOK, error) {
	//TODO implement me
	panic("implement me")
}
