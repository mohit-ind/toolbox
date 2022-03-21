# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.18.8] - 2022-01-03

### Added
- sftp upload function under /services/sftp packages

## [1.18.7] - 2021-12-23

### Added

- Emailer services under /services/email packages
- Covered every function with test cases

## [1.18.6] - 2021-11-22

### Added

- WKHTMLTOPDF_PATH to pdf-generator, so the path to the wkhtmltopdf tool can be explicitly set

## [1.18.5] - 2021-11-22

### Updated

- updated changelog.md
- updated pdf-generator

## [1.18.4] - 2021-11-19

### Updated

- updated changelog.md
- updated pdf-generator

## [1.18.3] - 2021-11-02

### Added

- CSV - generator function
- test-cases file to cover my changes.

### Updated

- updated .gitignore file

## [1.18.2] - 2021-10-23

### Added

- Download file from S3 to memory
- Proper S3 file uploader, also capable to upload data from memory

## [1.18.1] - 2021-09-20

### Added

- shell package with a command structure, which is capable to execute shell commands
- Ping function to Safe Json

## [1.18.0] - 2021-09-19

### Added

- S3 File Downloader to the aws package
- Safe unzip file
- Safe Json added to database package
- TestFile function to file package

### Changed

- Repository's Go version from 1.15 to 1.17

## [1.17.8] - 2021-09-17

### Changed

- updated changelog.md
- changes header colour in excel

## [1.17.7] - 2021-09-15

### Changed

- updated changelog.md
- added excel batch processing functionality
- added slice-chunck function for every type of slice

## [1.17.6] - 2021-09-09

### Changed

- updated changelog.md
- changes excel header colour from blue to black(default-colour)

## [1.17.5] - 2021-07-02

### Changed

- pdf generator html filename from Nano to Unix Nano timestamp

## [1.17.4] - 2021-06-29

### Fixed

- Merged in pdf generator

## [1.17.3] - 2021-06-01

### Updated

- updated column added check if file already exists or not
- updates test cases for excel function.
- updated columns needed for createAndSaveexcel func.

## [1.17.2] - 2021-05-28

### Added

- added test cases for aws s3-upload functions.
- added s3-upload function in createAndSaveexcel func.

## [1.17.1] - 2021-05-26

### Changed

- path in my test-cases for excel-test functions.

### Added

- added test cases for excel functions.
- added excel functionality to create and save excel file on given path.

### Updated

- updated .gitignore
- write database package in scipts to escape it.
- updated package info into README.MD file.

## [1.17.0] - 2021-05-21

### Added

- logger: NewGormLogger method
- database: GORM Microsoft SQL Server connection

## [1.16.1] - 2021-05-11

### Changed

- coverage: now the Skip list uses HasSuffix instead of Contains when deciding if a package needs to be skipped

## [1.16.0] - 2021-05-01

### Added

- database: helper functions to create SQLX Database Pool connection https://github.com/jmoiron/sqlx

## [1.15.3] - 2021-04-29

### Fixed

- Fixed a typo in NewDatabasePoolFromSecret function's name

## [1.15.2] - 2021-03-23

- updated log lines as per PR comments and CHANGELOG.md also.
- added extra fields in middleware to take API_GW_BASE_URL from micro-service where it is used.
- added logger field to get component logger instance which used to get correct logs according to micro-service.

## [1.15.1] - 2021-03-20

### Added

- logger/NewCommonLoggerFromConfiguration a helper function to create the CommonLogger from an AppConfig object

### Changed

- AppConfig Get method to Lookup
- AppConfig MustGet method to Get

### Fixed

