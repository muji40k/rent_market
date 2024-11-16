
.PHONY: default
default: unit_tests

.PHONY: unit_tests
unit_tests:
	docker compose -f docker/docker-compose.yml run backend_unit_tests
	ALLURE_OUTPUT_PATH=$(PWD)/backend/allure-report/ $(MAKE) -C ./backend/ report

