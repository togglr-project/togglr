package rest

import (
	"context"
	"errors"
	"testing"

	"github.com/rom8726/etoggl/internal/domain"
	generatedapi "github.com/rom8726/etoggl/internal/generated/server"
	mockcontract "github.com/rom8726/etoggl/test_mocks/internal_/contract"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRestAPI_AddProject(t *testing.T) {
	t.Run("successful project creation", func(t *testing.T) {
		mockProjectsUseCase := mockcontract.NewMockProjectsUseCase(t)

		api := &RestAPI{
			projectsUseCase: mockProjectsUseCase,
		}

		req := &generatedapi.AddProjectRequest{
			Name:        "New Project",
			Description: "A new test project",
		}

		mockProjectsUseCase.EXPECT().
			CreateProject(mock.Anything, "New Project", "A new test project").
			Return(domain.Project{}, nil)

		resp, err := api.AddProject(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, resp)

		_, ok := resp.(*generatedapi.AddProjectCreated)
		require.True(t, ok)
	})

	t.Run("successful project creation without team", func(t *testing.T) {
		mockProjectsUseCase := mockcontract.NewMockProjectsUseCase(t)

		api := &RestAPI{
			projectsUseCase: mockProjectsUseCase,
		}

		req := &generatedapi.AddProjectRequest{
			Name:        "New Project",
			Description: "A new test project",
		}

		mockProjectsUseCase.EXPECT().
			CreateProject(mock.Anything, "New Project", "A new test project").
			Return(domain.Project{}, nil)

		resp, err := api.AddProject(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, resp)

		_, ok := resp.(*generatedapi.AddProjectCreated)
		require.True(t, ok)
	})

	t.Run("permission denied", func(t *testing.T) {
		mockProjectsUseCase := mockcontract.NewMockProjectsUseCase(t)

		api := &RestAPI{
			projectsUseCase: mockProjectsUseCase,
		}

		req := &generatedapi.AddProjectRequest{
			Name:        "New Project",
			Description: "A new test project",
		}

		mockProjectsUseCase.EXPECT().
			CreateProject(mock.Anything, "New Project", "A new test project").
			Return(domain.Project{}, domain.ErrPermissionDenied)

		resp, err := api.AddProject(context.Background(), req)

		require.Error(t, err)
		require.Nil(t, resp)
		assert.Equal(t, domain.ErrPermissionDenied, err)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockProjectsUseCase := mockcontract.NewMockProjectsUseCase(t)

		api := &RestAPI{
			projectsUseCase: mockProjectsUseCase,
		}

		req := &generatedapi.AddProjectRequest{
			Name:        "New Project",
			Description: "A new test project",
		}

		mockProjectsUseCase.EXPECT().
			CreateProject(mock.Anything, "New Project", "A new test project").
			Return(domain.Project{}, domain.ErrUserNotFound)

		resp, err := api.AddProject(context.Background(), req)

		require.Error(t, err)
		require.Nil(t, resp)
		assert.Equal(t, domain.ErrUserNotFound, err)
	})

	t.Run("project name already exists", func(t *testing.T) {
		mockProjectsUseCase := mockcontract.NewMockProjectsUseCase(t)

		api := &RestAPI{
			projectsUseCase: mockProjectsUseCase,
		}

		req := &generatedapi.AddProjectRequest{
			Name:        "Existing Project",
			Description: "A project with existing name",
		}

		mockProjectsUseCase.EXPECT().
			CreateProject(mock.Anything, "Existing Project", "A project with existing name").
			Return(domain.Project{}, domain.ErrEntityNotFound)

		resp, err := api.AddProject(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("permission check failed with unexpected error", func(t *testing.T) {
		mockProjectsUseCase := mockcontract.NewMockProjectsUseCase(t)

		api := &RestAPI{
			projectsUseCase: mockProjectsUseCase,
		}

		req := &generatedapi.AddProjectRequest{
			Name:        "New Project",
			Description: "A new test project",
		}

		unexpectedErr := errors.New("database error")
		mockProjectsUseCase.EXPECT().
			CreateProject(mock.Anything, "New Project", "A new test project").
			Return(domain.Project{}, unexpectedErr)

		resp, err := api.AddProject(context.Background(), req)

		require.Error(t, err)
		require.Nil(t, resp)
		assert.Equal(t, unexpectedErr, err)
	})

	t.Run("create project failed with unexpected error", func(t *testing.T) {
		mockProjectsUseCase := mockcontract.NewMockProjectsUseCase(t)

		api := &RestAPI{
			projectsUseCase: mockProjectsUseCase,
		}

		req := &generatedapi.AddProjectRequest{
			Name:        "New Project",
			Description: "A new test project",
		}

		unexpectedErr := errors.New("database error")
		mockProjectsUseCase.EXPECT().
			CreateProject(mock.Anything, "New Project", "A new test project").
			Return(domain.Project{}, unexpectedErr)

		resp, err := api.AddProject(context.Background(), req)

		require.Error(t, err)
		require.Nil(t, resp)
		assert.Equal(t, unexpectedErr, err)
	})
}
