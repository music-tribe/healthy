package healthy

import (
	"time"

	"github.com/hellofresh/health-go/v5"
)

const defaultName = "healthcheck"

type Checker struct {
	Name      string
	CheckFunc health.CheckFunc
	// Timeout defaults to 2 seconds to match health package: https://github.com/hellofresh/health-go/blob/v5.1.0/health.go#L114
	Timeout time.Duration
}

func NewChecker(name string, checkFunc health.CheckFunc) Checker {
	if name == "" {
		name = defaultName
	}
	if checkFunc == nil {
		checkFunc = stubChecker()
	}
	return Checker{
		Name:      name,
		CheckFunc: checkFunc,
		Timeout:   2 * time.Second,
	}
}

func NewCheckerWithTimeout(name string, checkFunc health.CheckFunc, timeout time.Duration) Checker {
	c := NewChecker(name, checkFunc)
	c.Timeout = timeout
	return c
}
