package handlers

import (
	"context"

	health "github.com/hellofresh/health-go/v5"
	"github.com/labstack/echo/v4"
	"github.com/music-tribe/healthy"
)

func Handler(svc healthy.Service) echo.HandlerFunc {
	opts := make([]health.Config, len(svc.Checkers()))

	i := 0
	for name, checker := range svc.Checkers() {
		opts[i] = health.Config{
			Name:  name,
			Check: getCheckFunc(checker),
		}
		i++
	}

	h, _ := health.New(health.WithComponent(health.Component{
		Name:    svc.Name(),
		Version: svc.Version(),
	}), health.WithChecks(opts...))

	return echo.WrapHandler(h.Handler())
}

func getCheckFunc(c healthy.Checker) func(ctx context.Context) error {
	return c.CheckFunc
}
