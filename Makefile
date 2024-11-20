
PERSONAL=1

PERSONAL_FLAGS=
DOCKER_ENV=$(shell ./scripts/expand-env.sh ./.env/env ./.env/env_${USER} ${env})

ifeq ($(PERSONAL), 1)
PERSONAL_FLAGS=-p ${USER}
endif

.PHONY: default
default: unit_tests

.PHONY: unit_tests unit_tests_build unit_tests_clean
unit_tests_build:
	docker compose -f docker/docker-compose.yml $(DOCKER_ENV) $(PERSONAL_FLAGS) build --no-cache backend_unit_tests test_db

unit_tests_clean:
	docker compose -f docker/docker-compose.yml $(DOCKER_ENV) $(PERSONAL_FLAGS) down --remove-orphans

unit_tests:
	docker compose -f docker/docker-compose.yml $(DOCKER_ENV) $(PERSONAL_FLAGS) run backend_unit_tests
	ALLURE_OUTPUT_PATH=$(PWD)/backend/allure-report/ $(MAKE) -C ./backend/ report

.PHONY: report_show
report_show:
	ALLURE_OUTPUT_PATH=$(PWD)/backend/allure-report/ $(MAKE) -C ./backend/ report_show

