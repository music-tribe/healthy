package handlers

import (
	"context"

	health "github.com/hellofresh/health-go/v5"
	"github.com/labstack/echo/v4"
	"github.com/music-tribe/healthy"
)

// Handler wraps the health-go Handler to be used with Echo.
// However this doesn't propogate errors up the stack.
func Handler(svc healthy.Service) echo.HandlerFunc {
	opts := make([]health.Config, len(svc.Checkers()))

	i := 0
	for name, checker := range svc.Checkers() {
		opts[i] = health.Config{
			Name:    name,
			Check:   getCheckFunc(checker),
			Timeout: checker.Timeout,
		}
		i++
	}

	// FIXME(alex) check the error here
	h, _ := buildHealthCheckContainer(svc)

	return echo.WrapHandler(h.Handler())
}

func getCheckFunc(c healthy.Checker) func(ctx context.Context) error {
	return c.CheckFunc
}

func buildHealthCheckContainer(svc healthy.Service) (*health.Health, error) {
	opts := make([]health.Config, len(svc.Checkers()))

	i := 0
	for name, checker := range svc.Checkers() {
		opts[i] = health.Config{
			Name:    name,
			Check:   getCheckFunc(checker),
			Timeout: checker.Timeout,
		}
		i++
	}

	h, err := health.New(health.WithComponent(health.Component{
		Name:    svc.Name(),
		Version: svc.Version(),
	}), health.WithChecks(opts...))
	return h, err
}
