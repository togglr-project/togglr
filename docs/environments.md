# План перехода к модели с окружениями

## Обзор изменений

Миграция `029_environments.up.sql` вводит концепцию окружений (environments) в систему Togglr. Основные изменения:

### Структурные изменения в БД

1. **Новая таблица `environments`**:
   - `id` (BIGSERIAL PRIMARY KEY)
   - `project_id` (UUID, ссылка на projects)
   - `key` (VARCHAR(20)) - ключ окружения (dev, stage, prod)
   - `name` (VARCHAR(50)) - человекочитаемое имя
   - `api_key` (UUID) - API ключ для окружения
   - Уникальность: `(project_id, key)`

2. **Новая таблица `feature_params`**:
   - `feature_id` (UUID, ссылка на features)
   - `environment_id` (BIGINT, ссылка на environments)
   - `enabled` (BOOLEAN) - состояние фичи в окружении
   - `default_value` (VARCHAR(128)) - значение по умолчанию
   - Первичный ключ: `(feature_id, environment_id)`

3. **Изменения в существующих таблицах**:
   - `projects`: удален столбец `api_key` (перенесен в environments)
   - `features`: удалены столбцы `enabled` и `default_variant` (перенесены в feature_params)
   - `rules`: добавлен `environment_id` (NOT NULL, default 0)
   - `flag_variants`: добавлен `environment_id` (NOT NULL, default 0)
   - `feature_schedules`: добавлен `environment_id` (NOT NULL, default 0)
   - `audit_log`: добавлен `environment_id` (NOT NULL)
   - `pending_changes`: добавлен `environment_id` (NOT NULL, default 0)

## План миграции

### Этап 1: Обновление доменных моделей

#### 1.1 Создание новых доменных моделей

**Файл: `internal/domain/environment.go`**
```go
package domain

import "time"

type EnvironmentID int64

type Environment struct {
    ID        EnvironmentID
    ProjectID ProjectID
    Key       string    // dev, stage, prod
    Name      string    // Development, Staging, Production
    APIKey    string
    CreatedAt time.Time
}

type EnvironmentDTO struct {
    Key  string
    Name string
}

func (id EnvironmentID) String() string {
    return fmt.Sprintf("%d", id)
}
```

**Файл: `internal/domain/feature_params.go`**
```go
package domain

import "time"

type FeatureParams struct {
    FeatureID     FeatureID
    EnvironmentID EnvironmentID
    Enabled       bool
    DefaultValue  *string
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

type FeatureParamsDTO struct {
    EnvironmentID EnvironmentID
    Enabled       bool
    DefaultValue  *string
}
```

#### 1.2 Обновление существующих моделей

**Обновить `internal/domain/feature.go`**:
- Удалить поля `Enabled` и `DefaultVariant` из структуры `Feature`
- Добавить метод для получения параметров фичи по окружению

**Обновить `internal/domain/rule.go`**:
- Добавить поле `EnvironmentID EnvironmentID`

**Обновить `internal/domain/flag_variants.go`**:
- Добавить поле `EnvironmentID EnvironmentID`

**Обновить `internal/domain/feature_schedule.go`**:
- Добавить поле `EnvironmentID EnvironmentID`

**Обновить `internal/domain/project.go`**:
- Удалить поле `APIKey` из структуры `Project`
- Добавить метод для получения API ключа по окружению

### Этап 2: Создание репозиториев

#### 2.1 Репозиторий окружений

**Файл: `internal/repository/environments/repository.go`**
```go
package environments

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/togglr-project/togglr/internal/domain"
    "github.com/togglr-project/togglr/pkg/db"
)

type Repository struct {
    db db.Tx
}

func New(pool *pgxpool.Pool) *Repository {
    return &Repository{db: pool}
}

// Create создает новое окружение
func (r *Repository) Create(ctx context.Context, env domain.Environment) (domain.Environment, error)

// GetByID получает окружение по ID
func (r *Repository) GetByID(ctx context.Context, id domain.EnvironmentID) (domain.Environment, error)

// GetByProjectIDAndKey получает окружение по project_id и key
func (r *Repository) GetByProjectIDAndKey(ctx context.Context, projectID domain.ProjectID, key string) (domain.Environment, error)

// ListByProjectID получает все окружения проекта
func (r *Repository) ListByProjectID(ctx context.Context, projectID domain.ProjectID) ([]domain.Environment, error)

// GetByAPIKey получает окружение по API ключу
func (r *Repository) GetByAPIKey(ctx context.Context, apiKey string) (domain.Environment, error)

// Update обновляет окружение
func (r *Repository) Update(ctx context.Context, env domain.Environment) (domain.Environment, error)

// Delete удаляет окружение
func (r *Repository) Delete(ctx context.Context, id domain.EnvironmentID) error
```

