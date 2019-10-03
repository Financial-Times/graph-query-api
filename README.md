# graph-query-api
_Should be the same as the github repo name but it isn't always._

[![Circle CI](https://circleci.com/gh/Financial-Times/graph-query-api/tree/master.png?style=shield)](https://circleci.com/gh/Financial-Times/graph-query-api/tree/master)[![Go Report Card](https://goreportcard.com/badge/github.com/Financial-Times/graph-query-api)](https://goreportcard.com/report/github.com/Financial-Times/graph-query-api) [![Coverage Status](https://coveralls.io/repos/github/Financial-Times/graph-query-api/badge.svg)](https://coveralls.io/github/Financial-Times/graph-query-api)

## Introduction

_What is this service and what is it for? What other services does it depend on_

UPP Golang Microservice Template short description - please amend

## Installation
      
_How can I install it_

Download the source code, dependencies and test dependencies:

        curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
        go get -u github.com/Financial-Times/graph-query-api
        cd $GOPATH/src/github.com/Financial-Times/graph-query-api
        dep ensure
        go build .

## Running locally
_How can I run it_

1. Run the tests and install the binary:

        dep ensure
        go test -race ./...
        go install

2. Run the binary (using the `help` flag to see the available optional arguments):

        $GOPATH/bin/graph-query-api [--help]

Options:

        --app-system-code="graph-query-api"            System Code of the application ($APP_SYSTEM_CODE)
        --app-name="graph-query-api"                   Application name ($APP_NAME)
        --port="8080"                                           Port to listen on ($APP_PORT)

## Build and deployment
_How can I build and deploy it (lots of this will be links out as the steps will be common)_

* Built by Docker Hub on merge to master: [coco/graph-query-api](https://hub.docker.com/r/coco/graph-query-api/)
* CI provided by CircleCI: [graph-query-api](https://circleci.com/gh/Financial-Times/graph-query-api)

## Service endpoints
_What are the endpoints offered by the service_

e.g.
### GET

Using curl:

    curl http://localhost:8080/_<INSERT SEPCIFIC URL HERE>_| json_pp`

_Explain what the response should represent_

Based on the following [google doc](_<INSERT API DOCUMETATION HERE>_)

## Utility endpoints
_Endpoints that are there for support or testing, e.g read endpoints on the writers_

## Healthchecks
Admin endpoints are:

`/__gtg`

`/__health`

`/__build-info`

_These standard endpoints do not need to be specifically documented._

_This section *should* however explain what checks are done to determine health and gtg status._

There are several checks performed:

_e.g._
* Checks that a connection can be made to Neo4j, using the neo4j url supplied as a parameter in service startup.

## Other information
_Anything else you want to add._

_e.g. (NB: this example may be something we want to extract as it's probably common to a lot of services)_

### Logging

* The application uses [logrus](https://github.com/sirupsen/logrus); the log file is initialised in [main.go](main.go).
* Logging requires an `env` app parameter, for all environments other than `local` logs are written to file.
* When running locally, logs are written to console. If you want to log locally to file, you need to pass in an env parameter that is != `local`.
* NOTE: `/__build-info` and `/__gtg` endpoints are not logged as they are called every second from varnish/vulcand and this information is not needed in logs/splunk.