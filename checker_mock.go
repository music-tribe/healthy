package healthy

import (
	"context"

	"github.com/hellofresh/health-go/v5"
)

func NewMockChecker(err error) health.CheckFunc {
	return func(ctx context.Context) error {
		return err
	}
}