#### 2.2 Репозиторий параметров фич

**Файл: `internal/repository/feature_params/repository.go`**
```go
package feature_params

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/togglr-project/togglr/internal/domain"
    "github.com/togglr-project/togglr/pkg/db"
)

type Repository struct {
    db db.Tx
}

func New(pool *pgxpool.Pool) *Repository {
    return &Repository{db: pool}
}

// Create создает параметры фичи для окружения
func (r *Repository) Create(ctx context.Context, params domain.FeatureParams) (domain.FeatureParams, error)

// GetByFeatureAndEnvironment получает параметры фичи для окружения
func (r *Repository) GetByFeatureAndEnvironment(ctx context.Context, featureID domain.FeatureID, envID domain.EnvironmentID) (domain.FeatureParams, error)

// ListByFeatureID получает все параметры фичи
func (r *Repository) ListByFeatureID(ctx context.Context, featureID domain.FeatureID) ([]domain.FeatureParams, error)

// ListByEnvironmentID получает все параметры окружения
func (r *Repository) ListByEnvironmentID(ctx context.Context, envID domain.EnvironmentID) ([]domain.FeatureParams, error)

// Update обновляет параметры фичи
func (r *Repository) Update(ctx context.Context, params domain.FeatureParams) (domain.FeatureParams, error)

// Delete удаляет параметры фичи
func (r *Repository) Delete(ctx context.Context, featureID domain.FeatureID, envID domain.EnvironmentID) error
```

#### 2.3 Обновление существующих репозиториев

**Обновить `internal/repository/features/repository.go`**:
- Удалить поля `enabled` и `default_variant` из запросов
- Добавить методы для работы с параметрами фич по окружениям

**Обновить `internal/repository/rules/repository.go`**:
- Добавить `environment_id` во все запросы
- Обновить методы фильтрации по окружению

**Обновить `internal/repository/flagvariants/repository.go`**:
- Добавить `environment_id` во все запросы
- Обновить методы фильтрации по окружению

**Обновить `internal/repository/featureschedules/repository.go`**:
- Добавить `environment_id` во все запросы
- Обновить методы фильтрации по окружению

**Обновить `internal/repository/projects/repository.go`**:
- Удалить поле `api_key` из запросов
- Добавить методы для работы с окружениями проекта

### Этап 3: Обновление контрактов API

#### 3.1 Новые контракты

**Файл: `internal/contract/environments.go`**
```go
package contract

import "github.com/togglr-project/togglr/internal/domain"

// CreateEnvironmentRequest запрос на создание окружения
type CreateEnvironmentRequest struct {
    Key  string `json:"key" validate:"required,min=1,max=20"`
    Name string `json:"name" validate:"required,min=1,max=50"`
}

// UpdateEnvironmentRequest запрос на обновление окружения
type UpdateEnvironmentRequest struct {
    Name string `json:"name" validate:"required,min=1,max=50"`
}

// EnvironmentResponse ответ с данными окружения
type EnvironmentResponse struct {
    ID        domain.EnvironmentID `json:"id"`
    ProjectID domain.ProjectID     `json:"project_id"`
    Key       string               `json:"key"`
    Name      string               `json:"name"`
    APIKey    string               `json:"api_key"`
    CreatedAt string               `json:"created_at"`
}

// ListEnvironmentsResponse ответ со списком окружений
type ListEnvironmentsResponse struct {
    Items      []EnvironmentResponse `json:"items"`
    Pagination PaginationResponse    `json:"pagination"`
}
```

