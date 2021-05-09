#!/bin/sh

FILES=$(go list $PWD/... | grep -v vendor)

mkdir -p artifacts
go vet $FILES
golint $FILES
go test $FILES -bench=. -race -timeout 10000ms -coverprofile artifacts/cover.out
go tool cover -html=artifacts/cover.out -o artifacts/cover.html
go tool cover -func=artifacts/cover.out

# To get codecov to accept us
go test -race -coverprofile=coverage.txt -covermode=atomic