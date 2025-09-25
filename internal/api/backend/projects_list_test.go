package apibackend

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
	mockcontract "github.com/togglr-project/togglr/test_mocks/internal_/contract"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRestAPI_ListProjects(t *testing.T) {
	t.Run("successful projects list", func(t *testing.T) {
		mockProjectsUseCase := mockcontract.NewMockProjectsUseCase(t)
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			projectsUseCase:    mockProjectsUseCase,
			permissionsService: mockPermissionsService,
		}

		allProjects := []domain.Project{
			{
				ID:          domain.ProjectID("1"),
				Name:        "Project 1",
				Description: "Description 1",
				CreatedAt:   time.Now(),
			},
			{
				ID:          domain.ProjectID("2"),
				Name:        "Project 2",
				Description: "Description 2",
				CreatedAt:   time.Now(),
			},
		}

		accessibleProjects := []domain.Project{
			allProjects[0],
			allProjects[1],
		}

		mockProjectsUseCase.EXPECT().
			List(mock.Anything).
			Return(allProjects, nil)

		mockPermissionsService.EXPECT().
			GetAccessibleProjects(mock.Anything, allProjects).
			Return(accessibleProjects, nil)

		resp, err := api.ListProjects(context.Background())

		require.NoError(t, err)
		require.NotNil(t, resp)

		listResp, ok := resp.(*generatedapi.ListProjectsResponse)
		require.True(t, ok)
		assert.Len(t, *listResp, 2)

		// Check first project (no team)
		assert.Equal(t, "1", (*listResp)[0].ID)
		assert.Equal(t, "Project 1", (*listResp)[0].Name)
		assert.Equal(t, "Description 1", (*listResp)[0].Description)

		// Check second project (with team)
		assert.Equal(t, "2", (*listResp)[1].ID)
		assert.Equal(t, "Project 2", (*listResp)[1].Name)
		assert.Equal(t, "Description 2", (*listResp)[1].Description)
	})

	t.Run("empty projects list", func(t *testing.T) {
		mockProjectsUseCase := mockcontract.NewMockProjectsUseCase(t)
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			projectsUseCase:    mockProjectsUseCase,
			permissionsService: mockPermissionsService,
		}

		allProjects := []domain.Project{}
		accessibleProjects := []domain.Project{}

		mockProjectsUseCase.EXPECT().
			List(mock.Anything).
			Return(allProjects, nil)

		mockPermissionsService.EXPECT().
			GetAccessibleProjects(mock.Anything, allProjects).
			Return(accessibleProjects, nil)

		resp, err := api.ListProjects(context.Background())

		require.NoError(t, err)
		require.NotNil(t, resp)

		listResp, ok := resp.(*generatedapi.ListProjectsResponse)
		require.True(t, ok)
		assert.Empty(t, *listResp)
	})

	t.Run("get all projects failed", func(t *testing.T) {
		mockProjectsUseCase := mockcontract.NewMockProjectsUseCase(t)

		api := &RestAPI{
			projectsUseCase: mockProjectsUseCase,
		}

		unexpectedErr := errors.New("database error")
		mockProjectsUseCase.EXPECT().
			List(mock.Anything).
			Return(nil, unexpectedErr)

		resp, err := api.ListProjects(context.Background())

		require.Error(t, err)
		require.Nil(t, resp)
		assert.Equal(t, unexpectedErr, err)
	})

	t.Run("filter projects failed", func(t *testing.T) {
		mockProjectsUseCase := mockcontract.NewMockProjectsUseCase(t)
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			projectsUseCase:    mockProjectsUseCase,
			permissionsService: mockPermissionsService,
		}

		allProjects := []domain.Project{
			{

				ID:        domain.ProjectID("1"),
				Name:      "Project 1",
				CreatedAt: time.Now(),
			},
		}

		unexpectedErr := errors.New("permission error")
		mockProjectsUseCase.EXPECT().
			List(mock.Anything).
			Return(allProjects, nil)

		mockPermissionsService.EXPECT().
			GetAccessibleProjects(mock.Anything, allProjects).
			Return(nil, unexpectedErr)

		resp, err := api.ListProjects(context.Background())

		require.Error(t, err)
		require.Nil(t, resp)
		assert.Equal(t, unexpectedErr, err)
	})
}