**Файл: `internal/contract/feature_params.go`**
```go
package contract

import "github.com/togglr-project/togglr/internal/domain"

// UpdateFeatureParamsRequest запрос на обновление параметров фичи
type UpdateFeatureParamsRequest struct {
    EnvironmentID domain.EnvironmentID `json:"environment_id" validate:"required"`
    Enabled       bool                 `json:"enabled"`
    DefaultValue  *string              `json:"default_value,omitempty"`
}

// FeatureParamsResponse ответ с параметрами фичи
type FeatureParamsResponse struct {
    FeatureID     domain.FeatureID     `json:"feature_id"`
    EnvironmentID domain.EnvironmentID `json:"environment_id"`
    Enabled       bool                 `json:"enabled"`
    DefaultValue  *string              `json:"default_value,omitempty"`
    CreatedAt     string               `json:"created_at"`
    UpdatedAt     string               `json:"updated_at"`
}

// ListFeatureParamsResponse ответ со списком параметров фичи
type ListFeatureParamsResponse struct {
    Items []FeatureParamsResponse `json:"items"`
}
```

#### 3.2 Обновление существующих контрактов

**Обновить `internal/contract/features.go`**:
- Удалить поля `enabled` и `default_variant` из структур
- Добавить поля для работы с окружениями

**Обновить `internal/contract/projects.go`**:
- Удалить поле `api_key` из структур
- Добавить методы для работы с окружениями

### Этап 4: Обновление use cases

#### 4.1 Новые use cases

**Файл: `internal/usecases/environments.go`**
```go
package usecases

import (
    "context"
    "github.com/togglr-project/togglr/internal/contract"
    "github.com/togglr-project/togglr/internal/domain"
    "github.com/togglr-project/togglr/internal/repository/environments"
)

type EnvironmentUseCase struct {
    envRepo *environments.Repository
}

func NewEnvironmentUseCase(envRepo *environments.Repository) *EnvironmentUseCase {
    return &EnvironmentUseCase{envRepo: envRepo}
}

// CreateEnvironment создает новое окружение
func (uc *EnvironmentUseCase) CreateEnvironment(ctx context.Context, projectID domain.ProjectID, req contract.CreateEnvironmentRequest) (contract.EnvironmentResponse, error)

// ListEnvironments получает список окружений проекта
func (uc *EnvironmentUseCase) ListEnvironments(ctx context.Context, projectID domain.ProjectID) ([]contract.EnvironmentResponse, error)

// GetEnvironment получает окружение по ID
func (uc *EnvironmentUseCase) GetEnvironment(ctx context.Context, id domain.EnvironmentID) (contract.EnvironmentResponse, error)

// UpdateEnvironment обновляет окружение
func (uc *EnvironmentUseCase) UpdateEnvironment(ctx context.Context, id domain.EnvironmentID, req contract.UpdateEnvironmentRequest) (contract.EnvironmentResponse, error)

// DeleteEnvironment удаляет окружение
func (uc *EnvironmentUseCase) DeleteEnvironment(ctx context.Context, id domain.EnvironmentID) error

// GetEnvironmentByAPIKey получает окружение по API ключу
func (uc *EnvironmentUseCase) GetEnvironmentByAPIKey(ctx context.Context, apiKey string) (contract.EnvironmentResponse, error)
```

**Файл: `internal/usecases/feature_params.go`**
```go
package usecases

import (
    "context"
    "github.com/togglr-project/togglr/internal/contract"
    "github.com/togglr-project/togglr/internal/domain"
    "github.com/togglr-project/togglr/internal/repository/feature_params"
)

type FeatureParamsUseCase struct {
    paramsRepo *feature_params.Repository
}

func NewFeatureParamsUseCase(paramsRepo *feature_params.Repository) *FeatureParamsUseCase {
    return &FeatureParamsUseCase{paramsRepo: paramsRepo}
}

// UpdateFeatureParams обновляет параметры фичи для окружения
func (uc *FeatureParamsUseCase) UpdateFeatureParams(ctx context.Context, featureID domain.FeatureID, req contract.UpdateFeatureParamsRequest) (contract.FeatureParamsResponse, error)

// GetFeatureParams получает параметры фичи для окружения
func (uc *FeatureParamsUseCase) GetFeatureParams(ctx context.Context, featureID domain.FeatureID, envID domain.EnvironmentID) (contract.FeatureParamsResponse, error)

// ListFeatureParams получает все параметры фичи
func (uc *FeatureParamsUseCase) ListFeatureParams(ctx context.Context, featureID domain.FeatureID) ([]contract.FeatureParamsResponse, error)
```

#### 4.2 Обновление существующих use cases

**Обновить `internal/usecases/features.go`**:
- Добавить методы для работы с параметрами фич по окружениям
- Обновить методы создания/обновления фич

**Обновить `internal/usecases/projects.go`**:
- Добавить методы для работы с окружениями проекта
- Обновить методы получения API ключей

