RED="\033[0;31m"
GREEN="\033[1;32m"
NOCOLOR="\033[0m"

.PHONY: docker-build-backend
docker-build-backend: ## Building Docker image for backend (scratch + curl)
	@echo "\nBuilding Docker image (scratch + curl)..."
	@docker build \
		--build-arg TOOL_VERSION=${TOOL_VERSION} \
		--build-arg TOOL_BUILD_TIME=${TOOL_BUILD_TIME} \
		--build-arg DOCKER_REGISTRY=${DOCKER_REGISTRY} \
		-t togglr-backend:latest -f ./Dockerfile.backend .
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"FAIL"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo ${GREEN}"Docker image 'togglr-backend' built successfully!"${NOCOLOR}

.PHONY: docker-push-backend
docker-push-backend: docker-build-backend ## Tagging and pushing Docker image for backend
	@echo "\nTagging Docker backend image..."
	@docker tag togglr-backend:latest $(DOCKER_REGISTRY)/rom8726/togglr-backend:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Tagging FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo "\nPushing Docker backend image to Docker Registry..."
	@docker push $(DOCKER_REGISTRY)/rom8726/togglr-backend:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Push FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@if [ -n "$(TOOL_VERSION)" ] && [ "$(TOOL_VERSION)" != "latest" ]; then \
		echo "\nTagging Docker backend image with version $(TOOL_VERSION)..."; \
		docker tag togglr-backend:latest $(DOCKER_REGISTRY)/rom8726/togglr-backend:$(TOOL_VERSION); \
		if [ $$? -ne 0 ] ; then \
			echo -e ${RED}"Version tagging FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
		echo "\nPushing Docker backend image with version $(TOOL_VERSION) to Docker Registry..."; \
		docker push $(DOCKER_REGISTRY)/rom8726/togglr-backend:$(TOOL_VERSION); \
		if [ $$? -ne 0 ] ; then \
			echo -e ${RED}"Version push FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
	fi
	@echo ${GREEN}"\nDocker backend image pushed to registry successfully!"${NOCOLOR}

.PHONY: docker-build-frontend
docker-build-frontend: ## Building Docker image for frontend
	@echo "\nBuilding Docker image frontend..."
	@cd ui && TOOL_VERSION=${TOOL_VERSION} TOOL_BUILD_TIME=${TOOL_BUILD_TIME} VITE_IS_DEMO=${VITE_IS_DEMO} make docker-build
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"FAIL"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo ${GREEN}"Docker image 'togglr-frontend' built successfully!"${NOCOLOR}

.PHONY: docker-push-frontend
docker-push-frontend: docker-build-frontend ## Tagging and pushing Docker image for frontend
	@echo "\nTagging Docker frontend image..."
	@docker tag togglr-frontend:latest $(DOCKER_REGISTRY)/rom8726/togglr-frontend:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Tagging FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo "\nPushing Docker frontend image to Docker Registry..."
	@docker push $(DOCKER_REGISTRY)/rom8726/togglr-frontend:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Push FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@if [ -n "$(TOOL_VERSION)" ] && [ "$(TOOL_VERSION)" != "latest" ]; then \
		echo "\nTagging Docker frontend image with version $(TOOL_VERSION)..."; \
		docker tag togglr-frontend:latest $(DOCKER_REGISTRY)/rom8726/togglr-frontend:$(TOOL_VERSION); \
		if [ $$? -ne 0 ] ; then \
			echo -e ${RED}"Version tagging FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
		echo "\nPushing Docker frontend image with version $(TOOL_VERSION) to Docker Registry..."; \
		docker push $(DOCKER_REGISTRY)/rom8726/togglr-frontend:$(TOOL_VERSION); \
		if [ $$? -ne 0 ] ; then \
			echo -e ${RED}"Version push FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
	fi
	@echo ${GREEN}"\nDocker frontend image pushed to registry successfully!"${NOCOLOR}

.PHONY: docker-build-reverse-proxy
docker-build-reverse-proxy: ## Building Docker image for reverse-proxy
	@echo "\nBuilding Docker image reverse-proxy..."
	@cd reverse-proxy && TOOL_VERSION=${TOOL_VERSION} TOOL_BUILD_TIME=${TOOL_BUILD_TIME} make docker-build
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"FAIL"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo ${GREEN}"Docker image 'togglr-reverse-proxy' built successfully!"${NOCOLOR}

.PHONY: docker-push-reverse-proxy
docker-push-reverse-proxy: docker-build-reverse-proxy ## Tagging and pushing Docker image for reverse-proxy
	@echo "\nTagging Docker reverse-proxy image..."
	@docker tag togglr-reverse-proxy:latest $(DOCKER_REGISTRY)/rom8726/togglr-reverse-proxy:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Tagging FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo "\nPushing Docker reverse-proxy image to Docker Registry..."
	@docker push $(DOCKER_REGISTRY)/rom8726/togglr-reverse-proxy:latest
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"Push FAILED"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@if [ -n "$(TOOL_VERSION)" ] && [ "$(TOOL_VERSION)" != "latest" ]; then \
		echo "\nTagging Docker reverse-proxy image with version $(TOOL_VERSION)..."; \
		docker tag togglr-reverse-proxy:latest $(DOCKER_REGISTRY)/rom8726/togglr-reverse-proxy:$(TOOL_VERSION); \
		if [ $$? -ne 0 ] ; then \
			echo -e ${RED}"Version tagging FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
		echo "\nPushing Docker reverse-proxy image with version $(TOOL_VERSION) to Docker Registry..."; \
		docker push $(DOCKER_REGISTRY)/rom8726/togglr-reverse-proxy:$(TOOL_VERSION); \
		if [ $$? -ne 0 ] ; then \
			echo -e ${RED}"Version push FAILED"${NOCOLOR} ; \
			exit 1 ; \
		fi; \
	fi
	@echo ${GREEN}"\nDocker reverse-proxy image pushed to registry successfully!"${NOCOLOR}
