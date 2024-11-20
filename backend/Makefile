GO_FILES = $(shell find . -type f -name "*.go")

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

.PHONY: default
default: build

.PHONY: build run
build: main
run: build
	./main

main: $(GO_FILES)
	go build ./cmd/main.go

.PHONY: unit_tests
unit_tests:
	$(MAKE) -C ${ALLURE_OUTPUT_PATH} before
	go test $(TEST_FLAGS)  ./internal/logic/services/implementations/... $(PIPE)
	go test $(TEST_FLAGS) -p 1  ./internal/repository/implementation/... $(PIPE)
	$(MAKE) -C ${ALLURE_OUTPUT_PATH} report

.PHONY: report_show
report_show:
	$(MAKE) -C ${ALLURE_OUTPUT_PATH} show

.PHONY: test
test:
	echo ${ALLURE_OUTPUT_PATH}

