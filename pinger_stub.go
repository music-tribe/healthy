package healthy

import (
	"github.com/music-tribe/errors"
)

const stubPingerErrorMessage = "stubPinger: if you're seeing this message, you've not added a `Pinger` to your `Checker`!"

type stubPinger struct{}

func (sp *stubPinger) Ping() error {
	return errors.NewCloudError(400, stubPingerErrorMessage)
}
