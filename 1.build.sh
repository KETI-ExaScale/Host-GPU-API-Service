#!/bin/bash

dockerID="ketidevit2"

go mod tidy
go mod vendor

go build -o build/bin/host-api-service -mod=vendor cmd/main.go
