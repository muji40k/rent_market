
FILES := $(shell find ./openapi/ -type f)

.PHONY: openapi

openapi: openapi.json

openapi.%: $(FILES)
	swagger-ui-watcher -b $@ ./openapi/main.json

