Принято ✅
Составлю подробное техзадание в формате, который можно отдать AI-агенту/разработчику для реализации RBAC-проверок на фронтенде. Оно будет опираться на `GET /api/v1/users/me` и поле `project_permissions`.

---

# 📑 ТЗ: Реализация RBAC на фронтенде

## 🎯 Цель

Обеспечить на фронтенде корректное разграничение доступа к функционалу Togglr в зависимости от роли пользователя и набора его прав (permissions), получаемых от backend.

---

## 📡 Источник данных

* `GET /api/v1/users/me` возвращает объект `User`, в котором есть поле:

  ```json
  "project_permissions": {
    "<project_id>": ["project.view", "feature.view", "feature.toggle", ...]
  }
  ```
* `is_superuser = true` означает, что пользователь имеет доступ ко всему функционалу без ограничений.

---

## 🛠 Фронтенд-слой авторизации

1. **Context/AuthStore**

    * Хранить текущего пользователя (`User`) и его `project_permissions` в глобальном сторе (например, Redux/Zustand/Pinia).
    * Добавить утилиту `hasPermission(projectId: string, perm: string): boolean`, которая:

        * Возвращает `true`, если `is_superuser = true`;
        * Возвращает `true`, если `perm` есть в массиве `project_permissions[projectId]`;
        * Иначе `false`.

2. **Хелперы/Guard-функции**

    * `canViewProject(projectId)` → проверяет `project.view`;
    * `canManageProject(projectId)` → проверяет `project.manage`;
    * `canToggleFeature(projectId)` → проверяет `feature.toggle`;
    * `canManageFeature(projectId)` → проверяет `feature.manage`;
    * `canManageSegment(projectId)` → проверяет `segment.manage`;
    * `canManageSchedule(projectId)` → проверяет `schedule.manage`;
    * `canViewAudit(projectId)` → проверяет `audit.view`;
    * `canManageMembership(projectId)` → проверяет `membership.manage`.

   ⚠️ Эти функции должны использовать `hasPermission` внутри.

---

## 🖥 Применение в UI

1. **Проекты**

    * Если нет `project.view` → проект не отображается в списке.
    * Если есть `project.view`, но нет `project.manage` → UI доступен только для чтения (без кнопок редактирования/удаления).

2. **Фичи**

    * `feature.view` → можно видеть список фич.
    * `feature.toggle` → кнопка включения/выключения активна. Если нет → disabled.
    * `feature.manage` → можно редактировать настройки фичи, добавлять варианты и т. д.

3. **Сегменты**

    * `segment.manage` → можно создавать, редактировать, удалять сегменты.
    * Если нет → скрыть/задизейблить UI для изменения сегментов.

4. **Расписания**

    * `schedule.manage` → доступ к настройке расписаний.

5. **Аудит**

    * `audit.view` → доступ к вкладке/разделу «История изменений».
    * Если нет → вкладка скрыта.

6. **Участники**

    * `membership.manage` → можно добавлять/удалять участников проекта и менять их роли.

7. **Глобальный суперюзер**

    * Если `is_superuser = true`, то фронт **игнорирует проверки** и показывает все элементы.

---

## 🧪 Тестирование

* Создать тестовые аккаунты с разными ролями (`project_owner`, `project_manager`, `project_member`).
* Проверить:

    * У `project_member` нет доступа к toggle фичей, но есть `feature.view`.
    * У `project_manager` есть доступ к toggle/manage features, segment.manage, schedule.manage, audit.view.
    * У `project_owner` полный доступ.
* Проверить суперюзера (`is_superuser=true`) — доступны все элементы.

---

## 🚀 Критерии готовности

* Все компоненты UI проверяют доступ через `canXxx()` функции.
* Пользователь никогда не видит кнопок/действий, на которые у него нет прав.
* При ручном вызове недоступного действия (например, через DevTools) запрос к API всё равно будет отклонён backend-ом (двойная защита).

---

Отлично 👍 Составил таблицу «Роль → Доступный функционал» на основе текущих прав. Она поможет фронтенду быстро ориентироваться, что показывать, а что скрывать.

---

# 📊 RBAC: Роли и доступный функционал
```json
[
  {
    "role_key": "project_manager",
    "role_name": "Project Manager",
    "permissions": "audit.view, feature.manage, feature.toggle, feature.view, project.view, schedule.manage, segment.manage"
  },
  {
    "role_key": "project_member",
    "role_name": "Project Member",
    "permissions": "feature.toggle, feature.view, project.view"
  },
  {
    "role_key": "project_owner",
    "role_name": "Project Owner",
    "permissions": "audit.view, feature.manage, feature.toggle, feature.view, membership.manage, project.manage, project.view, schedule.manage, segment.manage"
  },
  {
    "role_key": "project_viewer",
    "role_name": "Project Viewer",
    "permissions": "feature.view, project.view"
  }
]
```
---

## 📝 Интерпретация

* **project_owner**: полный доступ ко всем функциям проекта + может создавать проекты.
* **project_manager**: почти полный доступ внутри проекта (фичи, сегменты, расписания, аудит, управление участниками), но не может создавать проекты.
* **project_member**: может просматривать проект и фичи и включать-выключать фичи.
* **project_viewer**: может только просматривать проект и фичи.

---

⚡ Важно: при `is_superuser = true` фронт показывает весь функционал независимо от роли.

---
