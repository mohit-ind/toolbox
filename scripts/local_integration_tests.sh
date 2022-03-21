#!/usr/bin/env bash

function log() {
    echo "$(date) - ${1}"
}

log "Starting test PostgreSQL database Docker container: migration-test-db"
docker run --name migration-test-db -e POSTGRES_PASSWORD=pass123 -p 5432:5432 -d postgres
log "Starting test Microsoft SQL Server Docker container: mssql-test"
docker run -e "ACCEPT_EULA=Y" -e "SA_PASSWORD=Pass1234" -p 1433:1433 --name mssql-test -h mssql-test -d mcr.microsoft.com/mssql/server:2019-latest
log "Wait for the databases to be reachable"
sleep 3

log "Running Go tests"
export CI=true
export TEST_DB_CONNECTION_STRING="user=postgres dbname=postgres password=pass123 sslmode=disable"
go clean -testcache && go test ./aws/... -coverprofile cp_aws.out > test_aws.results
go clean -testcache && go test ./database/... -coverprofile cp_database.out > test_database.results
go clean -testcache && go test ./migration/... -coverprofile cp_migration.out > test_migration.results

log "Removing test PostgreSQL database Docker container: migration-test-db"
docker rm -f migration-test-db

log "Removing test Microsoft SQL Server Docker container: mssql-test"
docker rm -f mssql-test
