#!/usr/bin/env bash

set -Eeou pipefail


# If CI is enabled (Bitbucket Pipelines) create a test PostgreSQL container
# and wait a bit so it will be available in the tests
if [[ ${CI:-} == "true" ]]; then
    echo "Installing wkhtmltopdf before tests..."
    apt update
    apt install -y wkhtmltopdf 

    echo "Starting test PostgreSQL Docker container..."
    docker run --name migration-test-db -e POSTGRES_PASSWORD=pass123 -p 5432:5432 -d postgres
    echo "Starting test Microsoft SQL Server Docker container..."
    docker run -e "ACCEPT_EULA=Y" -e "SA_PASSWORD=Pass1234" -p 1433:1433 --name mssql-test -h mssql-test -d mcr.microsoft.com/mssql/server:2019-latest
    echo "Waiting for the databases to be availeable..."
    sleep 6
fi

# Clean test chache, and run all the tests, capture coverage profile.
go clean -testcache &&  go test ./... -json -coverprofile=go-test-coverage.out | tee go-test-report.json

# go clean -testcache &&  go test ./... -coverprofile cp.out | tee test.results


# Run coverage check against the results
# go run scripts/coverage.go test.results
