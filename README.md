# Healthy
Healthy is a simple adapter that enables us to easily check the healh of our services and their dependencies.
The HTTP handler side is built on top of the excellent Health Go package and is adapted to our router of choice - Echo.

## In Use
### Installation
```go get github.com/music-tribe/healthy```

### Setup
Lets say you wanted to check the health of a database. To begin with we need to create a `CheckFunc` method of the form `func(context.Context) error` on that database - the internals of this method are up to you, you just want to confirm that you can connect to that database...
```golang

type MyDatabase struct{}

func (d *MyDatabase) CheckFunc() func(context.Context) error {
    // confirm your connection in here making sure to return any errors
}
```
We can now create our health check as follows...
```golang
package main

import (
    "github.com/music-tribe/healthy"
)

db := new(MyDatabase)
databaseCheck := healthy.NewChecker("database", db.CheckFunc())

healthService := healthy.New("my-service", "1.45.7", databaseCheck)

api := NewAPI(healthService)
if err := api.Run("8080"); err != nil {
    panic(err)
}
```
We can pass this `healthService` into our echo Handler within our API...
```golang
package api

import (
    "github.com/labstack/echo/v4"
    "github.com/music-tribe/healthy/handlers"
)

type api struct{
    healthService healthy.Service
}

func NewApi(hSvc healthy.Service) ports.API {
    return &api{
        healthService: hSvc,
    }
}

func (a *api) Run(port string) error {
    e := echo.New()

    // pass the healthy handler in here...
    e.GET("/healthz", handlers.Handler(api.healthService))

    return e.Start(":"+port)
}

```

#### Shutdown Checker
This package provides a `ShutdownChecker`. Use this to respond to liveness probes whilst shutting down. e.g.

```golang
package main

import (
    "github.com/music-tribe/healthy"
)

isShuttingDownChecker, shutdownChecker := healthy.NewShutdownChecker("tunnelShutdownChecker", logger, "tunnel", "tunnel service is shutting down")
healthService, err := healthy.New(
    serviceName,
    version,
    healthy.NewMongoDbCheckerWithConnectionString("database", dbConnString),
    shutdownChecker,
)
if err != nil {
    return err
}

//
package api

import (
    "github.com/labstack/echo/v4"
    "github.com/music-tribe/healthy/handlers"
)

type api struct{
    healthService   healthy.Service
    isShuttingDownChecker *healthy.ShutdownChecker
}

func NewApi(hSvc healthy.Service, isShuttingDownChecker *healthy.ShutdownChecker) ports.API {
    return &api{
        healthService:  hSvc,
        isShuttingDownChecker: isShuttingDownChecker,
    }
}

func (api *api) Shutdown(ctx context.Context) error {
	api.isShuttingDownChecker.SetShutdown()
	return api.router.Shutdown(ctx)
}
```

## Development

To run the test suite use command `go test -v`.

> NOTE: generated mocks are added to source control so downstream packages can compile and test this. To update mocks run `make mocks`.

### Pre-commit

Install [pre-commit](https://pre-commit.com/) with `pre-commit install`.

Ensure you have the tools required:
* [golangci-lint](https://golangci-lint.run/)
* [goreturns](https://github.com/sqs/goreturns)
* [gosec](https://github.com/securego/gosec)
* [staticcheck](https://staticcheck.dev/)