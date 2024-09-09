package healthy

import (
	"context"

	"github.com/hellofresh/health-go/v5"
	"github.com/music-tribe/errors"
)

const stubCheckerErrorMessage = "stubChecker: if you're seeing this message, you've not added a `CheckFunc` to your `Checker`!"

func stubChecker() health.CheckFunc {
	return func(ctx context.Context) error {
		return errors.NewCloudError(400, stubCheckerErrorMessage)
	}
}
