_COMPOSE_BASE=docker compose -f dev/docker-compose.yml --project-name ${NAMESPACE} --env-file dev/platform.env
_COMPOSE_EXT=docker compose -f dev/docker-compose.ext.yml --project-name ${NAMESPACE} --env-file dev/platform.env

dev-up: ## Up the environment in docker compose (base services only)
	${_COMPOSE_BASE} up -d

dev-up-ext: ## Up the environment in docker compose with LDAP services
	./dev/generate-ext.sh
	${_COMPOSE_EXT} up -d

dev-down: ## Down the environment in docker compose
	${_COMPOSE_BASE} down --remove-orphans
	${_COMPOSE_EXT} down --remove-orphans

dev-clean: ## Down the environment in docker compose with image cleanup
	${_COMPOSE_BASE} down --remove-orphans -v --rmi all
	${_COMPOSE_EXT} down --remove-orphans -v --rmi all

dev-build-proxy: ## Building togglr-reverse-proxy
	${_COMPOSE_BASE} build togglr-reverse-proxy

dev-build-frontend: ## Building togglr-frontend
	${_COMPOSE_BASE} build togglr-frontend

dev-build-backend: ## Building togglr-backend
	${_COMPOSE_BASE} build togglr-backend
