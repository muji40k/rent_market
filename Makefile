
PERSONAL=1

PERSONAL_FLAGS=
DOCKER_ENV=$(shell ./scripts/expand-env.sh ./.env/env ./.env/env_${USER} ${env})

ifeq ($(PERSONAL), 1)
PERSONAL_FLAGS=-p ${USER}
endif

.PHONY: default
default: unit_tests

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

.PHONY: report_show
report_show:
	ALLURE_OUTPUT_PATH=$(PWD)/backend/allure-report/ $(MAKE) -C ./backend/ report_show

