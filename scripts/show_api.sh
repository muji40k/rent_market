#! /bin/bash

if [ -z "${1+x}" ]; then
    HOST=localhost
else
    HOST=$1
fi

echo "HOST: $HOST"

function task {
    print_log "Swagger up"
    swagger-ui-watcher -h $HOST --no-open openapi/main.json
    print_log "Swagger down"
}

export HOST
export -f task

ROOT=./openapi/ ./scripts/watch.sh

