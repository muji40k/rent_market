#! /bin/bash

function task {
    go run ./scripts/openapi/walk_paths.go paths | jq  --indent 4 > openapi/paths.json
}

export -f task

ROOT=./openapi/paths/ ./scripts/openapi/watch.sh

