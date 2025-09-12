//nolint:interfacebloat // it's ok here
package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"

	"github.com/rom8726/etoggl/internal/config"
	"github.com/rom8726/etoggl/internal/contract"
	generatedapi "github.com/rom8726/etoggl/internal/generated/server"
)

var _ generatedapi.Handler = (*RestAPI)(nil)

type RestAPI struct {
	config             *config.Config
	tokenizer          contract.Tokenizer
	usersUseCase       contract.UsersUseCase
	projectsUseCase    contract.ProjectsUseCase
	ldapService        contract.LDAPService
	ldapUseCase        contract.LDAPSyncUseCase
	settingsUseCase    contract.SettingsUseCase
	licenseUseCase     contract.LicenseUseCase
	productInfoUseCase contract.ProductInfoUseCase
	permissionsService contract.PermissionsService
}

func New(
	config *config.Config,
	usersService contract.UsersUseCase,
	tokenizer contract.Tokenizer,
	projectsUseCase contract.ProjectsUseCase,
	ldapService contract.LDAPService,
	ldapUseCase contract.LDAPSyncUseCase,
	settingsUseCase contract.SettingsUseCase,
	licenseUseCase contract.LicenseUseCase,
	productInfoUseCase contract.ProductInfoUseCase,
) *RestAPI {
	return &RestAPI{
		config:             config,
		usersUseCase:       usersService,
		tokenizer:          tokenizer,
		projectsUseCase:    projectsUseCase,
		ldapService:        ldapService,
		ldapUseCase:        ldapUseCase,
		settingsUseCase:    settingsUseCase,
		licenseUseCase:     licenseUseCase,
		productInfoUseCase: productInfoUseCase,
	}
}

func (r *RestAPI) NewError(_ context.Context, err error) *generatedapi.ErrorStatusCode {
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
