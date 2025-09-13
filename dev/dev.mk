_COMPOSE=docker compose -f dev/docker-compose.yml --project-name ${NAMESPACE} --env-file dev/platform.env

dev-up: ## Up the environment in docker compose
	${_COMPOSE} up -d

dev-down: ## Down the environment in docker compose
	${_COMPOSE} down --remove-orphans

dev-clean: ## Down the environment in docker compose with image cleanup
	${_COMPOSE} down --remove-orphans -v --rmi all

dev-build-proxy: ## Building etoggle-reverse-proxy
	${_COMPOSE} build etoggle-reverse-proxy

dev-build-frontend: ## Building etoggle-frontend
	${_COMPOSE} build etoggle-frontend

dev-build-backend: ## Building etoggle-backend
	${_COMPOSE} build etoggle-backend
