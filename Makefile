
PERSONAL=1

PERSONAL_FLAGS=
DOCKER_ENV=$(shell ./scripts/expand-env.sh ./.env/env ./.env/env_${USER} ${env})

ifeq ($(PERSONAL), 1)
PERSONAL_FLAGS=-p ${USER}
endif

BACKEND_BUILD_FLAGS=

ifeq ($(NO_CACHE), 1)
BACKEND_BUILD_FLAGS += --no-cache
endif

.PHONY: default
default: backend

.PHONY: backend
backend:
	docker compose -f docker/docker-compose.yml $(DOCKER_ENV) build $(BACKEND_BUILD_FLAGS)
	-docker compose -f docker/docker-compose.yml $(DOCKER_ENV) up nginx backend \
		backend_ro1 backend_ro2 postgresql_db pgadmin \
		mirror_nginx mirror_backend mirror_postgresql_db
	docker compose -f docker/docker-compose.yml $(DOCKER_ENV) down

.PHONY: unit_tests unit_tests_run unit_tests_build unit_tests_clean
unit_tests_build:
	docker compose -f docker/docker-compose.yml $(DOCKER_ENV) $(PERSONAL_FLAGS) build --no-cache backend_unit_tests test_db

unit_tests_clean:
	docker compose -f docker/docker-compose.yml $(DOCKER_ENV) $(PERSONAL_FLAGS) down --remove-orphans

unit_tests_run:
	-docker compose -f docker/docker-compose.yml $(DOCKER_ENV) $(PERSONAL_FLAGS) run backend_unit_tests
	ALLURE_OUTPUT_PATH=$(PWD)/backend/allure-report/ $(MAKE) -C ./backend/ report

unit_tests:
	$(MAKE) unit_tests_clean
	$(MAKE) unit_tests_build
	-$(MAKE) unit_tests_run
	$(MAKE) unit_tests_clean

.PHONY: integration_tests integration_tests_run integration_tests_build integration_tests_clean
integration_tests_build:
	docker compose -f docker/docker-compose.yml $(DOCKER_ENV) $(PERSONAL_FLAGS) build --no-cache backend_integration_tests test_db

integration_tests_clean:
	docker compose -f docker/docker-compose.yml $(DOCKER_ENV) $(PERSONAL_FLAGS) down --remove-orphans

integration_tests_run:
	-docker compose -f docker/docker-compose.yml $(DOCKER_ENV) $(PERSONAL_FLAGS) run backend_integration_tests
	ALLURE_OUTPUT_PATH=$(PWD)/backend/allure-report/ $(MAKE) -C ./backend/ report

integration_tests:
	$(MAKE) integration_tests_clean
	$(MAKE) integration_tests_build
	-$(MAKE) integration_tests_run
	$(MAKE) integration_tests_clean

.PHONY: e2e_tests e2e_tests_run e2e_tests_build e2e_tests_clean
e2e_tests_build:
	docker compose -f docker/docker-compose.yml $(DOCKER_ENV) $(PERSONAL_FLAGS) build --no-cache backend_e2e_tests test_db

e2e_tests_clean:
	docker compose -f docker/docker-compose.yml $(DOCKER_ENV) $(PERSONAL_FLAGS) down --remove-orphans

e2e_tests_run:
	-docker compose -f docker/docker-compose.yml $(DOCKER_ENV) $(PERSONAL_FLAGS) run backend_e2e_tests
	ALLURE_OUTPUT_PATH=$(PWD)/backend/allure-report/ $(MAKE) -C ./backend/ report

e2e_tests:
	$(MAKE) e2e_tests_clean
	$(MAKE) e2e_tests_build
	-$(MAKE) e2e_tests_run
	$(MAKE) e2e_tests_clean

.PHONY: bdd_e2e_tests bdd_e2e_tests_run bdd_e2e_tests_build bdd_e2e_tests_clean
bdd_e2e_tests_build:
	docker compose -f docker/docker-compose.yml $(DOCKER_ENV) $(PERSONAL_FLAGS) build --no-cache backend_bdd_e2e_tests test_db

bdd_e2e_tests_clean:
	docker compose -f docker/docker-compose.yml $(DOCKER_ENV) $(PERSONAL_FLAGS) down --remove-orphans

bdd_e2e_tests_run:
	-docker compose -f docker/docker-compose.yml $(DOCKER_ENV) $(PERSONAL_FLAGS) run backend_bdd_e2e_tests

bdd_e2e_tests:
	$(MAKE) bdd_e2e_tests_clean
	$(MAKE) bdd_e2e_tests_build
	-$(MAKE) bdd_e2e_tests_run
	$(MAKE) bdd_e2e_tests_clean

.PHONY: sandbox sandbox_down
sandbox:
	docker compose -f docker/docker-compose.yml --env-file .env/sandbox build backend_sandbox debug_db
	docker compose -f docker/docker-compose.yml --env-file .env/sandbox up -d backend_sandbox debug_db

sandbox_down:
	docker compose -f docker/docker-compose.yml --env-file .env/sandbox down --remove-orphans backend_sandbox debug_db

.PHONY: report_show
report_show:
	ALLURE_OUTPUT_PATH=$(PWD)/backend/allure-report/ $(MAKE) -C ./backend/ report_show

