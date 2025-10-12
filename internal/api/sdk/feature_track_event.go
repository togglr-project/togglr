package apisdk

import (
	"context"

	generatedapi "github.com/togglr-project/togglr/internal/generated/sdkserver"
)

func (s *SDKRestAPI) TrackFeatureEvent(
	ctx context.Context,
	req *generatedapi.TrackRequest,
	params generatedapi.TrackFeatureEventParams,
) (generatedapi.TrackFeatureEventRes, error) {
	// TODO implement me
	panic("implement me")
}