- How config loads the .env file(s), now it won't try to load any file on 0 input
- CreateSampleFile now renders a correct .env file (no more inline # comments)

### Updated

- examples/config-global With the current configs

## [1.15.0] - 2021-03-19

### Added

- database package to help create database pools by env, DatabaseInfo or SecretsManager Secret

## [1.14.0] - 2021-03-19

### Added

- models/DatabaseInfo a common model for building DSN-like connection string fro PostgreSQL databases, from environment variables or SecretManager secrets
- APP_LOG_LEVEL APP_LOG_DEV APP_LOG_FORMAT_ERRORS to default configs
- CreateSampleFile function to AppConfig

### Changed

- SecretsManager now uses models/DatabaseInfo to build to connection string from secret

### Removed

- Old database configs and variables

## [1.13.1] - 2021-03-17

### Added

- Logger: auto detect logging level by the LOG_LEVEL env var
- Logger: auto detect development logging (text logger) by the LOG_DEV env var
- Logger: auto detect error formatting by the LOG_FORMAT_ERRORS env var
- An example for a simple global logger

## [1.13.0] - 2021-03-11

### Added

- Alternative SecretsManager client, with unit and integration tests

## [1.12.0] - 2021-03-10

### Added

- APP_DB_SECRET_NAME to App configs
- DBSecretName to AppConfig

### Changed

- config.GetHostName now uses os.Hostname to determine the hostname, instead of the HOSTNAME env var
- AppConfig Hostname now uses GetHostName instead of the APP_HOSTNAME variable

### Removed

- APP_HOSTNAME configuration item
- EC2Identifier from AppConfig (now it will be returned by Hostname if EC2_ID is set)

### Updated

- Examples with new configs

## [1.11.1] - 2021-03-09

### Added

- GetHostName function to config package

## [1.11.0] - 2021-02-26

### Added

- added aws package for Secrets Manager
- added test function for Secrets Manager

## [1.10.1] - 2021-02-11

### Fixed

- ParseDutchPhoneNumber now does not accepts both +31 and 06 prefixes

## [1.10.0] - 2021-02-08

### Added

- ITSUP-971 - Validator package
- ITSUP-971 - ParseDutchPhoneNumber function added to the Validator package

## [1.9.0] - 2021-01-23

### Added

- APP_EC2_ID and APP_HEALTHCHECK_PORT to app configs.

### Changed

- APP_HOST to APP_HOSTNAME in app configs
- APP_DB_PASSWORD to APP_DB_PASS in app configs
- changed the config tests and example accordingly
- constant valid value lists from string to interface{} so validator.In can take them as parameters

### Fixed

- APP_DB_SSL modes according to lib/pq https://www.postgresql.org/docs/current/libpq-ssl.html

## [1.8.5] - 2021-01-18

### Added

- Dev, Test, Staging Old-Staging and Acceptance URLs to GetBaseURL

## [1.8.4] - 2021-01-18

### Added

- Added extra context-key for UMS_Golden_source

## [1.8.5] - 2021-01-20

### Added

- Added extra context-key for forced_logout from Decode Token API in token-based-auth middleware

## [1.8.4] - 2021-01-18

### Added

- Added extra context-key for UMS_Golden_source

## [1.8.3] - 2021-01-10

### Added

- SSL mode to configs and constants and connection string

## [1.8.2] - 2021-01-10

### Updated

- UMS Connector comments

## [1.8.1] - 2021-01-10

### Added

- UMS Connector
- tests to middlewares/RespondWithError

### Modified

- NewMockServer now returns proper http.Response
- Separated constant url-query-keys, service-paths and basic-errors to three different files

## [1.8.0] - 2021-01-08

### Added

- tests package for HTTP testing
- rest client tests
- coverage package for checking Go code test coverage
- scripts/coverage.go to check test coverage
- migration package to manage SQL migrations
- migrator CLI
- an example for using the migrator CLI
- PostgreSQL Docker container to Go tests in Pipelines

### Updated

- Updated README.md with connectors, models, rest, tests and coverage package

## [1.7.2] - 2020-12-18

### Changed

- Refactored package structure: enforcing flat design
- Naming of context keys
- Token paths

## [1.7.1] - 2020-12-17

### Added

- Context keys in middlewares
- Middleware for JWT Token
- Middleware for Label based auth
- Connector for calling Token API of UMS
- Connector for calling Label API of UMS

## [1.7.0] - 2020-12-13

### Added

- config package
- constants package
- examples for cli, config, logger, request-logger middleware

### Changed

- The Logger no longer satisfies the logrus.FieldLogger interface (this is because of the reporting caller issue)

### Updated

- README.md updated

### Removed

- file logger (every service should use only the CommonLogger)

## [1.6.1] - 2020-12-04

### Added

- GetDefaultLogger

## [1.6.0] - 2020-12-04

### Added

- Default logger

## [1.5.0] - 2020-12-04

### Added

- Logger tests

### Changed

- Logging middleware now accepts logrus.FieldLogger instead of a Logrus Logger

## [1.4.0] - 2020-12-03

### Added

- Common Logger
- Component Logger

## [1.3.0] - 2020-12-02

### Added

- CLI Package
- slack.sh script
- Tag pipeline

## [1.2.0] - 2020-11-18

### Added

- FileLogger

## [1.1.0] - 2020-11-18

### Added

- Updated README.md

### Changed

- Extra Logrus fields added to the LoggingMiddleware as parameters

## [1.0.0] - 2020-11-17

### Added

- Bitbucket Pipelines - default and testing pipe

### Fixed

- tests fixed

## [0.0.1] - 2020-11-17

### Added

- README.md
- CHANGELOG.md
- go modules
- request logger middleware

### Fixed

- change the html file naming from Unix to UnixNano
