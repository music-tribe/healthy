package healthy

import (
	"github.com/hellofresh/health-go/v5"
)

const defaultName = "healthcheck"

type Checker struct {
	Name      string
	CheckFunc health.CheckFunc
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
	}
}
