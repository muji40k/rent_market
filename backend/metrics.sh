#! /bin/bash

echo gocyclo -over ${CYCLOMATIC} .
gocyclo -over ${CYCLOMATIC} . | tee /dev/fd/2 | [ 0 -eq $(wc -l) ]

echo go vet --unreachable ./...
go vet --unreachable ./...

