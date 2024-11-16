
.PHONY: default unit_tests report_show

SHUFFLE=1
RACE=0
SHORT=1

PIPE=
TEST_FLAGS=

ifeq ($(SHUFFLE), 1)
TEST_FLAGS += -shuffle on
endif

ifeq ($(RACE), 1)
TEST_FLAGS += -race
endif

ifeq ($(SHORT), 1)
PIPE=| grep -E -v "no (test files|tests to run)"
endif

default: unit_tests

unit_tests:
	$(MAKE) -f allure-report/Makefile before
	ALLURE_OUTPUT_PATH=$(PWD)/allure-report/ TEST_DB_PORT=5434 go test $(TEST_FLAGS)  ./internal/logic/services/implementations/... $(PIPE)
	ALLURE_OUTPUT_PATH=$(PWD)/allure-report/ TEST_DB_PORT=5434 go test $(TEST_FLAGS) -p 1  ./internal/repository/implementation/... $(PIPE)
	$(MAKE) -f allure-report/Makefile report

report_show:
	$(MAKE) -f allure-report/Makefile show

