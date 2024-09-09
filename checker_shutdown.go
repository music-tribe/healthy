package healthy

import (
	"context"

	"github.com/hellofresh/health-go/v5"
	"github.com/music-tribe/errors"
	"github.com/music-tribe/healthy/atomic_bool"
)

// ShutdownChecker allows us to respond to liveness probes to report shutting down
type ShutdownChecker struct {
	isShuttingDown atomic_bool.AtomicBool
	logger         Logger
	component      string
	shutdownLog    string
}

func (c *ShutdownChecker) SetShutdown() {
	c.isShuttingDown.Set(true)
	c.logger.Info(c.component, "SetShutdown", c.shutdownLog)
}

func (c *ShutdownChecker) Get() *atomic_bool.AtomicBool {
	return &c.isShuttingDown
}

func (c *ShutdownChecker) GetShutdown() bool {
	return c.isShuttingDown.Get()
}

func shutdownCheckerFunc(c *ShutdownChecker) health.CheckFunc {
	return func(ctx context.Context) error {
		if c.isShuttingDown.Get() {
			return errors.NewCloudError(503, "shutting down")
		}
		return nil
	}
}

func NewShutdownChecker(name string, logger Logger, component, shutdownLog string) (*ShutdownChecker, Checker) {
	p := &ShutdownChecker{
		isShuttingDown: atomic_bool.AtomicBool{},
		logger:         logger,
		component:      component,
		shutdownLog:    shutdownLog,
	}
	return p, NewChecker(name, shutdownCheckerFunc(p))
}
