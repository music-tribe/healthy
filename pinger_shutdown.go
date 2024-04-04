package healthy

import (
	"github.com/music-tribe/errors"
	"github.com/music-tribe/healthy/atomic_bool"
)

// ShutdownPinger allows us to respond to liveness probes to report shutting down
type ShutdownPinger struct {
	isShuttingDown atomic_bool.AtomicBool
	logger         Logger
	component      string
	shutdownLog    string
}

func NewShutdownPinger(logger Logger, component, shutdownLog string) *ShutdownPinger {
	return &ShutdownPinger{
		isShuttingDown: atomic_bool.AtomicBool{},
		logger:         logger,
		component:      component,
		shutdownLog:    shutdownLog,
	}
}

func (p *ShutdownPinger) SetShutdown() {
	p.isShuttingDown.Set(true)
	p.logger.Info(p.component, "SetShutdown", p.shutdownLog)
}

func (p *ShutdownPinger) Get() *atomic_bool.AtomicBool {
	return &p.isShuttingDown
}

func (p *ShutdownPinger) GetShutdown() bool {
	return p.isShuttingDown.Get()
}

func (p *ShutdownPinger) Ping() error {
	if p.isShuttingDown.Get() {
		return errors.NewCloudError(503, "shutting down")
	}
	return nil
}

func NewShutdownChecker(name string, pinger *ShutdownPinger) Checker {
	return NewChecker(name, pinger)
}
