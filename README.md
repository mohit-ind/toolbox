# Toolbox

This is a collection of handy packages which can be imported into various Go Applications within the Appventurez organization

## Usage
---
### Development
set the **GOPRIVATE** variable, so Go can download the module:
```BASH
export GOPRIVATE="github.com/toolbox"
```
You may include this in your ```.bashrc``` or ```.bash_profile``` so you do not need to set it every time you want to download a new version of the package.  

---
### Docker
To use it in a **Docker** environment you need to create an [App Password in Bitbucket](https://bitbucket.org/account/settings/app-passwords/) and
```BASH
export DOCKER_NETRC="machine bitbucket.org login your_user_name password your_app_password"
```
Then add this to your **Dockerfile**
```Dockerfile
ARG DOCKER_NETRC

RUN echo "${DOCKER_NETRC}" > ~/.netrc
```
And this to your **docker-compose.yml** (if there is any)
```YAML
services:
  app:
    build:
      args:
        - DOCKER_NETRC
```
---
### CI/CD
In Bitbucket Pipelines set a Secure variable with key: **DOCKER_NETRC** and value: ```"machine bitbucket.org login your_user_name password your_app_password"``` Then and add this script line, before trying to build / test the application
```BASH
echo "${DOCKER_NETRC}" > ~/.netrc
```
---
### Upgrade to a newer version
If you'd like to use a newer version, edit the version string in your ```go.mod``` file
```
require (
    github.com/toolbox v1.6.1
)
```
and run ```go mod tidy```
***  

## Packages

You may import packages from the toolbox with
```go
import (
    "github.com/toolboxcli"
    "github.com/toolboxconfig"
    "github.com/toolboxlogger"
    ...
)
```
### [Tests](tests)
The Tests package is meant to use for HTTP Request and Response testing

The methodology used in this package to test outgoing HTTP Requests is heavily based on [this article](http://hassansin.github.io/Unit-Testing-http-client-in-Go)

**It is very important to keep the toolbox repository's test-coverage above 85%**
 - Keep things flat, do not introduce nested packages and dependencies
 - One function / method should do one thing
 - There should be no side effects of a function / method calls (if there is any there should be a way to test it)
 - If a package / object / method / function introduces a dependency (external URLs, database , logger, etc.) it has to provide a convenient way to inject that said dependency

---
### [Coverage](coverage)
Coverage is a small package that can analyse the output of the `go test` command (check how it is called in the Makefile -> test/ci-test.sh). The point of this package is to check in CI if a repository's test coverage has been meet with the required standards.

---
### [Constants](constants)
The constants package provides constant values that all application should use. These are mainly environment variable names

---
### [Config](config)
The config package provides configuration primitives and the AppConfig object, which can be used to load, validate and retrieve configuration items

---
### [CLI](cli)
The cli package provides dead simple tools to build a command line interface for your application

---
### [Validator](validator)
The validator provides a set of validation functions, in alignment with go-ozzo/ozzo-validation.

---
### [File](file)
The file package provides utility functions to manage files on the local filesystem.

---
### [Logger](logger)
The logger package provides a common logger which should be used by all services. It requires the service-name, version, environment and hostname to be set. These fields will be added to all log entries. In debug mode every log entry will contain the caller function with filename and line-number.

Use ```github.com/pkg/errors``` to wrap and propagate errors in your application. Use the logger's WithError method to log errors from the application (this will allow the unwrapping of errors, with correct error-trace)

---
### [Connectors](connectors)
The Connectors package provides different connectors across the Appventurez APIs  

---
### [Models](models)
The Models package contains commonly used Go objects with appropriate Json and db Struct Tags for encoding

---
### [Rest](rest)
The Rest package provides a simple HTTP Client to interact with external services and Appventurez APIs (imho it is better to use a package like [Sling](https://github.com/dghubble/sling) or [Gentleman]())

---
### [Services](services)
Package services provides various new services using golang 3rd party and standard libraries

---
#### [Email](services/email)
Package email provides all essential object and settings to send email using smtp

---
### [shell](shell)
Package shell provides a command structure which is able to execute shell commands

---
### [Middlewares](middlewares)
The middlewares package provides useful middleware functions for web applications. The **LoggingMiddleware** accepts a toolbox/logger.Logger as an argument, and creates a middleware that logs every Request with path, method, status and duration.  

---
### [Database](database)
The Database package provides helper methods to connect to an SQL database, ping it, set the connection pool and get the connection object.

---
### [Migration](migration)
Package migration provides a Migrator object (with a command line interface) which can execute and roll back SQL migration scripts against a PostgreSQL database  

---

### [Document - generator](document-generator)
This package provides function to create and save excel and pdf files.

---
#### [Excel](document-generator/excel)
This package have model which is used to create and save excel file on given path. 

---
#### [PDF](document-generator/pdf)
This package have model which is used to create and save pdf file on given path. 

---

####
## Examples

There is an example for the usage of every package under [examples](examples)

## Commitment
You are encouraged to use the toolbox repository all across the Go web services within appventurez. **Pull Requests** are welcome, but they must have a good reason to be in the toolbox, and **they have to be covered at least 85% by tests**.

## [AWS](Secrets Manager)
Package AWS provides common function to default external configurations and, populates an AWS Config with the values from the external configurations also provide one function to get the values from AWS Secrets-Manager.