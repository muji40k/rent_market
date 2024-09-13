#! /bin/bash

function task {
    print_log "Swagger up"
    swagger-ui-watcher --no-open openapi/main.json
    print_log "Swagger down"
}

export -f task

ROOT=./openapi/ ./scripts/watch.sh