### Этап 5: Обновление API handlers

#### 5.1 Новые handlers

**Файл: `internal/api/backend/environments.go`**
```go
package backend

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/togglr-project/togglr/internal/contract"
    "github.com/togglr-project/togglr/internal/domain"
    "github.com/togglr-project/togglr/internal/usecases"
)

type EnvironmentHandler struct {
    envUC *usecases.EnvironmentUseCase
}

func NewEnvironmentHandler(envUC *usecases.EnvironmentUseCase) *EnvironmentHandler {
    return &EnvironmentHandler{envUC: envUC}
}

// CreateEnvironment создает новое окружение
func (h *EnvironmentHandler) CreateEnvironment(c *gin.Context)

// ListEnvironments получает список окружений проекта
func (h *EnvironmentHandler) ListEnvironments(c *gin.Context)

// GetEnvironment получает окружение по ID
func (h *EnvironmentHandler) GetEnvironment(c *gin.Context)

// UpdateEnvironment обновляет окружение
func (h *EnvironmentHandler) UpdateEnvironment(c *gin.Context)

// DeleteEnvironment удаляет окружение
func (h *EnvironmentHandler) DeleteEnvironment(c *gin.Context)
```

**Файл: `internal/api/backend/feature_params.go`**
```go
package backend

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/togglr-project/togglr/internal/contract"
    "github.com/togglr-project/togglr/internal/domain"
    "github.com/togglr-project/togglr/internal/usecases"
)

type FeatureParamsHandler struct {
    paramsUC *usecases.FeatureParamsUseCase
}

func NewFeatureParamsHandler(paramsUC *usecases.FeatureParamsUseCase) *FeatureParamsHandler {
    return &FeatureParamsHandler{paramsUC: paramsUC}
}

// UpdateFeatureParams обновляет параметры фичи
func (h *FeatureParamsHandler) UpdateFeatureParams(c *gin.Context)

// GetFeatureParams получает параметры фичи
func (h *FeatureParamsHandler) GetFeatureParams(c *gin.Context)

// ListFeatureParams получает все параметры фичи
func (h *FeatureParamsHandler) ListFeatureParams(c *gin.Context)
```

#### 5.2 Обновление существующих handlers

**Обновить `internal/api/backend/features.go`**:
- Добавить поддержку окружений в методы работы с фичами
- Обновить методы создания/обновления фич

**Обновить `internal/api/backend/projects.go`**:
- Добавить методы для работы с окружениями проекта
- Обновить методы получения API ключей

### Этап 6: Обновление OpenAPI схемы

#### 6.1 Новые endpoints

```yaml
# Окружения
/projects/{project_id}/environments:
  get:
    summary: List project environments
    parameters:
      - name: project_id
        in: path
        required: true
        schema:
          type: string
          format: uuid
    responses:
      200:
        description: List of environments
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ListEnvironmentsResponse'

  post:
    summary: Create environment
    parameters:
      - name: project_id
        in: path
        required: true
        schema:
          type: string
          format: uuid
    requestBody:
      required: true
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/CreateEnvironmentRequest'
    responses:
      201:
        description: Environment created
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/EnvironmentResponse'

/environments/{environment_id}:
  get:
    summary: Get environment
    parameters:
      - name: environment_id
        in: path
        required: true
        schema:
          type: integer
    responses:
      200:
        description: Environment details
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/EnvironmentResponse'

  put:
    summary: Update environment
    parameters:
      - name: environment_id
        in: path
        required: true
        schema:
          type: integer
    requestBody:
      required: true
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/UpdateEnvironmentRequest'
    responses:
      200:
        description: Environment updated
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/EnvironmentResponse'

  delete:
    summary: Delete environment
    parameters:
      - name: environment_id
        in: path
        required: true
        schema:
          type: integer
    responses:
      204:
        description: Environment deleted

# Параметры фич
/features/{feature_id}/params:
  get:
    summary: List feature parameters
    parameters:
      - name: feature_id
        in: path
        required: true
        schema:
          type: string
          format: uuid
    responses:
      200:
        description: List of feature parameters
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ListFeatureParamsResponse'

  post:
    summary: Update feature parameters
    parameters:
      - name: feature_id
        in: path
        required: true
        schema:
          type: string
          format: uuid
    requestBody:
      required: true
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/UpdateFeatureParamsRequest'
    responses:
      200:
        description: Feature parameters updated
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/FeatureParamsResponse'
```

