# Быстрый справочник по добавлению эндпоинтов в OpenAPI

## Структура файлов
- `server.yml` - полная спецификация (6871 строк)
- `server_base.yml` - базовая структура с примерами и основными схемами
- `QUICK_REFERENCE.md` - этот файл

## Основные паттерны эндпоинтов

### 1. Стандартная структура GET эндпоинта
```yaml
/api/v1/resource:
  get:
    summary: List resources
    operationId: ListResources
    parameters:
      - name: page
        in: query
        required: false
        schema:
          type: integer
          minimum: 1
          default: 1
      - name: per_page
        in: query
        required: false
        schema:
          type: integer
          minimum: 1
          default: 20
    responses:
      '200':
        description: List of resources
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ListResourcesResponse'
      '401':
        $ref: '#/components/responses/Unauthorized'
      '403':
        $ref: '#/components/responses/Forbidden'
      '500':
        $ref: '#/components/responses/InternalServerError'
    security:
      - bearerAuth: []
```

### 2. Стандартная структура POST эндпоинта
```yaml
/api/v1/resource:
  post:
    summary: Create resource
    operationId: CreateResource
    requestBody:
      required: true
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/CreateResourceRequest'
    responses:
      '201':
        description: Resource created
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ResourceResponse'
      '400':
        $ref: '#/components/responses/BadRequest'
      '401':
        $ref: '#/components/responses/Unauthorized'
      '403':
        $ref: '#/components/responses/Forbidden'
      '500':
        $ref: '#/components/responses/InternalServerError'
    security:
      - bearerAuth: []
```

### 3. Стандартная структура PUT/PATCH эндпоинта
```yaml
/api/v1/resource/{id}:
  put:
    summary: Update resource
    operationId: UpdateResource
    parameters:
      - name: id
        in: path
        required: true
        schema:
          type: string
    requestBody:
      required: true
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/UpdateResourceRequest'
    responses:
      '200':
        description: Resource updated
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ResourceResponse'
      '400':
        $ref: '#/components/responses/BadRequest'
      '401':
        $ref: '#/components/responses/Unauthorized'
      '403':
        $ref: '#/components/responses/Forbidden'
      '404':
        $ref: '#/components/responses/NotFound'
      '500':
        $ref: '#/components/responses/InternalServerError'
    security:
      - bearerAuth: []
```

### 4. Стандартная структура DELETE эндпоинта
```yaml
/api/v1/resource/{id}:
  delete:
    summary: Delete resource
    operationId: DeleteResource
    parameters:
      - name: id
        in: path
        required: true
        schema:
          type: string
    responses:
      '204':
        description: Resource deleted
      '401':
        $ref: '#/components/responses/Unauthorized'
      '403':
        $ref: '#/components/responses/Forbidden'
      '404':
        $ref: '#/components/responses/NotFound'
      '500':
        $ref: '#/components/responses/InternalServerError'
    security:
      - bearerAuth: []
```

## Стандартные коды ответов

### Успешные ответы
- `200` - OK (GET, PUT, PATCH)
- `201` - Created (POST)
- `204` - No Content (DELETE)

### Ошибки клиента
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict
- `429` - Too Many Requests

### Ошибки сервера
- `500` - Internal Server Error

## Стандартные схемы ошибок

Все ошибки наследуются от базовой схемы `Error`:

```yaml
Error:
  type: object
  properties:
    error:
      type: object
      properties:
        message:
          type: string
  required: [error]
```

### Специфичные ошибки:
- `ErrorBadRequest` - 400
- `ErrorUnauthorized` - 401
- `ErrorPermissionDenied` - 403
- `ErrorNotFound` - 404
- `ErrorConflict` - 409
- `ErrorTooManyRequests` - 429
- `ErrorInternalServerError` - 500

## Стандартные схемы

### Пагинация
```yaml
Pagination:
  type: object
  properties:
    total:
      type: integer
      format: uint
    page:
      type: integer
      format: uint
      minimum: 1
    per_page:
      type: integer
      format: uint
      minimum: 1
  required: [total, page, per_page]
```

### Сортировка
```yaml
SortOrder:
  type: string
  enum: [asc, desc]
```

### Успешный ответ
```yaml
SuccessResponse:
  type: object
  properties:
    message:
      type: string
```

## Группы эндпоинтов в проекте

1. **Authentication** (`/api/v1/auth/*`)
   - login, refresh, forgot-password, reset-password
   - 2fa/verify
   - sso/callback, sso/initiate, sso/providers

2. **SAML** (`/api/v1/saml/*`)
   - metadata, acs

3. **License** (`/api/v1/license/*`)
   - status, create/update

4. **Users** (`/api/v1/users/*`)
   - me, change-password, license-acceptance
   - list, get by id, superuser, active

5. **Projects** (`/api/v1/projects/*`)
   - list, create, get by id, update, delete
   - features, rules, tags, categories, etc.

## Быстрые шаблоны для копирования

### GET список с фильтрацией
```yaml
parameters:
  - name: text_selector
    in: query
    required: false
    description: Case-insensitive text search
    schema:
      type: string
  - name: sort_by
    in: query
    required: false
    schema:
      type: string
      enum: [name, created_at, updated_at]
  - name: sort_order
    in: query
    required: false
    schema:
      $ref: '#/components/schemas/SortOrder'
  - name: page
    in: query
    required: false
    schema:
      type: integer
      minimum: 1
      default: 1
  - name: per_page
    in: query
    required: false
    schema:
      type: integer
      minimum: 1
      default: 20
```

### Стандартные ответы
```yaml
responses:
  '200':
    description: Success
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/ListResponse'
  '400':
    description: Bad request
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/ErrorBadRequest'
  '401':
    description: Unauthorized
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/ErrorUnauthorized'
  '403':
    description: Permission denied
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/ErrorPermissionDenied'
  '404':
    description: Not found
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/ErrorNotFound'
  '500':
    description: Internal server error
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/ErrorInternalServerError'
  default:
    description: Unexpected error
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/Error'
```

## Рекомендации

1. **Всегда используйте `server_base.yml`** как основу для понимания структуры
2. **Копируйте существующие эндпоинты** и адаптируйте под новые
3. **Используйте стандартные схемы ошибок** из `server_base.yml`
4. **Следуйте паттернам именования**: `ListResources`, `CreateResource`, `UpdateResource`, `DeleteResource`
5. **Добавляйте `security: - bearerAuth: []`** ко всем защищенным эндпоинтам
6. **Используйте `operationId`** для генерации клиентского кода

## Где искать примеры

- В `server_base.yml` есть полный пример эндпоинта `/api/v1/projects/{project_id}/features`
- В `server.yml` ищите похожие эндпоинты по группам (auth, users, projects)
- Используйте grep для поиска паттернов: `grep -A 10 "operationId: List" specs/server.yml`
