package apisdk

import (
	"context"
	"time"

	generatedapi "github.com/rom8726/etoggle/internal/generated/sdkserver"
)

func (*SDKRestAPI) SdkV1HealthGet(context.Context) (generatedapi.SdkV1HealthGetRes, error) {
	return &generatedapi.HealthResponse{
		Status:     generatedapi.HealthResponseStatusOk,
		ServerTime: time.Now(),
	}, nil
}
