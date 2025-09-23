package apisdk

import (
	"context"
	"time"

	generatedapi "github.com/togglr-project/togglr/internal/generated/sdkserver"
)

func (*SDKRestAPI) SdkV1HealthGet(context.Context) (generatedapi.SdkV1HealthGetRes, error) {
	return &generatedapi.HealthResponse{
		Status:     generatedapi.HealthResponseStatusOk,
		ServerTime: time.Now(),
	}, nil
}
