package apibackend

import (
	"context"
	"errors"
	"testing"

	appctx "github.com/rom8726/etoggle/internal/context"
	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
	mockcontract "github.com/rom8726/etoggle/test_mocks/internal_/contract"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRestAPI_AddProject(t *testing.T) {
	t.Run("successful project creation", func(t *testing.T) {
		mockProjectsUseCase := mockcontract.NewMockProjectsUseCase(t)
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			projectsUseCase:    mockProjectsUseCase,
			permissionsService: mockPermissionsService,
		}

		req := &generatedapi.AddProjectRequest{
			Name:        "New Project",
			Description: "A new test project",
		}

		mockPermissionsService.EXPECT().
			HasGlobalPermission(mock.Anything, domain.PermProjectCreate).
			Return(true, nil)
		mockProjectsUseCase.EXPECT().
			CreateProject(mock.Anything, "New Project", "A new test project").
			Return(domain.Project{}, nil)

		userID := domain.UserID(1)
		ctx := context.Background()
		ctx = appctx.WithUserID(ctx, userID)

		resp, err := api.AddProject(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)

		_, ok := resp.(*generatedapi.AddProjectCreated)
		require.True(t, ok)
	})

	t.Run("permission denied", func(t *testing.T) {
		mockProjectsUseCase := mockcontract.NewMockProjectsUseCase(t)
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			projectsUseCase:    mockProjectsUseCase,
			permissionsService: mockPermissionsService,
		}

		req := &generatedapi.AddProjectRequest{
			Name:        "New Project",
			Description: "A new test project",
		}

		mockPermissionsService.EXPECT().
			HasGlobalPermission(mock.Anything, domain.PermProjectCreate).
			Return(false, nil)

		userID := domain.UserID(1)
		ctx := context.Background()
		ctx = appctx.WithUserID(ctx, userID)

		resp, err := api.AddProject(ctx, req)

		require.Nil(t, err)
		require.NotNil(t, resp)

		_, ok := resp.(*generatedapi.ErrorPermissionDenied)
		require.True(t, ok)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockProjectsUseCase := mockcontract.NewMockProjectsUseCase(t)
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			projectsUseCase:    mockProjectsUseCase,
			permissionsService: mockPermissionsService,
		}

		req := &generatedapi.AddProjectRequest{
			Name:        "New Project",
			Description: "A new test project",
		}

		ctx := context.Background()

		resp, err := api.AddProject(ctx, req)

		require.Nil(t, err)
		require.NotNil(t, resp)

		_, ok := resp.(*generatedapi.ErrorUnauthorized)
		require.True(t, ok)
	})

	t.Run("project name already exists", func(t *testing.T) {
		mockProjectsUseCase := mockcontract.NewMockProjectsUseCase(t)
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			projectsUseCase:    mockProjectsUseCase,
			permissionsService: mockPermissionsService,
		}

		req := &generatedapi.AddProjectRequest{
			Name:        "Existing Project",
			Description: "A project with existing name",
		}

		mockPermissionsService.EXPECT().
			HasGlobalPermission(mock.Anything, domain.PermProjectCreate).
			Return(true, nil)
		mockProjectsUseCase.EXPECT().
			CreateProject(mock.Anything, "Existing Project", "A project with existing name").
			Return(domain.Project{}, domain.ErrEntityAlreadyExists)

		userID := domain.UserID(1)
		ctx := context.Background()
		ctx = appctx.WithUserID(ctx, userID)

		resp, err := api.AddProject(ctx, req)

		require.Nil(t, err)
		require.NotNil(t, resp)

		_, ok := resp.(*generatedapi.Error)
		require.True(t, ok)
	})

	t.Run("permission check failed with unexpected error", func(t *testing.T) {
		mockProjectsUseCase := mockcontract.NewMockProjectsUseCase(t)
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			projectsUseCase:    mockProjectsUseCase,
			permissionsService: mockPermissionsService,
		}

		req := &generatedapi.AddProjectRequest{
			Name:        "New Project",
			Description: "A new test project",
		}

		unexpectedErr := errors.New("database error")
		mockPermissionsService.EXPECT().
			HasGlobalPermission(mock.Anything, domain.PermProjectCreate).
			Return(true, nil)
		mockProjectsUseCase.EXPECT().
			CreateProject(mock.Anything, "New Project", "A new test project").
			Return(domain.Project{}, unexpectedErr)

		userID := domain.UserID(1)
		ctx := context.Background()
		ctx = appctx.WithUserID(ctx, userID)
		resp, err := api.AddProject(ctx, req)

		require.Error(t, err)
		require.Nil(t, resp)
		assert.Equal(t, unexpectedErr, err)
	})
}
