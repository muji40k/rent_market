#! /bin/bash

echo gocyclo -over ${CYCLOMATIC} .
gocyclo -over ${CYCLOMATIC} .
rc=$?

echo go vet --unreachable ./...
go vet --unreachable ./... && exit $rc

