package rest

import (
	generatedapi "github.com/rom8726/etoggl/internal/generated/server"
)

var _ (generatedapi.Handler) = (*RestAPI)(nil) // TODO: implement!

type RestAPI struct {
}

func New() *RestAPI {
	return &RestAPI{}
}