#### 6.2 Новые схемы

```yaml
components:
  schemas:
    Environment:
      type: object
      properties:
        id:
          type: integer
          format: int64
        project_id:
          type: string
          format: uuid
        key:
          type: string
          maxLength: 20
        name:
          type: string
          maxLength: 50
        api_key:
          type: string
          format: uuid
        created_at:
          type: string
          format: date-time
      required:
        - id
        - project_id
        - key
        - name
        - api_key
        - created_at

    CreateEnvironmentRequest:
      type: object
      properties:
        key:
          type: string
          maxLength: 20
        name:
          type: string
          maxLength: 50
      required:
        - key
        - name

    UpdateEnvironmentRequest:
      type: object
      properties:
        name:
          type: string
          maxLength: 50
      required:
        - name

    FeatureParams:
      type: object
      properties:
        feature_id:
          type: string
          format: uuid
        environment_id:
          type: integer
          format: int64
        enabled:
          type: boolean
        default_value:
          type: string
          maxLength: 128
        created_at:
          type: string
          format: date-time
        updated_at:
          type: string
          format: date-time
      required:
        - feature_id
        - environment_id
        - enabled
        - created_at
        - updated_at

    UpdateFeatureParamsRequest:
      type: object
      properties:
        environment_id:
          type: integer
          format: int64
        enabled:
          type: boolean
        default_value:
          type: string
          maxLength: 128
      required:
        - environment_id
        - enabled
```

### Этап 7: Обновление фронтенда

#### 7.1 Новые компоненты

**Файл: `webui/src/components/environments/EnvironmentCard.tsx`**
```tsx
import React from 'react';
import { Card, CardContent, Typography, Box, Chip, IconButton, Tooltip } from '@mui/material';
import { Edit as EditIcon, Delete as DeleteIcon, ContentCopy as CopyIcon } from '@mui/icons-material';

interface EnvironmentCardProps {
  environment: Environment;
  onEdit: (environment: Environment) => void;
  onDelete: (environment: Environment) => void;
  onCopyApiKey: (apiKey: string) => void;
}

const EnvironmentCard: React.FC<EnvironmentCardProps> = ({
  environment,
  onEdit,
  onDelete,
  onCopyApiKey
}) => {
  return (
    <Card>
      <CardContent>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Box>
            <Typography variant="h6">{environment.name}</Typography>
            <Typography variant="body2" color="text.secondary">
              {environment.key}
            </Typography>
          </Box>
          <Box sx={{ display: 'flex', gap: 1 }}>
            <Tooltip title="Copy API Key">
              <IconButton onClick={() => onCopyApiKey(environment.api_key)}>
                <CopyIcon />
              </IconButton>
            </Tooltip>
            <IconButton onClick={() => onEdit(environment)}>
              <EditIcon />
            </IconButton>
            <IconButton onClick={() => onDelete(environment)}>
              <DeleteIcon />
            </IconButton>
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
};

export default EnvironmentCard;
```

**Файл: `webui/src/components/environments/EnvironmentDialog.tsx`**
```tsx
import React from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  Box
} from '@mui/material';

interface EnvironmentDialogProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: CreateEnvironmentRequest) => void;
  environment?: Environment;
  isEditing?: boolean;
}

const EnvironmentDialog: React.FC<EnvironmentDialogProps> = ({
  open,
  onClose,
  onSubmit,
  environment,
  isEditing = false
}) => {
  const [formData, setFormData] = React.useState({
    key: '',
    name: ''
  });

  React.useEffect(() => {
    if (environment && isEditing) {
      setFormData({
        key: environment.key,
        name: environment.name
      });
    } else {
      setFormData({ key: '', name: '' });
    }
  }, [environment, isEditing, open]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <form onSubmit={handleSubmit}>
        <DialogTitle>
          {isEditing ? 'Edit Environment' : 'Create Environment'}
        </DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, pt: 1 }}>
            <TextField
              label="Key"
              value={formData.key}
              onChange={(e) => setFormData({ ...formData, key: e.target.value })}
              disabled={isEditing}
              required
              helperText="Unique identifier for the environment (e.g., dev, stage, prod)"
            />
            <TextField
              label="Name"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              required
              helperText="Human-readable name for the environment"
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose}>Cancel</Button>
          <Button type="submit" variant="contained">
            {isEditing ? 'Update' : 'Create'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};

export default EnvironmentDialog;
```

