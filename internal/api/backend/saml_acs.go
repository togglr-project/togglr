package apibackend

import (
	"context"
	"errors"
	"log/slog"
	"net/url"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) ConsumeSAMLAssertion(
	ctx context.Context,
	req *generatedapi.ConsumeSAMLAssertionReq,
) (generatedapi.ConsumeSAMLAssertionRes, error) {
	rawReq := appcontext.RawRequest(ctx)

	rawReq.PostForm = make(map[string][]string)
	rawReq.PostForm.Set("SAMLResponse", req.SAMLResponse)
	rawReq.PostForm.Set("RelayState", req.RelayState)

	accessToken, refreshToken, _, err := r.usersUseCase.SSOCallback(
		ctx, domain.SSOProviderNameADSaml, rawReq, req.SAMLResponse, req.RelayState)
	if err != nil {
		slog.Error("SSO assert failed", "error", err)

		if errors.Is(err, domain.ErrInactiveUser) {
			return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
				Message: generatedapi.NewOptString("user is inactive"),
			}}, nil
		}

		return nil, err
	}

	return &generatedapi.ConsumeSAMLAssertionFound{
		Location: generatedapi.NewOptString(r.buildFrontLoginSuccessLocation(accessToken, refreshToken)),
	}, nil
}

func (r *RestAPI) buildFrontLoginSuccessLocation(accessToken, refreshToken string) string {
	values := url.Values{}
	values.Set("access_token", accessToken)
	values.Set("refresh_token", refreshToken)

	return r.config.FrontendURL + "/auth/saml/success?" + values.Encode()
}
