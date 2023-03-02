# Healthy
Healthy is a simple adapter that enables us to easily check the healh of our services and their dependencies.
The HTTP handler side is built on top of the excellent Health Go package and is adapted to our router of choice - Echo.

## In Use
### Installation
```go get github.com/music-tribe/healthy```

### Setup
Lets say you wanted to check the health of a database. To begin with we need to create a `Ping() error` method on that database - the internals of this method are up to you, you just want to confirm that you can connect to that database...
```golang

type MyDatabase struct{}

func (d *MyDatabase)Ping() error {
    // confirm your connection in here making sure to return any errors
}
```
The `Ping` method that we've added means that our database now satisfies the `Pinger` interface and we can create our health check as follows...
```golang
package main

import (
    "github.com/music-tribe/healthy"
)

db := new(MyDatabase)
databaseCheck := healthy.NewChecker("database", db)

healthService := healthy.New("my-service", "1.45.7", databaseCheck)

api := NewAPI(healthService)
if err := api.Run("8080"); err != nil {
    panic(err)
}
```
We can pass then pass this `healthService` into our echo Handler within our API...
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