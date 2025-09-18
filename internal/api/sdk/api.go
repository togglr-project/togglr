package apisdk

import (
	"context"
	"errors"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"

	"github.com/rom8726/etoggle/internal/contract"
	generatedapi "github.com/rom8726/etoggle/internal/generated/sdkserver"
)

var _ generatedapi.Handler = (*SDKRestAPI)(nil)

type SDKRestAPI struct {
	featureProcessor contract.FeatureProcessor
}

func New(
	featureProcessor contract.FeatureProcessor,
) *SDKRestAPI {
	return &SDKRestAPI{
		featureProcessor: featureProcessor,
	}
}

func (s *SDKRestAPI) NewError(_ context.Context, err error) *generatedapi.ErrorStatusCode {
	code := http.StatusInternalServerError
	errMessage := err.Error()

	var secError *ogenerrors.SecurityError
	if errors.As(err, &secError) {
		code = http.StatusUnauthorized
		errMessage = "unauthorized"
	}

	return &generatedapi.ErrorStatusCode{
		StatusCode: code,
		Response: generatedapi.Error{
			Error: generatedapi.ErrorError{
				Message: generatedapi.NewOptString(errMessage),
			},
		},
	}
}
