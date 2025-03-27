package healthy

import (
	"encoding/json"

	"github.com/hellofresh/health-go/v5"
)

type HealthCheckError struct {
	Code  int
	Check health.Check
}

func NewHealthCheckError(code int, check health.Check) *HealthCheckError {
	he := &HealthCheckError{Code: code, Check: check}
	return he
}

func (he *HealthCheckError) Error() string {
	byt, _ := json.MarshalIndent(he.Check, "", "  ")

	return string(byt)
}
