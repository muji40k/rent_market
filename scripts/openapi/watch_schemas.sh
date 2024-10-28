#! /bin/bash

function task {
    go run ./scripts/openapi/walk_paths.go schemas | jq  --indent 4 > openapi/components/schemas.json
}

export -f task

ROOT=./openapi/components/schemas/ ./scripts/openapi/watch.sh

