Добавлены следующие ручки и структуры в OpenAPI спецификацию, сгенерирован сервер.
Необходимо реализовать методы repository (internal/repository/rbac/repository.go) для работы этих ручек, usecases (новый, по пути internal/usecases/rbac).
Методы хандлеров есть в internal/api/backend/api.go, их необходимо разнести по своим файлам (делать по аналогии с другими хандлерами) и реализовать.
При совершении действий над таблицей memberships необходимо писать лог в membership_audit (аудит). Смотри пример других репозиториев, как они пишут в audit_log таблицу (прямые вызовы Write() в auditlog repository, транзакционно, транзакция открывается в usecases).
Использовать гайд .junie/guideline.md.

```yaml
paths:
  /api/v1/roles:
    get:
      summary: List all roles
      operationId: ListRoles
      responses:
        '200':
          description: List of roles
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Role'
        '401':
          description: Unauthorized
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorUnauthorized'
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
      security:
        - bearerAuth: []

  /api/v1/roles/{role_id}/permissions:
    get:
      summary: Get permissions for a role
      operationId: GetRolePermissions
      parameters:
        - name: role_id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: List of permissions for role
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Permission'
        '401':
          description: Unauthorized
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorUnauthorized'
        '404':
          description: Role not found
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
      security:
        - bearerAuth: []

  /api/v1/permissions:
    get:
      summary: List all permissions
      operationId: ListPermissions
      responses:
        '200':
          description: List of all permissions
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Permission'
        '401':
          description: Unauthorized
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorUnauthorized'
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
      security:
        - bearerAuth: []

  /api/v1/roles/permissions:
    get:
      summary: List permissions for all roles
      operationId: ListRolePermissions
      responses:
        '200':
          description: Map of role to permissions
          content:
            application/json:
              schema:
                type: array
                items:
                  type: object
                  properties:
                    role:
                      $ref: '#/components/schemas/Role'
                    permissions:
                      type: array
                      items:
                        $ref: '#/components/schemas/Permission'
        '401':
          description: Unauthorized
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorUnauthorized'
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
      security:
        - bearerAuth: []

  /api/v1/projects/{project_id}/memberships:
    get:
      summary: List memberships for project
      operationId: ListProjectMemberships
      parameters:
        - name: project_id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: List of memberships
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Membership'
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
          description: Project not found
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
      security:
        - bearerAuth: []
    post:
      summary: Add membership to project
      operationId: CreateProjectMembership
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
              $ref: '#/components/schemas/CreateMembershipRequest'
      responses:
        '201':
          description: Membership created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Membership'
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
          description: Project not found
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
      security:
        - bearerAuth: []

  /api/v1/projects/{project_id}/memberships/{membership_id}:
    get:
      summary: Get membership
      operationId: GetProjectMembership
      parameters:
        - name: project_id
          in: path
          required: true
          schema:
            type: string
            format: uuid
        - name: membership_id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Membership
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Membership'
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
          description: Membership not found
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
      security:
        - bearerAuth: []
    put:
      summary: Update membership
      operationId: UpdateProjectMembership
      parameters:
        - name: project_id
          in: path
          required: true
          schema:
            type: string
            format: uuid
        - name: membership_id
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
              $ref: '#/components/schemas/UpdateMembershipRequest'
      responses:
        '200':
          description: Membership updated
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Membership'
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
          description: Membership not found
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
      security:
        - bearerAuth: []
    delete:
      summary: Delete membership
      operationId: DeleteProjectMembership
      parameters:
        - name: project_id
          in: path
          required: true
          schema:
            type: string
            format: uuid
        - name: membership_id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '204':
          description: Membership deleted
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
          description: Membership not found
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
      security:
        - bearerAuth: []

components:
  schemas:
    Role:
      type: object
      properties:
        id:
          type: string
          format: uuid
        key:
          type: string
        name:
          type: string
        description:
          type: string
      required: [id, key, name, description]

    Permission:
      type: object
      properties:
        id:
          type: string
          format: uuid
        key:
          type: string
          example: "feature.toggle"
        name:
          type: string
          example: "Toggle features"
      required: [id, key, name]

    Membership:
      type: object
      properties:
        id:
          type: string
          format: uuid
        user_id:
          type: integer
          format: int64
        project_id:
          type: string
          format: uuid
        role_id:
          type: string
          format: uuid
        role_key:
          type: string
          example: "project_manager"
        role_name:
          type: string
          example: "Project Manager"
        created_at:
          type: string
          format: date-time
      required: [id, user_id, project_id, role_id, role_key, role_name, created_at]

    CreateMembershipRequest:
      type: object
      properties:
        user_id:
          type: integer
          format: int64
          example: 42
        role_id:
          type: string
          format: uuid
          example: "c4f1a8d0-3b75-4c9f-9181-bbba66ad635d"
      required: [user_id, role_id]

    UpdateMembershipRequest:
      type: object
      properties:
        role_id:
          type: string
          format: uuid
          example: "c4f1a8d0-3b75-4c9f-9181-bbba66ad635d"
      required: [role_id]
```

---

Сейчас такая схема БД:
```sql
create table memberships
(
    id         uuid                     default gen_random_uuid() not null
        primary key,
    project_id uuid                                               not null
        references projects
            on delete cascade,
    user_id    integer                                            not null
        references users
            on delete cascade,
    role_id    uuid                                               not null
        references roles
            on delete restrict,
    created_at timestamp with time zone default now()             not null,
    updated_at timestamp with time zone default now()             not null,
    constraint membership_unique
        unique (project_id, user_id)
);

create table membership_audit
(
    id            bigserial
        primary key,
    membership_id uuid,
    actor_user_id integer,
    action        varchar(50)                            not null,
    old_value     jsonb,
    new_value     jsonb,
    created_at    timestamp with time zone default now() not null
);

create table roles
(
    id          uuid                     default gen_random_uuid() not null
        primary key,
    key         varchar(50)                                        not null
        unique,
    name        varchar(50)                                        not null,
    description varchar(300),
    created_at  timestamp with time zone default now()             not null
);

create table role_permissions
(
    id            uuid default gen_random_uuid() not null
        primary key,
    role_id       uuid                           not null
        references roles
            on delete cascade,
    permission_id uuid                           not null
        references permissions
            on delete cascade,
    constraint role_permissions_unique
        unique (role_id, permission_id)
);

create table permissions
(
    id   uuid default gen_random_uuid() not null
        primary key,
    key  varchar(50)                    not null
        unique,
    name varchar(50)                    not null
);
```