#### 7.2 Обновление существующих компонентов

**Обновить `webui/src/components/features/FeatureCard.tsx`**:
- Добавить отображение состояния фичи по окружениям
- Добавить переключатель окружения

**Обновить `webui/src/pages/ProjectPage.tsx`**:
- Добавить вкладку "Environments"
- Добавить селектор окружения для фильтрации фич

**Обновить `webui/src/pages/ProjectSettingsPage.tsx`**:
- Удалить отображение API ключа проекта
- Добавить управление окружениями

#### 7.3 Новые страницы

**Файл: `webui/src/pages/ProjectEnvironmentsPage.tsx`**
```tsx
import React, { useState } from 'react';
import { Box, Paper, Typography, Button, CircularProgress } from '@mui/material';
import { Add as AddIcon } from '@mui/icons-material';
import { useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import EnvironmentCard from '../components/environments/EnvironmentCard';
import EnvironmentDialog from '../components/environments/EnvironmentDialog';
import apiClient from '../api/apiClient';

const ProjectEnvironmentsPage: React.FC = () => {
  const { projectId = '' } = useParams();
  const queryClient = useQueryClient();

  const { data: environments, isLoading } = useQuery({
    queryKey: ['environments', projectId],
    queryFn: async () => {
      const res = await apiClient.listProjectEnvironments(projectId);
      return res.data.items;
    },
    enabled: !!projectId,
  });

  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingEnvironment, setEditingEnvironment] = useState<Environment | null>(null);

  const createMutation = useMutation({
    mutationFn: async (data: CreateEnvironmentRequest) => {
      const res = await apiClient.createEnvironment(projectId, data);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['environments', projectId] });
      setDialogOpen(false);
    },
  });

  const updateMutation = useMutation({
    mutationFn: async ({ id, data }: { id: number; data: UpdateEnvironmentRequest }) => {
      const res = await apiClient.updateEnvironment(id, data);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['environments', projectId] });
      setDialogOpen(false);
      setEditingEnvironment(null);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: number) => {
      await apiClient.deleteEnvironment(id);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['environments', projectId] });
    },
  });

  const handleCreate = () => {
    setEditingEnvironment(null);
    setDialogOpen(true);
  };

  const handleEdit = (environment: Environment) => {
    setEditingEnvironment(environment);
    setDialogOpen(true);
  };

  const handleDelete = (environment: Environment) => {
    if (window.confirm(`Are you sure you want to delete environment "${environment.name}"?`)) {
      deleteMutation.mutate(environment.id);
    }
  };

  const handleCopyApiKey = (apiKey: string) => {
    navigator.clipboard.writeText(apiKey);
  };

  const handleSubmit = (data: CreateEnvironmentRequest | UpdateEnvironmentRequest) => {
    if (editingEnvironment) {
      updateMutation.mutate({ id: editingEnvironment.id, data: data as UpdateEnvironmentRequest });
    } else {
      createMutation.mutate(data as CreateEnvironmentRequest);
    }
  };

  return (
    <AuthenticatedLayout showBackButton backTo={`/projects/${projectId}`}>
      <Paper sx={{ p: 2 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
          <Typography variant="h6">Environments</Typography>
          <Button variant="contained" startIcon={<AddIcon />} onClick={handleCreate}>
            Add Environment
          </Button>
        </Box>

        {isLoading && (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        )}

        {environments && environments.length > 0 ? (
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
            {environments.map((env) => (
              <EnvironmentCard
                key={env.id}
                environment={env}
                onEdit={handleEdit}
                onDelete={handleDelete}
                onCopyApiKey={handleCopyApiKey}
              />
            ))}
          </Box>
        ) : (
          !isLoading && <Typography>No environments yet.</Typography>
        )}
      </Paper>

      <EnvironmentDialog
        open={dialogOpen}
        onClose={() => setDialogOpen(false)}
        onSubmit={handleSubmit}
        environment={editingEnvironment}
        isEditing={!!editingEnvironment}
      />
    </AuthenticatedLayout>
  );
};

export default ProjectEnvironmentsPage;
```

### Этап 8: Обновление SDK

#### 8.1 Обновление SDK для работы с окружениями

