package apibackend

import (
	"context"
	"errors"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"

	"github.com/togglr-project/togglr/internal/config"
	"github.com/togglr-project/togglr/internal/contract"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

var _ generatedapi.Handler = (*RestAPI)(nil)

type RestAPI struct {
	config                  *config.Config
	tokenizer               contract.Tokenizer
	usersUseCase            contract.UsersUseCase
	projectsUseCase         contract.ProjectsUseCase
	ldapService             contract.LDAPService
	ldapUseCase             contract.LDAPSyncUseCase
	settingsUseCase         contract.SettingsUseCase
	licenseUseCase          contract.LicenseUseCase
	productInfoUseCase      contract.ProductInfoUseCase
	permissionsService      contract.PermissionsService
	featuresUseCase         contract.FeaturesUseCase
	environmentsUseCase     contract.EnvironmentsUseCase
	flagVariantsUseCase     contract.FlagVariantsUseCase
	rulesUseCase            contract.RulesUseCase
	featureSchedulesUseCase contract.FeatureSchedulesUseCase
	ruleAttributesUseCase   contract.RuleAttributesUseCase
	segmentsUseCase         contract.SegmentsUseCase
	featureProcessor        contract.FeatureProcessor
	categoriesUseCase       contract.CategoriesUseCase
	tagsUseCase             contract.TagsUseCase
	featureTagsUseCase      contract.FeatureTagsUseCase
	pendingChangesUseCase   contract.PendingChangesUseCase
	guardService            contract.GuardService
	guardEngine             contract.GuardEngine
	projectSettingsUseCase  contract.ProjectSettingsUseCase
	dashboardUseCase        contract.DashboardUseCase
	membershipsUseCase      contract.MembershipsUseCase
	auditLogRepo            contract.AuditLogRepository
	errorReportsUseCase     contract.ErrorReportsUseCase
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
	featuresUseCase contract.FeaturesUseCase,
	environmentsUseCase contract.EnvironmentsUseCase,
	flagVariantsUseCase contract.FlagVariantsUseCase,
	rulesUseCase contract.RulesUseCase,
	featureSchedulesUseCase contract.FeatureSchedulesUseCase,
	permissionsService contract.PermissionsService,
	ruleAttributesUseCase contract.RuleAttributesUseCase,
	segmentsUseCase contract.SegmentsUseCase,
	featureProcessor contract.FeatureProcessor,
	categoriesUseCase contract.CategoriesUseCase,
	tagsUseCase contract.TagsUseCase,
	featureTagsUseCase contract.FeatureTagsUseCase,
	pendingChangesUseCase contract.PendingChangesUseCase,
	guardService contract.GuardService,
	guardEngine contract.GuardEngine,
	projectSettingsUseCase contract.ProjectSettingsUseCase,
	dashboardUseCase contract.DashboardUseCase,
	membershipsUseCase contract.MembershipsUseCase,
	auditLogRepo contract.AuditLogRepository,
	errorReportsUseCase contract.ErrorReportsUseCase,
) *RestAPI {
	return &RestAPI{
		config:                  config,
		usersUseCase:            usersService,
		tokenizer:               tokenizer,
		projectsUseCase:         projectsUseCase,
		ldapService:             ldapService,
		ldapUseCase:             ldapUseCase,
		settingsUseCase:         settingsUseCase,
		licenseUseCase:          licenseUseCase,
		productInfoUseCase:      productInfoUseCase,
		permissionsService:      permissionsService,
		featuresUseCase:         featuresUseCase,
		environmentsUseCase:     environmentsUseCase,
		flagVariantsUseCase:     flagVariantsUseCase,
		rulesUseCase:            rulesUseCase,
		featureSchedulesUseCase: featureSchedulesUseCase,
		ruleAttributesUseCase:   ruleAttributesUseCase,
		segmentsUseCase:         segmentsUseCase,
		featureProcessor:        featureProcessor,
		categoriesUseCase:       categoriesUseCase,
		tagsUseCase:             tagsUseCase,
		featureTagsUseCase:      featureTagsUseCase,
		pendingChangesUseCase:   pendingChangesUseCase,
		guardService:            guardService,
		guardEngine:             guardEngine,
		projectSettingsUseCase:  projectSettingsUseCase,
		dashboardUseCase:        dashboardUseCase,
		membershipsUseCase:      membershipsUseCase,
		errorReportsUseCase:     errorReportsUseCase,
		auditLogRepo:            auditLogRepo,
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