**Файл: `internal/api/sdk/environments.go`**
```go
package sdk

import (
    "context"
    "github.com/togglr-project/togglr/internal/domain"
    "github.com/togglr-project/togglr/internal/repository/environments"
)

type SDKEnvironmentService struct {
    envRepo *environments.Repository
}

func NewSDKEnvironmentService(envRepo *environments.Repository) *SDKEnvironmentService {
    return &SDKEnvironmentService{envRepo: envRepo}
}

// GetEnvironmentByAPIKey получает окружение по API ключу для SDK
func (s *SDKEnvironmentService) GetEnvironmentByAPIKey(ctx context.Context, apiKey string) (domain.Environment, error) {
    return s.envRepo.GetByAPIKey(ctx, apiKey)
}
```

#### 8.2 Обновление существующих SDK сервисов

**Обновить `internal/api/sdk/features.go`**:
- Добавить поддержку окружений в методы оценки фич
- Обновить методы для работы с параметрами фич

### Этап 9: Миграция данных

#### 9.1 Скрипт миграции данных

**Файл: `migrations/029_environments_data_migration.sql`**
```sql
-- Создание окружений для существующих проектов
INSERT INTO environments (project_id, key, name, api_key, created_at)
SELECT 
    p.id as project_id,
    'prod' as key,
    'Production' as name,
    p.api_key as api_key,
    NOW() as created_at
FROM projects p
WHERE p.api_key IS NOT NULL;

-- Создание параметров фич для существующих фич
INSERT INTO feature_params (feature_id, environment_id, enabled, default_value, created_at, updated_at)
SELECT 
    f.id as feature_id,
    e.id as environment_id,
    f.enabled as enabled,
    f.default_variant as default_value,
    NOW() as created_at,
    NOW() as updated_at
FROM features f
CROSS JOIN environments e
WHERE e.project_id = f.project_id;

-- Обновление правил с environment_id
UPDATE rules r
SET environment_id = e.id
FROM environments e
WHERE e.project_id = r.project_id;

-- Обновление вариантов флагов с environment_id
UPDATE flag_variants fv
SET environment_id = e.id
FROM environments e
WHERE e.project_id = fv.project_id;

-- Обновление расписаний фич с environment_id
UPDATE feature_schedules fs
SET environment_id = e.id
FROM environments e
WHERE e.project_id = fs.project_id;

-- Обновление аудита с environment_id
UPDATE audit_log al
SET environment_id = e.id
FROM environments e
WHERE e.project_id = al.project_id;

-- Обновление ожидающих изменений с environment_id
UPDATE pending_changes pc
SET environment_id = e.id
FROM environments e
WHERE e.project_id = pc.project_id;
```

### Этап 10: Тестирование

#### 10.1 Unit тесты

- Тесты для новых доменных моделей
- Тесты для новых репозиториев
- Тесты для новых use cases
- Тесты для новых API handlers

#### 10.2 Integration тесты

- Тесты для работы с окружениями
- Тесты для работы с параметрами фич
- Тесты для миграции данных

#### 10.3 E2E тесты

- Тесты для фронтенда с окружениями
- Тесты для SDK с окружениями

## Порядок выполнения

1. **Этап 1-2**: Создание доменных моделей и репозиториев
2. **Этап 3**: Обновление контрактов API
3. **Этап 4**: Создание use cases
4. **Этап 5**: Создание API handlers
5. **Этап 6**: Обновление OpenAPI схемы
6. **Этап 7**: Обновление фронтенда
7. **Этап 8**: Обновление SDK
8. **Этап 9**: Миграция данных
9. **Этап 10**: Тестирование

## Обратная совместимость

Для обеспечения обратной совместимости:

1. **API ключи проектов**: Создать окружение "prod" для каждого проекта с существующим API ключом
2. **Состояние фич**: Перенести `enabled` и `default_variant` в параметры фич для окружения "prod"
3. **Правила и варианты**: Связать с окружением "prod" по умолчанию
4. **SDK**: Обновить SDK для работы с окружениями, но сохранить обратную совместимость

## Риски и митигация

### Риски:
1. **Потеря данных** при миграции
2. **Нарушение работы SDK** после обновления
3. **Сложность фронтенда** с новыми концепциями

### Митигация:
1. **Резервное копирование** перед миграцией
2. **Поэтапное развертывание** с возможностью отката
3. **Тщательное тестирование** всех компонентов
4. **Документация** для пользователей

## Заключение

Данный план обеспечивает полный переход к модели с окружениями, сохраняя при этом обратную совместимость и обеспечивая плавную миграцию существующих данных. Все изменения структурированы по этапам, что позволяет контролировать процесс и минимизировать риски.
